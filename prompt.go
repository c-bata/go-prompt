package prompt

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"os"
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
	ctx               context.Context
	cancel            context.CancelFunc
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

	p.setUp()
	defer p.tearDown()

	p.renderer.Render(p.buf, p.completion)

	bufchan := make(chan []byte, 128)
	go p.readBuffer(p.ctx, bufchan)

	exitCh := make(chan int)
	winchan := make(chan *WinSize)
	go p.handleSignals(p.ctx, p.cancel, winchan)

	for {
		select {
		case b := <-bufchan:
			if shouldExit, e := p.feed(b); shouldExit {
				p.renderer.BreakLine(p.buf)
				p.cancel()
				return
			} else if e != nil {
				// Stop goroutine to run readBuffer function
				p.cancel()

				// Unset raw mode
				// Reset to Blocking mode because returned EAGAIN when still set non-blocking mode.
				p.in.TearDown()
				p.executor(e.input)

				p.completion.Update(*p.buf.Document())
				p.renderer.Render(p.buf, p.completion)

				// Set raw mode
				p.in.Setup()
				ctx, cancel := context.WithCancel(context.Background())
				p.ctx = ctx
				p.cancel = cancel
				go p.readBuffer(p.ctx, bufchan)
				go p.handleSignals(p.ctx, p.cancel, winchan)
			} else {
				p.completion.Update(*p.buf.Document())
				p.renderer.Render(p.buf, p.completion)
			}
		case w := <-winchan:
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

	// completion
	completing := p.completion.Completing()
	p.handleCompletionKeyBinding(key, completing)

	switch key {
	case Enter, ControlJ, ControlM:
		p.renderer.BreakLine(p.buf)

		exec = &Exec{input: p.buf.Text()}
		log.Printf("[History] %s", p.buf.Text())
		p.buf = NewBuffer()
		if exec.input != "" {
			p.history.Add(exec.input)
		}
	case ControlC:
		p.renderer.BreakLine(p.buf)
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

	p.setUp()
	defer p.tearDown()

	p.renderer.Render(p.buf, p.completion)
	bufchan := make(chan []byte, 128)
	go p.readBuffer(p.ctx, bufchan)

	for {
		select {
		case b := <-bufchan:
			if shouldExit, e := p.feed(b); shouldExit {
				p.renderer.BreakLine(p.buf)
				p.cancel()
				return ""
			} else if e != nil {
				return e.input
			} else {
				p.completion.Update(*p.buf.Document())
				p.renderer.Render(p.buf, p.completion)
			}
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func (p *Prompt) readBuffer(ctx context.Context, bufCh chan []byte) {
	log.Printf("[INFO] readBuffer start")
	for {
		select {
		case <-ctx.Done():
			log.Print("[INFO] stop readBuffer")
			return
		default:
			if b, err := p.in.Read(); err == nil && !(len(b) == 1 && b[0] == 0) {
				bufCh <- b
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (p *Prompt) setUp() {
	p.in.Setup()
	p.renderer.Setup()
	p.renderer.UpdateWinSize(p.in.GetWinSize())

	ctx, cancel := context.WithCancel(context.Background())
	p.ctx = ctx
	p.cancel = cancel
}

func (p *Prompt) tearDown() {
	p.cancel()
	p.in.TearDown()
	p.renderer.TearDown()
}
