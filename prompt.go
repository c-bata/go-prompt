package prompt

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"os"
	"sync"
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
	renderer          *Render
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

	// Run SignalHandler goroutine to receive OS signals that window size is changed and kill this process.
	sh := NewSignalHandler()
	go sh.Run(ctx, cancel)

	// Run renderer process
	wg.Add(1)
	go func() {
		p.renderer.Run(ctx, p.buf, p.completion, p.in.GetWinSize())
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
			if shouldExit, e := p.feed(b); shouldExit {
				return
			} else if e != nil {
				// Stop goroutine to run readBuffer function to unset raw mode and non-blocking mode.
				// Because returned EAGAIN when still set non-blocking mode.
				ip.Pause <- true
				p.executor(e.input)

				p.completion.Update(*p.buf.Document())
				p.renderer.Render <- RenderRequest{
					buffer:     p.buf,
					completion: p.completion,
				}
				ip.Pause <- false
			} else {
				p.completion.Update(*p.buf.Document())
				p.renderer.Render <- RenderRequest{
					buffer:     p.buf,
					completion: p.completion,
				}
			}
		case <-sh.SigWinch:
			p.renderer.WinSize <- p.in.GetWinSize()
		case <-ctx.Done():
			return
		}
	}
}

func (p *Prompt) feed(b []byte) (shouldExit bool, exec *Exec) {
	key := p.in.GetKey(b)

	// completion
	completing := p.completion.Completing()
	p.handleCompletionKeyBinding(key, completing)

	switch key {
	case Enter, ControlJ, ControlM:
		p.renderer.breakLine(p.buf)

		exec = &Exec{input: p.buf.Text()}
		log.Printf("[History] %s", p.buf.Text())
		p.buf = NewBuffer()
		if exec.input != "" {
			p.history.Add(exec.input)
		}
	case ControlC:
		p.renderer.breakLine(p.buf)
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

	// Run SignalHandler goroutine to receive OS signals that window size is changed and kill this process.
	sh := NewSignalHandler()
	go sh.Run(ctx, cancel)

	// Run renderer goroutine to render the buffer, suggests.
	wg.Add(1)
	go func() {
		p.renderer.Run(ctx, p.buf, p.completion, p.in.GetWinSize())
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
			if shouldExit, e := p.feed(b); shouldExit {
				return ""
			} else if e != nil {
				return e.input
			} else {
				p.completion.Update(*p.buf.Document())
				p.renderer.Render <- RenderRequest{
					buffer:     p.buf,
					completion: p.completion,
				}
			}
		case <-sh.SigWinch:
			p.renderer.WinSize <- p.in.GetWinSize()
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}
