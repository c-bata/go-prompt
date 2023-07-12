package prompt

import (
	"runtime"
	"strings"

	"github.com/elk-language/go-prompt/debug"
	istrings "github.com/elk-language/go-prompt/strings"
	runewidth "github.com/mattn/go-runewidth"
)

// Render to render prompt information from state of Buffer.
type Render struct {
	out               Writer
	prefixCallback    PrefixCallback
	breakLineCallback func(*Document)
	title             string
	row               uint16
	col               istrings.StringWidth

	previousCursor Position

	// colors,
	prefixTextColor              Color
	prefixBGColor                Color
	inputTextColor               Color
	inputBGColor                 Color
	previewSuggestionTextColor   Color
	previewSuggestionBGColor     Color
	suggestionTextColor          Color
	suggestionBGColor            Color
	selectedSuggestionTextColor  Color
	selectedSuggestionBGColor    Color
	descriptionTextColor         Color
	descriptionBGColor           Color
	selectedDescriptionTextColor Color
	selectedDescriptionBGColor   Color
	scrollbarThumbColor          Color
	scrollbarBGColor             Color
}

// Setup to initialize console output.
func (r *Render) Setup() {
	if r.title != "" {
		r.out.SetTitle(r.title)
		debug.AssertNoError(r.out.Flush())
	}
}

// getCurrentPrefix to get current prefix.
// If live-prefix is enabled, return live-prefix.
func (r *Render) getCurrentPrefix() string {
	return r.prefixCallback()
}

func (r *Render) renderPrefix() {
	r.out.SetColor(r.prefixTextColor, r.prefixBGColor, false)
	r.out.WriteString("\r")
	r.out.WriteString(r.getCurrentPrefix())
	r.out.SetColor(DefaultColor, DefaultColor, false)
}

// Close to clear title and erasing.
func (r *Render) Close() {
	r.out.ClearTitle()
	r.out.EraseDown()
	debug.AssertNoError(r.out.Flush())
}

func (r *Render) prepareArea(lines int) {
	for i := 0; i < lines; i++ {
		r.out.ScrollDown()
	}
	for i := 0; i < lines; i++ {
		r.out.ScrollUp()
	}
}

// UpdateWinSize called when window size is changed.
func (r *Render) UpdateWinSize(ws *WinSize) {
	r.row = ws.Row
	r.col = istrings.StringWidth(ws.Col)
}

func (r *Render) renderWindowTooSmall() {
	r.out.CursorGoTo(0, 0)
	r.out.EraseScreen()
	r.out.SetColor(DarkRed, White, false)
	r.out.WriteString("Your console window is too small...")
}

func (r *Render) renderCompletion(buf *Buffer, completions *CompletionManager) {
	suggestions := completions.GetSuggestions()
	if len(suggestions) == 0 {
		return
	}
	prefix := r.getCurrentPrefix()
	formatted, width := formatSuggestions(
		suggestions,
		int(r.col)-runewidth.StringWidth(prefix)-1, // -1 means a width of scrollbar
	)
	// +1 means a width of scrollbar.
	width++

	windowHeight := len(formatted)
	if windowHeight > int(completions.max) {
		windowHeight = int(completions.max)
	}
	formatted = formatted[completions.verticalScroll : completions.verticalScroll+windowHeight]
	r.prepareArea(windowHeight)

	cursor := positionAtEndOfString(prefix+buf.Document().TextBeforeCursor(), r.col)
	x := cursor.X
	if x+width >= r.col {
		cursor = r.backward(cursor, x+width-r.col)
	}

	contentHeight := len(completions.tmp)

	fractionVisible := float64(windowHeight) / float64(contentHeight)
	fractionAbove := float64(completions.verticalScroll) / float64(contentHeight)

	scrollbarHeight := int(clamp(float64(windowHeight), 1, float64(windowHeight)*fractionVisible))
	scrollbarTop := int(float64(windowHeight) * fractionAbove)

	isScrollThumb := func(row int) bool {
		return scrollbarTop <= row && row <= scrollbarTop+scrollbarHeight
	}

	selected := completions.selected - completions.verticalScroll
	cursorColumnSpacing := cursor

	r.out.SetColor(White, Cyan, false)
	for i := 0; i < windowHeight; i++ {
		alignNextLine(r, cursorColumnSpacing.X)

		if i == selected {
			r.out.SetColor(r.selectedSuggestionTextColor, r.selectedSuggestionBGColor, true)
		} else {
			r.out.SetColor(r.suggestionTextColor, r.suggestionBGColor, false)
		}
		r.out.WriteString(formatted[i].Text)

		if i == selected {
			r.out.SetColor(r.selectedDescriptionTextColor, r.selectedDescriptionBGColor, false)
		} else {
			r.out.SetColor(r.descriptionTextColor, r.descriptionBGColor, false)
		}
		r.out.WriteString(formatted[i].Description)

		if isScrollThumb(i) {
			r.out.SetColor(DefaultColor, r.scrollbarThumbColor, false)
		} else {
			r.out.SetColor(DefaultColor, r.scrollbarBGColor, false)
		}
		r.out.WriteString(" ")
		r.out.SetColor(DefaultColor, DefaultColor, false)

		c := cursor.Add(Position{X: width})
		r.lineWrap(&c)
		r.backward(c, width)
	}

	if x+width >= r.col {
		r.out.CursorForward(int(x + width - r.col))
	}

	r.out.CursorUp(windowHeight)
	r.out.SetColor(DefaultColor, DefaultColor, false)
}

