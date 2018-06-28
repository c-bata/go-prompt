package prompt

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	envDebugLogPath = "GO_PROMPT_LOG_PATH"
)

// Executor is called when user input something text.
type Executor func(string)

// Completer should return the suggest item from Document.
type Completer func(Document) []Suggest

// Prompt is core struct of go-prompt.
type Prompt struct {
	in                ConsoleParser
	buf               *Buffer
	rendererOptions   []RendererOption
	executor          Executor
	history           *History
	completion        *CompletionManager
	keyBindings       []KeyBind
	ASCIICodeBindings []ASCIICodeBind
	keyBindMode       KeyBindMode
}

// Exec is the struct contains user input context.
type Exec struct {
	input string
}

// Run starts prompt.
func (p *Prompt) Run() {
	if l := os.Getenv(envDebugLogPath); l == "" {
		log.SetOutput(ioutil.Discard)
	} else if f, err := os.OpenFile(l, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666); err != nil {
		log.SetOutput(ioutil.Discard)
	} else {
		defer f.Close()
		log.SetOutput(f)
		log.Println("[INFO] Logging is enabled.")
	}

	var wg sync.WaitGroup
	defer wg.Wait()

	// Application context. If canceled, all worker goroutine stopped.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Run signal handler goroutine to receive OS signals for quitting process.
	sigquit := make(chan os.Signal, 1)
	signal.Notify(sigquit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-sigquit:
				cancel()
			}
		}
	}()

	// Run renderer process
	renderer := NewRenderer(consoleWriter, RenderRequest{buffer: p.buf, completion: p.completion}, p.rendererOptions...)
	wg.Add(1)
	go func() {
		renderer.Run(ctx)
		wg.Done()
	}()

	// Run InputProcessor goroutine to read user input from Keyboard.
	ip := NewInputProcessor(p.in)
	wg.Add(1)
	go func() {
		ip.Run(ctx)
		wg.Done()
	}()

	for {
		select {
		case b := <-ip.UserInput:
			if shouldExit, e := p.feed(renderer, b); shouldExit {
				return
			} else if e != nil {
				// Stop goroutine to run readBuffer function to unset raw mode and non-blocking mode.
				// Because returned EAGAIN when still set non-blocking mode.
				ip.Pause <- true
				p.executor(e.input)

				p.completion.Update(*p.buf.Document())
				renderer.Render <- RenderRequest{
					buffer:     p.buf,
					completion: p.completion,
				}
				ip.Pause <- false
			} else {
				p.completion.Update(*p.buf.Document())
				renderer.Render <- RenderRequest{
					buffer:     p.buf,
					completion: p.completion,
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (p *Prompt) feed(renderer *Renderer, b []byte) (shouldExit bool, exec *Exec) {
	key := GetKey(b)

	// completion
	completing := p.completion.Completing()
	p.handleCompletionKeyBinding(key, completing)

	switch key {
	case Enter, ControlJ, ControlM:
		renderer.breakLine(p.buf)

		exec = &Exec{input: p.buf.Text()}
		log.Printf("[History] %s", p.buf.Text())
		p.buf = NewBuffer()
		if exec.input != "" {
			p.history.Add(exec.input)
		}
	case ControlC:
		renderer.breakLine(p.buf)
		p.buf = NewBuffer()
		p.history.Clear()
	case Up, ControlP:
		if !completing { // Don't use p.completion.Completing() because it takes double operation when switch to selected=-1.
			if newBuf, changed := p.history.Older(p.buf); changed {
				p.buf = newBuf
			}
		}
	case Down, ControlN:
		if !completing { // Don't use p.completion.Completing() because it takes double operation when switch to selected=-1.
			if newBuf, changed := p.history.Newer(p.buf); changed {
				p.buf = newBuf
			}
			return
		}
	case ControlD:
		if p.buf.Text() == "" {
			shouldExit = true
			return
		}
	case NotDefined:
		if p.handleASCIICodeBinding(b) {
			return
		}
		p.buf.InsertText(string(b), false, true)
	}

	p.handleKeyBinding(key)
	return
}

func (p *Prompt) handleCompletionKeyBinding(key Key, completing bool) {
	switch key {
	case Down:
		if completing {
			p.completion.Next()
		}
	case Tab, ControlI:
		p.completion.Next()
	case Up:
		if completing {
			p.completion.Previous()
		}
	case BackTab:
		p.completion.Previous()
	default:
		if s, ok := p.completion.GetSelectedSuggestion(); ok {
			w := p.buf.Document().GetWordBeforeCursorUntilSeparator(p.completion.wordSeparator)
			if w != "" {
				p.buf.DeleteBeforeCursor(len([]rune(w)))
			}
			p.buf.InsertText(s.Text, false, true)
		}
		p.completion.Reset()
	}
}

func (p *Prompt) handleKeyBinding(key Key) {
	for i := range commonKeyBindings {
		kb := commonKeyBindings[i]
		if kb.Key == key {
			kb.Fn(p.buf)
		}
	}

	if p.keyBindMode == EmacsKeyBind {
		for i := range emacsKeyBindings {
			kb := emacsKeyBindings[i]
			if kb.Key == key {
				kb.Fn(p.buf)
			}
		}
	}

	// Custom key bindings
	for i := range p.keyBindings {
		kb := p.keyBindings[i]
		if kb.Key == key {
			kb.Fn(p.buf)
		}
	}
}

func (p *Prompt) handleASCIICodeBinding(b []byte) bool {
	checked := false
	for _, kb := range p.ASCIICodeBindings {
		if bytes.Compare(kb.ASCIICode, b) == 0 {
			kb.Fn(p.buf)
			checked = true
		}
	}
	return checked
}

// Input just returns user input text.
func (p *Prompt) Input() string {
	if l := os.Getenv(envDebugLogPath); l == "" {
		log.SetOutput(ioutil.Discard)
	} else if f, err := os.OpenFile(l, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666); err != nil {
		log.SetOutput(ioutil.Discard)
	} else {
		defer f.Close()
		log.SetOutput(f)
		log.Println("[INFO] Logging is enabled.")
	}

	var wg sync.WaitGroup
	defer wg.Wait()

	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Run signal handler goroutine to receive OS signals for quitting process.
	sigquit := make(chan os.Signal, 1)
	signal.Notify(sigquit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-sigquit:
				close(sigquit)
				signal.Stop(sigquit)
				cancel()
			}
		}
	}()

	// Run renderer process
	renderer := NewRenderer(consoleWriter, RenderRequest{buffer: p.buf, completion: p.completion}, p.rendererOptions...)
	wg.Add(1)
	go func() {
		renderer.Run(ctx)
		wg.Done()
	}()

	// Run InputProcessor goroutine to read user input from Keyboard.
	ip := NewInputProcessor(p.in)

	wg.Add(1)
	go func() {
		ip.Run(ctx)
		wg.Done()
	}()

	for {
		select {
		case b := <-ip.UserInput:
			if shouldExit, e := p.feed(renderer, b); shouldExit {
				return ""
			} else if e != nil {
				return e.input
			} else {
				p.completion.Update(*p.buf.Document())
				renderer.Render <- RenderRequest{
					buffer:     p.buf,
					completion: p.completion,
				}
			}
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}
