package prompt

import (
	"bytes"
	"os"
	"strings"
	"syscall"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/elk-language/go-prompt/debug"
	istrings "github.com/elk-language/go-prompt/strings"
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

// ExecuteOnEnterCallback is a function that receives
// user input after Enter has been pressed
// and determines whether the input should be executed.
// If this function returns true, the Executor callback will be called
// otherwise a newline will be added to the buffer containing user input
// and optionally indentation made up of `indentSize * indent` spaces.
type ExecuteOnEnterCallback func(prompt *Prompt, indentSize int) (indent int, execute bool)

// Completer is a function that returns
// a slice of suggestions for the given Document.
//
// startChar and endChar represent the indices of the first and last rune of the text
// that the suggestions were generated for and that should be replaced by the selected suggestion.
type Completer func(Document) (suggestions []Suggest, startChar, endChar istrings.RuneNumber)

// Prompt is a core struct of go-prompt.
type Prompt struct {
	reader                 Reader
	buffer                 *Buffer
	renderer               *Renderer
	executor               Executor
	history                HistoryInterface
	lexer                  Lexer
	completion             *CompletionManager
	keyBindings            []KeyBind
	ASCIICodeBindings      []ASCIICodeBind
	keyBindMode            KeyBindMode
	exitChecker            ExitChecker
	executeOnEnterCallback ExecuteOnEnterCallback
	interruptCallback      InterruptCallback
	completionOnDown       bool
	skipClose              bool
	completionReset        bool
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
		p.completion.Update(*p.buffer.Document())
	}

	p.renderer.Render(p.buffer, p.completion, p.lexer)

	bufCh := make(chan []byte, 128)
	stopReadBufCh := make(chan struct{})
	go p.readBuffer(bufCh, stopReadBufCh)

	exitCh := make(chan os.Signal)
	winSizeCh := make(chan *WinSize)
	stopHandleSignalCh := make(chan struct{})
	go p.handleSignals(exitCh, winSizeCh, stopHandleSignalCh)

	for {
		select {
		case b := <-bufCh:
			if shouldExit, rerender, input := p.feed(b); shouldExit {
				p.renderer.BreakLine(p.buffer, p.lexer)
				stopReadBufCh <- struct{}{}
				stopHandleSignalCh <- struct{}{}
				return
			} else if input != nil {
				// Stop goroutine to run readBuffer function
				stopReadBufCh <- struct{}{}
				stopHandleSignalCh <- struct{}{}

				// Unset raw mode
				// Reset to Blocking mode because returned EAGAIN when still set non-blocking mode.
				debug.AssertNoError(p.reader.Close())
				p.executor(input.input)

				p.completion.Update(*p.buffer.Document())

				p.renderer.Render(p.buffer, p.completion, p.lexer)

				if p.exitChecker != nil && p.exitChecker(input.input, true) {
					p.skipClose = true
					return
				}
				// Set raw mode
				debug.AssertNoError(p.reader.Open())
				go p.readBuffer(bufCh, stopReadBufCh)
				go p.handleSignals(exitCh, winSizeCh, stopHandleSignalCh)
			} else if rerender {
				if p.completion.shouldUpdate {
					p.completion.Update(*p.buffer.Document())
				}
				p.renderer.Render(p.buffer, p.completion, p.lexer)
			}
		case w := <-winSizeCh:
			p.renderer.UpdateWinSize(w)
			p.buffer.resetStartLine()
			p.buffer.recalculateStartLine(p.renderer.UserInputColumns(), int(p.renderer.row))
			p.renderer.Render(p.buffer, p.completion, p.lexer)
		case signal := <-exitCh:
			p.renderer.BreakLine(p.buffer, p.lexer)
			if p.interruptCallback != nil {
				p.interruptCallback(p, signal)
				return
			}

			p.Close()
			switch signal {
			case syscall.SIGINT, syscall.SIGQUIT:
				os.Exit(0)
			case syscall.SIGTERM:
				os.Exit(1)
			}
		}
	}
}