// Render renders to the console.
func (r *Render) Render(buffer *Buffer, completion *CompletionManager, lexer Lexer) {
	// In situations where a pseudo tty is allocated (e.g. within a docker container),
	// window size via TIOCGWINSZ is not immediately available and will result in 0,0 dimensions.
	if r.col == 0 {
		return
	}
	defer func() { debug.AssertNoError(r.out.Flush()) }()
	r.clear(r.previousCursor)

	line := buffer.Text()
	prefix := r.getCurrentPrefix()
	prefixWidth := istrings.StringWidth(runewidth.StringWidth(prefix))
	cursor := positionAtEndOfString(prefix+line, r.col)

	// prepare area
	y := cursor.Y

	h := y + 1 + int(completion.max)
	if h > int(r.row) || completionMargin > int(r.col) {
		r.renderWindowTooSmall()
		return
	}

	// Rendering
	r.out.HideCursor()
	defer r.out.ShowCursor()

	r.renderPrefix()

	if lexer != nil {
		r.lex(lexer, line)
	} else {
		r.out.SetColor(r.inputTextColor, r.inputBGColor, false)
		r.out.WriteString(line)
	}

	r.out.SetColor(DefaultColor, DefaultColor, false)

	r.lineWrap(&cursor)

	targetCursor := buffer.DisplayCursorPosition(r.col)
	if targetCursor.Y == 0 {
		targetCursor.X += prefixWidth
	}
	cursor = r.move(cursor, targetCursor)

	r.renderCompletion(buffer, completion)
	if suggest, ok := completion.GetSelectedSuggestion(); ok {
		cursor = r.backward(cursor, istrings.StringWidth(runewidth.StringWidth(buffer.Document().GetWordBeforeCursorUntilSeparator(completion.wordSeparator))))

		r.out.SetColor(r.previewSuggestionTextColor, r.previewSuggestionBGColor, false)
		r.out.WriteString(suggest.Text)
		r.out.SetColor(DefaultColor, DefaultColor, false)
		cursor.X += istrings.StringWidth(runewidth.StringWidth(suggest.Text))
		endOfSuggestionPos := cursor

		rest := buffer.Document().TextAfterCursor()

		if lexer != nil {
			r.lex(lexer, rest)
		} else {
			r.out.WriteString(rest)
		}

		r.out.SetColor(DefaultColor, DefaultColor, false)

		cursor = cursor.Join(positionAtEndOfString(rest, r.col))

		r.lineWrap(&cursor)

		cursor = r.move(cursor, endOfSuggestionPos)
	}
	r.previousCursor = cursor
}

// lex processes the given input with the given lexer
// and writes the result
func (r *Render) lex(lexer Lexer, input string) {
	lexer.Init(input)
	s := input

	for {
		token, ok := lexer.Next()
		if !ok {
			break
		}

		a := strings.SplitAfter(s, token.Lexeme())
		s = strings.TrimPrefix(s, a[0])

		r.out.SetColor(token.Color(), r.inputBGColor, false)
		r.out.WriteString(a[0])
	}
}

// BreakLine to break line.
func (r *Render) BreakLine(buffer *Buffer, lexer Lexer) {
	// Erasing and Render
	cursor := positionAtEndOfString(buffer.Document().TextBeforeCursor()+r.getCurrentPrefix(), r.col)
	r.clear(cursor)

	r.renderPrefix()

	if lexer != nil {
		r.lex(lexer, buffer.Document().Text+"\n")
	} else {
		r.out.SetColor(r.inputTextColor, r.inputBGColor, false)
		r.out.WriteString(buffer.Document().Text + "\n")
	}

	r.out.SetColor(DefaultColor, DefaultColor, false)

	debug.AssertNoError(r.out.Flush())
	if r.breakLineCallback != nil {
		r.breakLineCallback(buffer.Document())
	}

	r.previousCursor = Position{}
}

// clear erases the screen from a beginning of input
// even if there is line break which means input length exceeds a window's width.
func (r *Render) clear(cursor Position) {
	r.move(cursor, Position{})
	r.out.EraseDown()
}

// backward moves cursor to backward from a current cursor position
// regardless there is a line break.
func (r *Render) backward(from Position, n istrings.StringWidth) Position {
	return r.move(from, Position{X: from.X - n, Y: from.Y})
}

// move moves cursor to specified position from the beginning of input
// even if there is a line break.
func (r *Render) move(from, to Position) Position {
	newPosition := from.Subtract(to)
	r.out.CursorUp(newPosition.Y)
	r.out.CursorBackward(int(newPosition.X))
	return to
}

func (r *Render) lineWrap(cursor *Position) {
	if runtime.GOOS != "windows" && cursor.X > 0 && cursor.X%r.col == 0 {
		cursor.X = 0
		cursor.Y += 1
		r.out.WriteRaw([]byte{'\n'})
	}
}

func clamp(high, low, x float64) float64 {
	switch {
	case high < x:
		return high
	case x < low:
		return low
	default:
		return x
	}
}

func alignNextLine(r *Render, col istrings.StringWidth) {
	r.out.CursorDown(1)
	r.out.WriteString("\r")
	r.out.CursorForward(int(col))
}
