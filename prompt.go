package prompt

import (
	"bytes"
	"os"
	"strconv"
	"time"
	"unicode/utf8"

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
	editBuf           *Buffer
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
	histSearch        bool
	histSearchFwd     bool
}

// Exec is the struct contains user input context.
type Exec struct {
	input string
}

// Run starts prompt.
func (p *Prompt) Run() {
	p.skipTearDown = false
	defer debug.Teardown()
	debug.Log("start prompt")
	p.freshBuffer(false)
	p.setUp()
	defer p.tearDown()

	if p.completion.showAtStart {
		p.completion.Update(*p.editBuf.Document())
	}

	p.renderer.Render(p.buf, p.completion)

	bufCh := make(chan []byte, 128)
	stopReadBufCh := make(chan struct{})
	go p.readBuffer(bufCh, stopReadBufCh)

	exitCh := make(chan int)
	winSizeCh := make(chan *WinSize)
	stopHandleSignalCh := make(chan struct{})
	go p.handleSignals(exitCh, winSizeCh, stopHandleSignalCh)

	for {
		select {
		case b := <-bufCh:
			if shouldExit, e := p.feed(b); shouldExit {
				p.renderer.BreakLine(p.buf)
				stopReadBufCh <- struct{}{}
				stopHandleSignalCh <- struct{}{}
				return
			} else if e != nil {
				// Stop goroutine to run readBuffer function
				stopReadBufCh <- struct{}{}
				stopHandleSignalCh <- struct{}{}

				// Unset raw mode
				// Reset to Blocking mode because returned EAGAIN when still set non-blocking mode.
				debug.AssertNoError(p.in.TearDown())
				p.executor(e.input)

				p.completion.Update(*p.editBuf.Document())

				p.renderer.Render(p.buf, p.completion)

				if p.exitChecker != nil && p.exitChecker(e.input, true) {
					p.skipTearDown = true
					return
				}
				// Set raw mode
				debug.AssertNoError(p.in.Setup())
				go p.readBuffer(bufCh, stopReadBufCh)
				go p.handleSignals(exitCh, winSizeCh, stopHandleSignalCh)
			} else {
				p.completion.Update(*p.editBuf.Document())
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

func (p *Prompt) SetHistory(hist *History) {
	p.history = hist
}

func (p *Prompt) freshBuffer(split bool) {
	p.buf = NewBuffer()
	p.histSearch = split
	p.completion.Enable(!split)
	if split {
		p.editBuf = NewBuffer()
	} else {
		p.editBuf = p.buf
	}
}

func (p *Prompt) updateHistSearch(finalize bool) {
	if p.histSearch {
		editText := p.editBuf.Text()
		found := p.history.Search(editText, p.histSearchFwd, false)
		p.buf = NewBuffer()
		if !finalize {
			p.buf.InsertText(editText+": ", false, true)
		}
		if found == "" {
			if finalize {
				p.buf.InsertText(editText, false, true)
			}
		} else {
			p.buf.InsertText(found, false, true)
		}
		//p.completion.Update(*NewBuffer().Document())
		p.renderer.Render(p.buf, p.completion)
	}
}

func (p *Prompt) feed(b []byte) (shouldExit bool, exec *Exec) {
	key := GetKey(b)
	p.editBuf.lastKeyStroke = key
	// completion
	completing := p.completion.Completing() && !p.histSearch
	p.handleCompletionKeyBinding(key, completing)

	switch key {
	case Enter, ControlJ, ControlM:
		if  p.histSearch {
			p.updateHistSearch(true)
		}
		p.renderer.BreakLine(p.buf)

		exec = &Exec{input: p.buf.Text()}
		p.freshBuffer(false)
		if exec.input != "" {
			p.history.Add(exec.input)
		}
	case ControlC:
		p.histSearch = false
		p.renderer.BreakLine(p.buf)
		p.freshBuffer(false)
		p.history.Clear()
		p.histSearch = false
	case Up, ControlP:
		if p.histSearch {
			p.histSearch = false
			p.freshBuffer(false)
		} else if !completing { // Don't use p.completion.Completing() because it takes double operation when switch to selected=-1.
			if newBuf, changed := p.history.Older(p.editBuf); changed {
				p.buf = newBuf
				p.editBuf = newBuf
			}
		}
	case Down, ControlN:
		if p.histSearch {
			p.histSearch = false
			p.freshBuffer(false)
		} else if !completing { // Don't use p.completion.Completing() because it takes double operation when switch to selected=-1.
			if newBuf, changed := p.history.Newer(p.editBuf); changed {
				p.buf = newBuf
				p.editBuf = newBuf
			}
			return
		}
	case Left, Right:
		if p.histSearch {
			p.updateHistSearch(true)
			p.editBuf = p.buf
			p.histSearch = false
			//p.renderer.BreakLine(p.buf)
			//p.completion.Reset()
			p.completion.Enable(true)
			p.completion.Update(*p.buf.Document())
			p.renderer.Render(p.buf, p.completion)
		}
	case ControlR:
		p.histSearchFwd = false
		if p.histSearch {
			p.history.Search(p.editBuf.Text(), p.histSearchFwd, true)
			p.updateHistSearch(false)
		} else {
			p.histSearch = true
			p.completion.Reset()
			p.completion.Enable(false)
			p.buf = NewBuffer()
			p.history.SearchReset(false)
			p.updateHistSearch(false)
		}
		return
	case ControlS:
		p.histSearchFwd = true
		if p.histSearch {
			p.history.Search(p.editBuf.Text(), p.histSearchFwd, true)
			p.updateHistSearch(false)
			return
		}
	case Escape:
		if p.histSearch {
			p.histSearch = false
			p.freshBuffer(false)
			return
		}
	case ControlD:
		if p.editBuf.Text() == "" {
			shouldExit = true
			return
		}
	case NotDefined:
		if p.handleASCIICodeBinding(b) {
			return
		}
		// check for unprintable characters (e.g. multiple simultaneous cursor-keys)
		roon, _ := utf8.DecodeRune(b)
		if strconv.IsPrint(roon) {
			p.editBuf.InsertText(string(b), false, true)
		} else {
			p.histSearch = false
			p.buf = p.editBuf
		}
	}

	shouldExit = p.handleKeyBinding(key)
	p.updateHistSearch(false)
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
			w := p.editBuf.Document().GetWordBeforeCursorUntilSeparator(p.completion.wordSeparator)
			if w != "" {
				p.editBuf.DeleteBeforeCursor(len([]rune(w)))
			}
			p.editBuf.InsertText(s.Text, false, true)
		}
		p.completion.Reset()
	}
}

func (p *Prompt) handleKeyBinding(key Key) bool {
	shouldExit := false
	for i := range commonKeyBindings {
		kb := commonKeyBindings[i]
		if kb.Key == key {
			kb.Fn(p.editBuf)
		}
	}

	if p.keyBindMode == EmacsKeyBind {
		for i := range emacsKeyBindings {
			kb := emacsKeyBindings[i]
			if kb.Key == key {
				kb.Fn(p.editBuf)
			}
		}
	}

	// Custom key bindings
	for i := range p.keyBindings {
		kb := p.keyBindings[i]
		if kb.Key == key {
			kb.Fn(p.editBuf)
		}
	}
	if p.exitChecker != nil && p.exitChecker(p.editBuf.Text(), false) {
		shouldExit = true
	}
	return shouldExit
}

func (p *Prompt) handleASCIICodeBinding(b []byte) bool {
	checked := false
	for _, kb := range p.ASCIICodeBindings {
		if bytes.Equal(kb.ASCIICode, b) {
			kb.Fn(p.editBuf)
			checked = true
		}
	}
	return checked
}

// Input just returns user input text.
func (p *Prompt) Input() string {
	defer debug.Teardown()
	debug.Log("start prompt")
	p.setUp()
	defer p.tearDown()

	if p.completion.showAtStart {
		p.completion.Update(*p.editBuf.Document())
	}

	p.renderer.Render(p.buf, p.completion)
	bufCh := make(chan []byte, 128)
	stopReadBufCh := make(chan struct{})
	go p.readBuffer(bufCh, stopReadBufCh)

	for {
		select {
		case b := <-bufCh:
			if shouldExit, e := p.feed(b); shouldExit {
				p.renderer.BreakLine(p.buf)
				stopReadBufCh <- struct{}{}
				return ""
			} else if e != nil {
				// Stop goroutine to run readBuffer function
				stopReadBufCh <- struct{}{}
				return e.input
			} else {
				p.completion.Update(*p.editBuf.Document())
				p.renderer.Render(p.buf, p.completion)
			}
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func (p *Prompt) readBuffer(bufCh chan []byte, stopCh chan struct{}) {
	debug.Log("start reading buffer")
	for {
		select {
		case <-stopCh:
			debug.Log("stop reading buffer")
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