// For debugging
// func log(format string, a ...any) {
// 	f, err := os.OpenFile("log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
// 	if err != nil {
// 		panic(fmt.Sprintf("error opening file: %v", err))
// 	}
// 	defer f.Close()
// 	if _, err := fmt.Fprintf(f, format+"\n", a...); err != nil {
// 		panic(err)
// 	}
// }

// Returns the configured indent size.
func (p *Prompt) IndentSize() int {
	return p.renderer.indentSize
}

// Get the number of columns that are available
// for user input.
func (p *Prompt) UserInputColumns() istrings.Width {
	return p.renderer.UserInputColumns()
}

// Returns the current amount of columns that the terminal can display.
func (p *Prompt) TerminalColumns() istrings.Width {
	return p.renderer.col
}

// Returns the current amount of rows that the terminal can display.
func (p *Prompt) TerminalRows() int {
	return p.renderer.row
}

// Returns the buffer struct.
func (p *Prompt) Buffer() *Buffer {
	return p.buffer
}

func (p *Prompt) feed(b []byte) (shouldExit bool, rerender bool, userInput *UserInput) {
	key := GetKey(b)
	p.buffer.lastKeyStroke = key
	// completion
	completing := p.completion.Completing()
	if p.handleCompletionKeyBinding(b, key, completing) {
		return false, true, nil
	}

	cols := p.renderer.UserInputColumns()
	rows := p.renderer.row

	switch key {
	case Enter, ControlJ, ControlM:
		indent, execute := p.executeOnEnterCallback(p, p.renderer.indentSize)
		if !execute {
			p.buffer.NewLine(cols, rows, false)

			var indentStrBuilder strings.Builder
			indentUnitCount := indent * p.renderer.indentSize
			for range indentUnitCount {
				indentStrBuilder.WriteRune(IndentUnit)
			}
			p.buffer.InsertTextMoveCursor(indentStrBuilder.String(), cols, rows, false)
			break
		}

		p.renderer.BreakLine(p.buffer, p.lexer)
		userInput = &UserInput{input: p.buffer.Text()}
		p.buffer = NewBuffer()
		if userInput.input != "" {
			p.history.Add(userInput.input)
		}
	case ControlC:
		p.renderer.BreakLine(p.buffer, p.lexer)
		p.buffer = NewBuffer()
		p.history.Clear()
	case Up, ControlP:
		line := p.buffer.Document().CursorPositionRow()
		if line > 0 {
			rerender = p.CursorUp(1)
			return false, rerender, nil
		}
		if completing {
			break
		}

		if newBuf, changed := p.history.Older(p.buffer, cols, rows); changed {
			p.buffer = newBuf
		}

	case Down, ControlN:
		endOfTextRow := p.buffer.Document().TextEndPositionRow()
		row := p.buffer.Document().CursorPositionRow()
		if endOfTextRow > row {
			rerender = p.CursorDown(1)
			return false, rerender, nil
		}

		if completing {
			break
		}

		if newBuf, changed := p.history.Newer(p.buffer, cols, rows); changed {
			p.buffer = newBuf
		}
		return false, true, nil
	case ControlD:
		if p.buffer.Text() == "" {
			return true, true, nil
		}
	case NotDefined:
		var checked bool
		checked, rerender = p.handleASCIICodeBinding(b, cols, rows)

		if checked {
			return false, rerender, nil
		}
		char, _ := utf8.DecodeRune(b)
		if unicode.IsControl(char) && !unicode.IsSpace(char) {
			return false, false, nil
		}

		p.buffer.InsertTextMoveCursor(string(b), cols, rows, false)
	}

	shouldExit, rerender = p.handleKeyBinding(key, cols, rows)
	return shouldExit, rerender, userInput
}

