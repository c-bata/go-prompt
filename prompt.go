package prompt

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	logfile      = "/tmp/go-prompt-debug.log"
	envEnableLog = "GO_PROMPT_ENABLE_LOG"
)

type Executor func(string)
type Completer func(string) []Suggest

type Prompt struct {
	in         ConsoleParser
	buf        *Buffer
	renderer   *Render
	executor   Executor
	history    *History
	completion *CompletionManager
}

type Exec struct {
	input string
}

func (p *Prompt) Run() {
	p.setUp()
	defer p.tearDown()

	p.renderer.Render(p.buf, p.completion)

	bufCh := make(chan []byte, 128)
	stopReadBufCh := make(chan struct{})
	go readBuffer(bufCh, stopReadBufCh)

	exitCh := make(chan int)
	winSizeCh := make(chan *WinSize)
	go handleSignals(p.in, exitCh, winSizeCh)

	for {
		select {
		case b := <-bufCh:
			if shouldExit, e := p.feed(b); shouldExit {
				p.renderer.BreakLine(p.buf)
				return
			} else if e != nil {
				// Stop goroutine to run readBuffer function
				stopReadBufCh <- struct{}{}

				// Unset raw mode
				// Reset to Blocking mode because returned EAGAIN when still set non-blocking mode.
				p.in.TearDown()
				p.executor(e.input)

				p.completion.Update(p.buf.Text())
				p.renderer.Render(p.buf, p.completion)

				// Set raw mode
				p.in.Setup()
				go readBuffer(bufCh, stopReadBufCh)
			} else {
				p.completion.Update(p.buf.Text())
				p.renderer.Render(p.buf, p.completion)
			}
		case w := <-winSizeCh:
			p.renderer.UpdateWinSize(w)
			p.renderer.Render(p.buf, p.completion)
		case code := <-exitCh:
			p.renderer.BreakLine(p.buf)
			p.tearDown()
			os.Exit(code)
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func (p *Prompt) feed(b []byte) (shouldExit bool, exec *Exec) {
	key := p.in.GetKey(b)

	switch key {
	case ControlJ, Enter:
		if s, ok := p.completion.GetSelectedSuggestion(); ok {
			w := p.buf.Document().GetWordBeforeCursor()
			if w != "" {
				p.buf.DeleteBeforeCursor(len([]rune(w)))
			}
			p.buf.InsertText(s.Text, false, true)
		}
		p.renderer.BreakLine(p.buf)

		exec = &Exec{input: p.buf.Text()}
		log.Printf("[History] %s", p.buf.Text())
		p.buf = NewBuffer()
		p.completion.Reset()
		if exec.input != "" {
			p.history.Add(exec.input)
		}
	case ControlC:
		p.renderer.BreakLine(p.buf)
		p.buf = NewBuffer()
		p.completion.Reset()
		p.history.Clear()
	case ControlD:
		shouldExit = true
	case Up:
		if !p.completion.Completing() {
			if newBuf, changed := p.history.Older(p.buf); changed {
				p.buf = newBuf
			}
			return
		}
		fallthrough
	case BackTab:
		p.completion.Previous()
	case Down:
		if !p.completion.Completing() {
			if newBuf, changed := p.history.Newer(p.buf); changed {
				p.buf = newBuf
			}
			return
		}
		fallthrough
	case Tab, ControlI:
		p.completion.Next()
	case Left:
		p.buf.CursorLeft(1)
	case Right:
		p.buf.CursorRight(1)
	case Backspace:
		if s, ok := p.completion.GetSelectedSuggestion(); ok {
			w := p.buf.Document().GetWordBeforeCursor()
			if w != "" {
				p.buf.DeleteBeforeCursor(len([]rune(w)))
			}
			p.buf.InsertText(s.Text, false, true)
			p.completion.Reset()
		}
		p.buf.DeleteBeforeCursor(1)
	case NotDefined:
		if s, ok := p.completion.GetSelectedSuggestion(); ok {
			w := p.buf.Document().GetWordBeforeCursor()
			if w != "" {
				p.buf.DeleteBeforeCursor(len([]rune(w)))
			}
			p.buf.InsertText(s.Text, false, true)
		}
		p.completion.Reset()
		p.buf.InsertText(string(b), false, true)
	default:
		p.completion.Reset()
	}
	return
}

func (p *Prompt) Input() string {
	p.setUp()
	defer p.tearDown()

	p.renderer.Render(p.buf, p.completion)
	bufCh := make(chan []byte, 128)
	stopReadBufCh := make(chan struct{})
	go readBuffer(bufCh, stopReadBufCh)

	for {
		select {
		case b := <-bufCh:
			if shouldExit, e := p.feed(b); shouldExit {
				p.renderer.BreakLine(p.buf)
				return ""
			} else if e != nil {
				// Stop goroutine to run readBuffer function
				stopReadBufCh <- struct{}{}
				return e.input
			} else {
				p.completion.Update(p.buf.Text())
				p.renderer.Render(p.buf, p.completion)
			}
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func (p *Prompt) setUp() {
	// Logging
	if os.Getenv(envEnableLog) != "true" {
		log.SetOutput(ioutil.Discard)
	} else if f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666); err != nil {
		log.SetOutput(ioutil.Discard)
	} else {
		defer f.Close()
		log.SetOutput(f)
		log.Println("[INFO] Logging is enabled.")
	}

	p.in.Setup()
	p.renderer.Setup()
	p.renderer.UpdateWinSize(p.in.GetWinSize())
}

func (p *Prompt) tearDown() {
	p.in.TearDown()
	p.renderer.TearDown()
}

func readBuffer(bufCh chan []byte, stopCh chan struct{}) {
	buf := make([]byte, 1024)

	log.Printf("[INFO] readBuffer start")
	for {
		time.Sleep(10 * time.Millisecond)
		select {
		case <-stopCh:
			log.Print("[INFO] stop readBuffer")
			return
		default:
			if n, err := syscall.Read(syscall.Stdin, buf); err == nil {
				bufCh <- buf[:n]
			}
		}
	}
}

func handleSignals(in ConsoleParser, exitCh chan int, winSizeCh chan *WinSize) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(
		sigCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGWINCH,
	)

	for {
		s := <-sigCh
		switch s {
		case syscall.SIGINT: // kill -SIGINT XXXX or Ctrl+c
			log.Println("[SIGNAL] Catch SIGINT")
			exitCh <- 0

		case syscall.SIGTERM: // kill -SIGTERM XXXX
			log.Println("[SIGNAL] Catch SIGTERM")
			exitCh <- 1

		case syscall.SIGQUIT: // kill -SIGQUIT XXXX
			log.Println("[SIGNAL] Catch SIGQUIT")
			exitCh <- 0

		case syscall.SIGWINCH:
			log.Println("[SIGNAL] Catch SIGWINCH")
			winSizeCh <- in.GetWinSize()
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}
