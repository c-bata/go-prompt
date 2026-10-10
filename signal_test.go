// +build !windows

package prompt

import (
	"io"
	"sync"
	"syscall"
	"testing"
	"time"
)

type testWriter struct {
	VT100Writer
}

func (w *testWriter) Flush() error {
	w.buffer = nil
	return nil
}

var _ ConsoleWriter = &testWriter{}

type testParser struct {
	mu      sync.Mutex
	inputs  [][]byte
	winSize *WinSize
}

func (p *testParser) Setup() error {
	return nil
}

func (p *testParser) TearDown() error {
	return nil
}

func (p *testParser) GetWinSize() *WinSize {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.winSize != nil {
		return p.winSize
	}
	return &WinSize{Row: 24, Col: 80}
}

func (p *testParser) Read() ([]byte, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.inputs) > 0 {
		b := p.inputs[0]
		p.inputs = p.inputs[1:]
		return b, nil
	}
	return nil, io.EOF
}

func (p *testParser) appendInput(b []byte) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.inputs = append(p.inputs, b)
}

func TestHandleSignals_WinSizeNonBlocking(t *testing.T) {
	parser := &testParser{
		winSize: &WinSize{Row: 30, Col: 100},
	}
	pt := &Prompt{
		in: parser,
	}

	exitCh := make(chan int)
	// Unbuffered channel with no reader to simulate reader busy or unbuffered channel.
	winSizeCh := make(chan *WinSize)
	stopCh := make(chan struct{})
	doneCh := make(chan struct{})

	go func() {
		pt.handleSignals(exitCh, winSizeCh, stopCh)
		close(doneCh)
	}()

	// Send SIGWINCH signal. In previous versions with blocking send, this would hang.
	if err := syscall.Kill(syscall.Getpid(), syscall.SIGWINCH); err != nil {
		t.Fatalf("failed to send SIGWINCH: %v", err)
	}

	// Allow a brief moment for the signal handler to process SIGWINCH.
	time.Sleep(50 * time.Millisecond)

	// Stopping handleSignals must succeed promptly without deadlock.
	stopCh <- struct{}{}

	select {
	case <-doneCh:
		// Succeeded
	case <-time.After(2 * time.Second):
		t.Fatal("handleSignals deadlocked on SIGWINCH with unbuffered winSizeCh")
	}
}

func TestHandleSignals_RapidConcurrentSignals(t *testing.T) {
	parser := &testParser{
		winSize: &WinSize{Row: 40, Col: 120},
	}
	pt := &Prompt{
		in: parser,
	}

	exitCh := make(chan int)
	winSizeCh := make(chan *WinSize, 128)
	stopCh := make(chan struct{})
	doneCh := make(chan struct{})

	go func() {
		pt.handleSignals(exitCh, winSizeCh, stopCh)
		close(doneCh)
	}()

	// Read from winSizeCh in background
	stopReader := make(chan struct{})
	readerDone := make(chan struct{})
	go func() {
		defer close(readerDone)
		for {
			select {
			case <-stopReader:
				return
			case <-winSizeCh:
			}
		}
	}()

	// Fire rapid SIGWINCH signals concurrently
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				_ = syscall.Kill(syscall.Getpid(), syscall.SIGWINCH)
				time.Sleep(2 * time.Millisecond)
			}
		}()
	}
	wg.Wait()

	// Stop handleSignals
	stopCh <- struct{}{}

	select {
	case <-doneCh:
		// Succeeded
	case <-time.After(2 * time.Second):
		t.Fatal("handleSignals deadlocked during rapid concurrent SIGWINCH signals")
	}

	close(stopReader)
	<-readerDone
}

func TestPrompt_Run_ConcurrentSIGWINCH(t *testing.T) {
	parser := &testParser{
		winSize: &WinSize{Row: 25, Col: 80},
	}
	writer := &testWriter{}

	executed := make(chan string, 5)
	executor := func(in string) {
		executed <- in
	}

	pt := &Prompt{
		in: parser,
		renderer: &Render{
			prefix:                       "> ",
			out:                          writer,
			livePrefixCallback:           func() (string, bool) { return "", false },
			prefixTextColor:              Blue,
			prefixBGColor:                DefaultColor,
			inputTextColor:               DefaultColor,
			inputBGColor:                 DefaultColor,
			previewSuggestionTextColor:   Green,
			previewSuggestionBGColor:     DefaultColor,
			suggestionTextColor:          White,
			suggestionBGColor:            Cyan,
			selectedSuggestionTextColor:  Black,
			selectedSuggestionBGColor:    Turquoise,
			descriptionTextColor:         Black,
			descriptionBGColor:           Turquoise,
			selectedDescriptionTextColor: White,
			selectedDescriptionBGColor:   Cyan,
			scrollbarThumbColor:          DarkGray,
			scrollbarBGColor:             Cyan,
		},
		buf:         NewBuffer(),
		executor:    executor,
		history:     NewHistory(),
		completion:  NewCompletionManager(func(Document) []Suggest { return nil }, 6),
		keyBindMode: EmacsKeyBind,
	}

	runDone := make(chan struct{})
	go func() {
		pt.Run()
		close(runDone)
	}()

	// Fire SIGWINCH signals concurrently while sending inputs
	stopSignals := make(chan struct{})
	signalsDone := make(chan struct{})
	go func() {
		defer close(signalsDone)
		for {
			select {
			case <-stopSignals:
				return
			default:
				_ = syscall.Kill(syscall.Getpid(), syscall.SIGWINCH)
				time.Sleep(5 * time.Millisecond)
			}
		}
	}()

	// Queue input to execute a command: first the command text, then Enter (0x0a)
	parser.appendInput([]byte("test_command"))
	parser.appendInput([]byte{0xa})

	select {
	case cmd := <-executed:
		if cmd != "test_command" {
			t.Fatalf("expected 'test_command', got %q", cmd)
		}
	case <-time.After(3 * time.Second):
		close(stopSignals)
		t.Fatal("prompt runner deadlocked while executing command with concurrent SIGWINCH")
	}

	// Send Control-D to exit prompt runner
	parser.appendInput([]byte{0x04})

	close(stopSignals)
	<-signalsDone

	select {
	case <-runDone:
		// Prompt exited cleanly
	case <-time.After(3 * time.Second):
		t.Fatal("prompt runner deadlocked on shutdown with concurrent SIGWINCH")
	}
}