func (p *Prompt) handleCompletionKeyBinding(b []byte, key Key, completing bool) (handled bool) {
	p.completion.shouldUpdate = true
	cols := p.renderer.UserInputColumns()
	rows := p.renderer.row
	completionLen := len(p.completion.tmp)
	p.completionReset = false

keySwitch:
	switch key {
	case Down:
		if completing || p.completionOnDown {
			p.updateSuggestions(func() {
				p.completion.Next()
			})
			return true
		}
	case ControlI:
		p.updateSuggestions(func() {
			p.completion.Next()
		})
		return true
	case Up:
		if completing {
			p.updateSuggestions(func() {
				p.completion.Previous()
			})
			return true
		}
	case Tab:
		if completionLen > 0 {
			// If there are any suggestions, select the next one
			p.updateSuggestions(func() {
				p.completion.Next()
			})

			return true
		}

		// if there are no suggestions insert indentation
		newBytes := make([]byte, 0, len(b))
		for _, byt := range b {
			switch byt {
			case '\t':
				for i := 0; i < p.renderer.indentSize; i++ {
					newBytes = append(newBytes, IndentUnit)
				}
			default:
				newBytes = append(newBytes, byt)
			}
		}
		p.buffer.InsertTextMoveCursor(string(newBytes), cols, rows, false)
		return true
	case BackTab:
		if completionLen > 0 {
			// If there are any suggestions, select the previous one
			p.updateSuggestions(func() {
				p.completion.Previous()
			})
			return true
		}

		text := p.buffer.Document().CurrentLineBeforeCursor()
		for _, char := range text {
			if char != IndentUnit {
				break keySwitch
			}
		}
		p.buffer.DeleteBeforeCursorRunes(istrings.RuneNumber(p.renderer.indentSize), cols, rows)
		return true
	default:
		if s, ok := p.completion.GetSelectedSuggestion(); ok {
			w := p.buffer.Document().GetWordBeforeCursorUntilSeparator(p.completion.wordSeparator)
			if w != "" {
				p.buffer.DeleteBeforeCursorRunes(istrings.RuneCountInString(w), cols, rows)
			}
			p.buffer.InsertTextMoveCursor(s.Text, cols, rows, false)
		}
		if completionLen > 0 {
			p.completionReset = true
		}
		p.completion.Reset()
	}
	return false
}

func (p *Prompt) updateSuggestions(fn func()) {
	cols := p.renderer.UserInputColumns()
	rows := p.renderer.row

	prevStart := p.completion.startCharIndex
	prevEnd := p.completion.endCharIndex
	prevSuggestion, prevSelected := p.completion.GetSelectedSuggestion()

	fn()

	p.completion.shouldUpdate = false
	newSuggestion, newSelected := p.completion.GetSelectedSuggestion()

	// do nothing
	if !prevSelected && !newSelected {
		return
	}

	// insert the new selection
	if !prevSelected {
		p.buffer.DeleteBeforeCursorRunes(p.completion.endCharIndex-p.completion.startCharIndex, cols, rows)
		p.buffer.InsertTextMoveCursor(newSuggestion.Text, cols, rows, false)
		return
	}
	// delete the previous selection
	if !newSelected {
		p.buffer.DeleteBeforeCursorRunes(
			istrings.RuneCountInString(prevSuggestion.Text)-(prevEnd-prevStart),
			cols,
			rows,
		)
		return
	}

	// delete previous selection and render the new one
	p.buffer.DeleteBeforeCursorRunes(
		istrings.RuneCountInString(prevSuggestion.Text),
		cols,
		rows,
	)

	p.buffer.InsertTextMoveCursor(newSuggestion.Text, cols, rows, false)
}

func (p *Prompt) handleKeyBinding(key Key, cols istrings.Width, rows int) (shouldExit bool, rerender bool) {
	var executed bool

	switch p.keyBindMode {
	case EmacsKeyBind:
		p.handleKeyBind(commonKeyBindings, key, &executed, &rerender)
		p.handleKeyBind(emacsKeyBindings, key, &executed, &rerender)
		p.handleKeyBind(p.keyBindings, key, &executed, &rerender)
	case CommonKeyBind:
		p.handleKeyBind(commonKeyBindings, key, &executed, &rerender)
		p.handleKeyBind(p.keyBindings, key, &executed, &rerender)
	case CustomOverrideCommonKeyBind:
		p.handleKeyBind(p.keyBindings, key, &executed, &rerender)
		if !executed {
			p.handleKeyBind(commonKeyBindings, key, &executed, &rerender)
		}
	case CustomOverrideEmacsKeyBind:
		p.handleKeyBind(p.keyBindings, key, &executed, &rerender)
		if !executed {
			p.handleKeyBind(emacsKeyBindings, key, &executed, &rerender)
		}
		if !executed {
			p.handleKeyBind(commonKeyBindings, key, &executed, &rerender)
		}
	}

	// Custom key bindings
	if p.exitChecker != nil && p.exitChecker(p.buffer.Text(), false) {
		shouldExit = true
	}
	if !executed && !rerender {
		rerender = true
	}
	return shouldExit, rerender
}

