package prompt

import (
	"bytes"
	"os"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/elk-language/go-prompt/internal/debug"
)

const inputBufferSize = 1024

// Executor is called when the user
// inputs a line of text.
type Executor func(string)

// ExitChecker is called after user input to check if prompt must stop and exit go-prompt Run loop.
// User input means: selecting/typing an entry, then, if said entry content matches the ExitChecker function criteria:
// - immediate exit (if breakline is false) without executor called
// - exit after typing <return> (meaning breakline is true), and the executor is called first, before exit.
// Exit means exit go-prompt (not the overall Go program)
type ExitChecker func(in string, breakline bool) bool

// Completer is a function that returns
// a slice of suggestions for the given Document.
type Completer func(Document) []Suggest

// Prompt is a core struct of go-prompt.
type Prompt struct {
	reader            Reader
	buf               *Buffer
	renderer          *Render
	executor          Executor
	history           *History
	lexer             Lexer
	completion        *CompletionManager
	keyBindings       []KeyBind
	ASCIICodeBindings []ASCIICodeBind
	keyBindMode       KeyBindMode
	completionOnDown  bool
	exitChecker       ExitChecker
	skipClose         bool
}

// UserInput is the struct that contains the user input context.
type UserInput struct {
	input string
}

// Run starts the prompt.
func (p *Prompt) Run() {
	p.skipClose = false
	defer debug.Close()
	debug.Log("start prompt")
	p.setup()
	defer p.Close()

	if p.completion.showAtStart {
		p.completion.Update(*p.buf.Document())
	}

	p.renderer.Render(p.buf, p.completion, p.lexer)

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
				p.renderer.BreakLine(p.buf, p.lexer)
				stopReadBufCh <- struct{}{}
				stopHandleSignalCh <- struct{}{}
				return
			} else if e != nil {
				// Stop goroutine to run readBuffer function
				stopReadBufCh <- struct{}{}
				stopHandleSignalCh <- struct{}{}

				// Unset raw mode
				// Reset to Blocking mode because returned EAGAIN when still set non-blocking mode.
				debug.AssertNoError(p.reader.Close())
				p.executor(e.input)

				p.completion.Update(*p.buf.Document())

				p.renderer.Render(p.buf, p.completion, p.lexer)

				if p.exitChecker != nil && p.exitChecker(e.input, true) {
					p.skipClose = true
					return
				}
				// Set raw mode
				debug.AssertNoError(p.reader.Open())
				go p.readBuffer(bufCh, stopReadBufCh)
				go p.handleSignals(exitCh, winSizeCh, stopHandleSignalCh)
			} else {
				p.completion.Update(*p.buf.Document())
				p.renderer.Render(p.buf, p.completion, p.lexer)
			}
		case w := <-winSizeCh:
			p.renderer.UpdateWinSize(w)
			p.renderer.Render(p.buf, p.completion, p.lexer)
		case code := <-exitCh:
			p.renderer.BreakLine(p.buf, p.lexer)
			p.Close()
			os.Exit(code)
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// func Log(format string, a ...any) {
// 	f, err := os.OpenFile("log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
// 	if err != nil {
// 		log.Fatalf("error opening file: %v", err)
// 	}
// 	defer f.Close()
// 	fmt.Fprintf(f, format, a...)
// }

func (p *Prompt) feed(b []byte) (shouldExit bool, userInput *UserInput) {
	key := GetKey(b)
	p.buf.lastKeyStroke = key
	// completion
	completing := p.completion.Completing()
	p.handleCompletionKeyBinding(key, completing)

	switch key {
	case Enter, ControlJ, ControlM:
		p.renderer.BreakLine(p.buf, p.lexer)

		userInput = &UserInput{input: p.buf.Text()}
		p.buf = NewBuffer()
		if userInput.input != "" {
			p.history.Add(userInput.input)
		}
	case ControlC:
		p.renderer.BreakLine(p.buf, p.lexer)
		p.buf = NewBuffer()
		p.history.Clear()
	case Up, ControlP:
		cursor := p.buf.Document().GetCursorPosition(int(p.renderer.col))
		if cursor.Y != 0 {
			p.buf.CursorUp(1)
			break
		}
		if completing {
			break
		}

		if newBuf, changed := p.history.Older(p.buf); changed {
			p.buf = newBuf
		}

	case Down, ControlN:
		endOfTextCursor := p.buf.Document().GetEndOfTextPosition(int(p.renderer.col))
		cursor := p.buf.Document().GetCursorPosition(int(p.renderer.col))
		if endOfTextCursor.Y > cursor.Y {
			p.buf.CursorDown(1)
			break
		}

		if completing {
			break
		}

		if newBuf, changed := p.history.Newer(p.buf); changed {
			p.buf = newBuf
		}
		return
	case ControlD:
		if p.buf.Text() == "" {
			shouldExit = true
			return
		}
	case NotDefined:
		if p.handleASCIICodeBinding(b) {
			return
		}
		char, _ := utf8.DecodeRune(b)
		if unicode.IsControl(char) {
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

// Input starts the prompt, lets the user
// input a single line and returns this line as a string.
func (p *Prompt) Input() string {
	defer debug.Close()
	debug.Log("start prompt")
	p.setup()
	defer p.Close()

	if p.completion.showAtStart {
		p.completion.Update(*p.buf.Document())
	}

	p.renderer.Render(p.buf, p.completion, p.lexer)
	bufCh := make(chan []byte, 128)
	stopReadBufCh := make(chan struct{})
	go p.readBuffer(bufCh, stopReadBufCh)

	for {
		select {
		case b := <-bufCh:
			if shouldExit, e := p.feed(b); shouldExit {
				p.renderer.BreakLine(p.buf, p.lexer)
				stopReadBufCh <- struct{}{}
				return ""
			} else if e != nil {
				// Stop goroutine to run readBuffer function
				stopReadBufCh <- struct{}{}
				return e.input
			} else {
				p.completion.Update(*p.buf.Document())
				p.renderer.Render(p.buf, p.completion, p.lexer)
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
			bytes := make([]byte, inputBufferSize)
			n, err := p.reader.Read(bytes)
			if err != nil {
				break
			}
			bytes = bytes[:n]
			if len(bytes) != 1 || bytes[0] != 0 {
				newBytes := make([]byte, len(bytes))
				for i, byt := range bytes {
					// translate raw mode \r into \n
					// to make pasting multiline text
					// work properly
					switch byt {
					case '\r':
						newBytes[i] = '\n'
					default:
						newBytes[i] = byt
					}
				}
				bufCh <- newBytes
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (p *Prompt) setup() {
	debug.AssertNoError(p.reader.Open())
	p.renderer.Setup()
	p.renderer.UpdateWinSize(p.reader.GetWinSize())
}

func (p *Prompt) Close() {
	if !p.skipClose {
		debug.AssertNoError(p.reader.Close())
	}
	p.renderer.Close()
}
