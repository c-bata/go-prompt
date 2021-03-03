package prompt

import (
	"bytes"
	"os"

	"github.com/c-bata/go-prompt/internal/debug"
)

// Executor is called when user input something text.
type Executor func(string)

// ExitChecker is called after user input to check if prompt must stop and exit go-prompt Run loop.
// User input means: selecting/typing an entry, then, if said entry content matches the ExitChecker function criteria:
// - immediate exit (if breakline is false) without executor called
// - exit after typing <return> (meaning breakline is true), and the executor is called first, before exit.
// Exit means exit go-prompt (not the overall Go program)
type ExitChecker func(in string, breakline bool) bool

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
	completionOnDown  bool
	exitChecker       ExitChecker
	skipTearDown      bool
}

// Exec is the struct contains user input context.
type Exec struct {
	input string
}

// Run starts prompt.
func (p *Prompt) Run() {
	p.skipTearDown = false
	defer debug.Teardown()
	debug.Log("Run start prompt")
	p.setUp()
	defer p.tearDown()

	if p.completion.showAtStart {
		p.completion.Update(*p.buf.Document())
	}

	p.renderer.Render(p.buf, p.completion)

	bufCh := make(chan []byte)
	go p.readBuffer(bufCh)

	exitCh := make(chan int)
	winSizeCh := make(chan *WinSize)
	stopHandleSignalCh := make(chan struct{})
	go p.handleSignals(exitCh, winSizeCh, stopHandleSignalCh)

	for {
		select {
		case b := <-bufCh:
			if shouldExit, e := p.feed(b); shouldExit {
				p.renderer.BreakLine(p.buf)
				stopHandleSignalCh <- struct{}{}
				return
			} else if e != nil {
				stopHandleSignalCh <- struct{}{}

				// Unset raw mode
				// Reset to Blocking mode because returned EAGAIN when still set non-blocking mode.
				debug.AssertNoError(p.in.TearDown())
				p.executor(e.input)

				p.completion.Update(*p.buf.Document())

				p.renderer.Render(p.buf, p.completion)

				if p.exitChecker != nil && p.exitChecker(e.input, true) {
					p.skipTearDown = true
					return
				}
				// Set raw mode
				debug.AssertNoError(p.in.Setup())
				go p.handleSignals(exitCh, winSizeCh, stopHandleSignalCh)
			} else {
				p.completion.Update(*p.buf.Document())
				p.renderer.Render(p.buf, p.completion)
			}
		case w := <-winSizeCh:
			p.renderer.UpdateWinSize(w)
			p.renderer.Render(p.buf, p.completion)
		case code := <-exitCh:
			p.renderer.BreakLine(p.buf)
			p.tearDown()
			os.Exit(code)
		}
	}
}

func (p *Prompt) feed(b []byte) (shouldExit bool, exec *Exec) {
	key := GetKey(b)
	p.buf.lastKeyStroke = key
	// completion
	completing := p.completion.Completing()
	p.handleCompletionKeyBinding(key, completing)

	switch key {
	case Enter, ControlJ, ControlM:
		p.renderer.BreakLine(p.buf)

		exec = &Exec{input: p.buf.Text()}
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

	shouldExit = p.handleKeyBinding(key)
	return
}

func (p *Prompt) handleCompletionKeyBinding(key Key, completing bool) {
	switch key {
	case Down:
		if completing || p.completionOnDown {
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

func (p *Prompt) handleKeyBinding(key Key) bool {
	shouldExit := false
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
	if p.exitChecker != nil && p.exitChecker(p.buf.Text(), false) {
		shouldExit = true
	}
	return shouldExit
}

func (p *Prompt) handleASCIICodeBinding(b []byte) bool {
	checked := false
	for _, kb := range p.ASCIICodeBindings {
		if bytes.Equal(kb.ASCIICode, b) {
			kb.Fn(p.buf)
			checked = true
		}
	}
	return checked
}

// Input just returns user input text.
func (p *Prompt) Input() string {
	defer debug.Teardown()
	debug.Log("Input start prompt")
	p.setUp()
	defer p.tearDown()

	if p.completion.showAtStart {
		p.completion.Update(*p.buf.Document())
	}

	p.renderer.Render(p.buf, p.completion)
	bufCh := make(chan []byte)
	go p.readBuffer(bufCh)

	for {
		select {
		case b := <-bufCh:
			if shouldExit, e := p.feed(b); shouldExit {
				p.renderer.BreakLine(p.buf)
				return ""
			} else if e != nil {
				return e.input
			} else {
				p.completion.Update(*p.buf.Document())
				p.renderer.Render(p.buf, p.completion)
			}
		}
	}
}

func (p *Prompt) readBuffer(bufCh chan []byte) {
	debug.Log("start reading buffer")

	// blocking read from tty, this groutine will not exit until os.Exit()
	for {
		if b, err := p.in.Read(); err == nil && !(len(b) == 1 && b[0] == 0) {
			bufCh <- b // blocking send
		}
	}
}

func (p *Prompt) setUp() {
	debug.AssertNoError(p.in.Setup())
	p.renderer.Setup()
	p.renderer.UpdateWinSize(p.in.GetWinSize())
}

func (p *Prompt) tearDown() {
	if !p.skipTearDown {
		debug.AssertNoError(p.in.TearDown())
	}
	p.renderer.TearDown()
}