func (p *Prompt) handleKeyBind(keyBinds []KeyBind, key Key, executed, rerender *bool) {
	for _, kb := range keyBinds {
		if kb.Key == key {
			result := kb.Fn(p)
			*executed = true
			if !*rerender {
				*rerender = result
			}
		}
	}
}

func (p *Prompt) handleASCIICodeBinding(b []byte, cols istrings.Width, rows int) (checked, rerender bool) {
	for _, kb := range p.ASCIICodeBindings {
		if bytes.Equal(kb.ASCIICode, b) {
			result := kb.Fn(p)
			if !rerender {
				rerender = result
			}
			checked = true
		}
	}
	return checked, rerender
}

// Input starts the prompt, lets the user
// input a single line and returns this line as a string.
func (p *Prompt) Input() string {
	defer debug.Close()
	debug.Log("start prompt")
	p.setup()
	defer p.Close()

	if p.completion.showAtStart {
		p.completion.Update(*p.buffer.Document())
	}

	p.renderer.Render(p.buffer, p.completion, p.lexer)
	bufCh := make(chan []byte, 128)
	stopReadBufCh := make(chan struct{})
	go p.readBuffer(bufCh, stopReadBufCh)

	for {
		select {
		case b := <-bufCh:
			if shouldExit, rerender, input := p.feed(b); shouldExit {
				p.renderer.BreakLine(p.buffer, p.lexer)
				stopReadBufCh <- struct{}{}
				return ""
			} else if input != nil {
				// Stop goroutine to run readBuffer function
				stopReadBufCh <- struct{}{}
				return input.input
			} else if rerender {
				p.completion.Update(*p.buffer.Document())
				p.renderer.Render(p.buffer, p.completion, p.lexer)
			}
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

const IndentUnit = ' '
const IndentUnitString = string(IndentUnit)

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
			// Log("%#v", bytes)
			if len(bytes) == 1 && bytes[0] == '\t' {
				// if only a single Tab key has been pressed
				// handle it as a keybind
				bufCh <- bytes
			} else if len(bytes) != 1 || bytes[0] != 0 {
				newBytes := make([]byte, 0, len(bytes))
				for _, byt := range bytes {
					switch byt {
					// translate raw mode \r into \n
					// to make pasting multiline text
					// work properly
					case '\r':
						newBytes = append(newBytes, '\n')
					// translate \t into two spaces
					// to avoid problems with cursor positions
					case '\t':
						for range p.renderer.indentSize {
							newBytes = append(newBytes, IndentUnit)
						}
					default:
						newBytes = append(newBytes, byt)
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

// Move to the left on the current line  by the given amount of graphemes (visible characters).
// Returns true when the view should be rerendered.
func (p *Prompt) CursorLeft(count istrings.GraphemeNumber) bool {
	return promptCursorHorizontalMove(p, p.buffer.CursorLeft, count)
}

// Move to the left on the current line by the given amount of runes.
// Returns true when the view should be rerendered.
func (p *Prompt) CursorLeftRunes(count istrings.RuneNumber) bool {
	return promptCursorHorizontalMove(p, p.buffer.CursorLeftRunes, count)
}

// Move the cursor to the right on the current line by the given amount of graphemes (visible characters).
// Returns true when the view should be rerendered.
func (p *Prompt) CursorRight(count istrings.GraphemeNumber) bool {
	return promptCursorHorizontalMove(p, p.buffer.CursorRight, count)
}

// Move the cursor to the right on the current line by the given amount of runes.
// Returns true when the view should be rerendered.
func (p *Prompt) CursorRightRunes(count istrings.RuneNumber) bool {
	return promptCursorHorizontalMove(p, p.buffer.CursorRightRunes, count)
}

type horizontalCursorModifier[CountT ~int] func(CountT, istrings.Width, int) bool

// Move to the left or right on the current line.
// Returns true when the view should be rerendered.
func promptCursorHorizontalMove[CountT ~int](p *Prompt, modifierFunc horizontalCursorModifier[CountT], count CountT) bool {
	b := p.buffer
	cols := p.renderer.UserInputColumns()
	previousCursor := b.DisplayCursorPosition(cols)

	rerender := modifierFunc(count, cols, p.renderer.row) || p.completionReset || len(p.completion.tmp) > 0
	if rerender {
		return true
	}

	newCursor := b.DisplayCursorPosition(cols)
	p.renderer.previousCursor = newCursor
	p.renderer.move(previousCursor, newCursor)
	p.renderer.flush()
	return false
}

// Move the cursor up.
// Returns true when the view should be rerendered.
func (p *Prompt) CursorUp(count int) bool {
	b := p.buffer
	cols := p.renderer.UserInputColumns()
	previousCursor := b.DisplayCursorPosition(cols)

	rerender := p.buffer.CursorUp(count, cols, p.renderer.row) || p.completionReset || len(p.completion.tmp) > 0
	if rerender {
		return true
	}

	newCursor := b.DisplayCursorPosition(cols)
	p.renderer.previousCursor = newCursor
	p.renderer.move(previousCursor, newCursor)
	p.renderer.flush()
	return false
}

// Move the cursor down.
// Returns true when the view should be rerendered.
func (p *Prompt) CursorDown(count int) bool {
	b := p.buffer
	cols := p.renderer.UserInputColumns()
	previousCursor := b.DisplayCursorPosition(cols)

	rerender := p.buffer.CursorDown(count, cols, p.renderer.row) || p.completionReset || len(p.completion.tmp) > 0
	if rerender {
		return true
	}

	newCursor := b.DisplayCursorPosition(cols)
	p.renderer.previousCursor = newCursor
	p.renderer.move(previousCursor, newCursor)
	p.renderer.flush()
	return false
}

// Deletes the specified number of graphemes before the cursor and returns the deleted text.
func (p *Prompt) DeleteBeforeCursor(count istrings.GraphemeNumber) string {
	return p.buffer.DeleteBeforeCursor(count, p.UserInputColumns(), p.renderer.row)
}

// Deletes the specified number of runes before the cursor and returns the deleted text.
func (p *Prompt) DeleteBeforeCursorRunes(count istrings.RuneNumber) string {
	return p.buffer.DeleteBeforeCursorRunes(count, p.UserInputColumns(), p.renderer.row)
}

// Deletes the specified number of graphemes and returns the deleted text.
func (p *Prompt) Delete(count istrings.GraphemeNumber) string {
	return p.buffer.Delete(count, p.UserInputColumns(), p.renderer.row)
}

// Deletes the specified number of runes and returns the deleted text.
func (p *Prompt) DeleteRunes(count istrings.RuneNumber) string {
	return p.buffer.DeleteRunes(count, p.UserInputColumns(), p.renderer.row)
}

// Insert string into the buffer without moving the cursor.
func (p *Prompt) InsertText(text string, overwrite bool) {
	p.buffer.InsertText(text, overwrite)
}

// Insert string into the buffer and move the cursor.
func (p *Prompt) InsertTextMoveCursor(text string, overwrite bool) {
	p.buffer.InsertTextMoveCursor(text, p.UserInputColumns(), p.renderer.row, overwrite)
}

func (p *Prompt) History() HistoryInterface {
	return p.history
}

func (p *Prompt) Close() {
	if !p.skipClose {
		debug.AssertNoError(p.reader.Close())
	}
	p.renderer.Close()
}
