package prompt

import (
	"strings"
	"unicode/utf8"

	"github.com/elk-language/go-prompt/debug"
	istrings "github.com/elk-language/go-prompt/strings"
)

const multilinePrefixCharacter = '.'

// Renderer to render prompt information from state of Buffer.
type Renderer struct {
	out               Writer
	prefixCallback    PrefixCallback
	breakLineCallback func(*Document)
	title             string
	row               uint16
	col               istrings.Width
	indentSize        int // How many spaces constitute a single indentation level

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

// Build a new Renderer.
func NewRenderer() *Renderer {
	defaultWriter := NewStdoutWriter()
	registerWriter(defaultWriter)

	return &Renderer{
		out:                          defaultWriter,
		indentSize:                   DefaultIndentSize,
		prefixCallback:               DefaultPrefixCallback,
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
	}
}

// Setup to initialize console output.
func (r *Renderer) Setup() {
	if r.title != "" {
		r.out.SetTitle(r.title)
		debug.AssertNoError(r.out.Flush())
	}
}

func (r *Renderer) renderPrefix(prefix string) {
	r.out.SetColor(r.prefixTextColor, r.prefixBGColor, false)
	if _, err := r.out.WriteString("\r"); err != nil {
		panic(err)
	}
	if _, err := r.out.WriteString(prefix); err != nil {
		panic(err)
	}
	r.out.SetColor(DefaultColor, DefaultColor, false)
}

// Close to clear title and erasing.
func (r *Renderer) Close() {
	r.out.ClearTitle()
	r.out.EraseDown()
	debug.AssertNoError(r.out.Flush())
}

func (r *Renderer) prepareArea(lines int) {
	for i := 0; i < lines; i++ {
		r.out.ScrollDown()
	}
	for i := 0; i < lines; i++ {
		r.out.ScrollUp()
	}
}

// UpdateWinSize called when window size is changed.
func (r *Renderer) UpdateWinSize(ws *WinSize) {
	r.row = ws.Row
	r.col = istrings.Width(ws.Col)
}

func (r *Renderer) renderWindowTooSmall() {
	r.out.CursorGoTo(0, 0)
	r.out.EraseScreen()
	r.out.SetColor(DarkRed, White, false)
	if _, err := r.out.WriteString("Your console window is too small..."); err != nil {
		panic(err)
	}
}

func (r *Renderer) renderCompletion(buf *Buffer, completions *CompletionManager) {
	suggestions := completions.GetSuggestions()
	if len(suggestions) == 0 {
		return
	}
	prefix := r.prefixCallback()
	prefixWidth := istrings.GetWidth(prefix)
	formatted, width := formatSuggestions(
		suggestions,
		r.col-istrings.GetWidth(prefix)-1, // -1 means a width of scrollbar
	)
	// +1 means a width of scrollbar.
	width++

	windowHeight := len(formatted)
	if windowHeight > int(completions.max) {
		windowHeight = int(completions.max)
	}
	formatted = formatted[completions.verticalScroll : completions.verticalScroll+windowHeight]
	r.prepareArea(windowHeight)

	cursor := positionAtEndOfString(prefix+buf.Document().TextBeforeCursor(), r.col-prefixWidth)
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
		if _, err := r.out.WriteString(formatted[i].Text); err != nil {
			panic(err)
		}

		if i == selected {
			r.out.SetColor(r.selectedDescriptionTextColor, r.selectedDescriptionBGColor, false)
		} else {
			r.out.SetColor(r.descriptionTextColor, r.descriptionBGColor, false)
		}
		if _, err := r.out.WriteString(formatted[i].Description); err != nil {
			panic(err)
		}

		if isScrollThumb(i) {
			r.out.SetColor(DefaultColor, r.scrollbarThumbColor, false)
		} else {
			r.out.SetColor(DefaultColor, r.scrollbarBGColor, false)
		}
		if _, err := r.out.WriteString(" "); err != nil {
			panic(err)
		}
		r.out.SetColor(DefaultColor, DefaultColor, false)

		c := cursor.Add(Position{X: width})
		r.backward(c, width)
	}

	if x+width >= r.col {
		r.out.CursorForward(int(x + width - r.col))
	}

	r.out.CursorUp(windowHeight)
	r.out.SetColor(DefaultColor, DefaultColor, false)
}

// Render renders to the console.
func (r *Renderer) Render(buffer *Buffer, completion *CompletionManager, lexer Lexer) {
	// In situations where a pseudo tty is allocated (e.g. within a docker container),
	// window size via TIOCGWINSZ is not immediately available and will result in 0,0 dimensions.
	if r.col == 0 {
		return
	}
	defer func() { debug.AssertNoError(r.out.Flush()) }()
	r.clear(r.previousCursor)

	text := buffer.Text()
	prefix := r.prefixCallback()
	prefixWidth := istrings.GetWidth(prefix)
	cursor := positionAtEndOfString(text, r.col-prefixWidth)
	cursor.X += prefixWidth

	// prepare area
	y := cursor.Y

	h := y + 1 + int(completion.max)
	if h > int(r.row) || completionMargin > r.col {
		r.renderWindowTooSmall()
		return
	}

	// Rendering
	r.out.HideCursor()
	defer r.out.ShowCursor()

	r.renderText(lexer, text)

	r.out.SetColor(DefaultColor, DefaultColor, false)

	targetCursor := buffer.DisplayCursorPosition(r.col - prefixWidth)
	targetCursor.X += prefixWidth
	// Log("col: %#v, targetCursor: %#v, cursor: %#v\n", r.col-prefixWidth, targetCursor, cursor)
	cursor = r.move(cursor, targetCursor)

	r.renderCompletion(buffer, completion)
	if suggest, ok := completion.GetSelectedSuggestion(); ok {
		cursor = r.backward(cursor, istrings.GetWidth(buffer.Document().GetWordBeforeCursorUntilSeparator(completion.wordSeparator)))

		r.out.SetColor(r.previewSuggestionTextColor, r.previewSuggestionBGColor, false)
		if _, err := r.out.WriteString(suggest.Text); err != nil {
			panic(err)
		}
		r.out.SetColor(DefaultColor, DefaultColor, false)
		cursor.X += istrings.GetWidth(suggest.Text)
		endOfSuggestionPos := cursor

		rest := buffer.Document().TextAfterCursor()

		r.renderText(lexer, text)

		r.out.SetColor(DefaultColor, DefaultColor, false)

		cursor = cursor.Join(positionAtEndOfString(rest, r.col-prefixWidth))

		cursor = r.move(cursor, endOfSuggestionPos)
	}
	r.previousCursor = cursor
}

func (r *Renderer) renderText(lexer Lexer, text string) {
	if lexer != nil {
		r.lex(lexer, text)
		return
	}

	prefix := r.prefixCallback()
	prefixWidth := istrings.GetWidth(prefix)
	col := r.col - prefixWidth
	multilinePrefix := r.getMultilinePrefix(prefix)
	firstIteration := true
	var lineBuffer strings.Builder
	var lineCharIndex istrings.Width

	for _, char := range text {
		if lineCharIndex >= col || char == '\n' {
			lineBuffer.WriteRune('\n')
			r.renderLine(prefix, lineBuffer.String(), r.inputTextColor)
			lineCharIndex = 0
			lineBuffer.Reset()
			if char != '\n' {
				lineBuffer.WriteRune(char)
				lineCharIndex += istrings.GetRuneWidth(char)
			}
			if firstIteration {
				prefix = multilinePrefix
				firstIteration = false
			}
			continue
		}

		lineBuffer.WriteRune(char)
		lineCharIndex += istrings.GetRuneWidth(char)
	}

	r.renderLine(prefix, lineBuffer.String(), r.inputTextColor)
}

func (r *Renderer) renderLine(prefix, line string, color Color) {
	r.renderPrefix(prefix)
	r.writeString(line, color)
}

func (r *Renderer) writeString(text string, color Color) {
	r.out.SetColor(color, r.inputBGColor, false)
	if _, err := r.out.WriteString(text); err != nil {
		panic(err)
	}
}

func (r *Renderer) getMultilinePrefix(prefix string) string {
	var spaceCount int
	var dotCount int
	var nonSpaceCharSeen bool
	for {
		if len(prefix) == 0 {
			break
		}
		char, size := utf8.DecodeLastRuneInString(prefix)
		prefix = prefix[:len(prefix)-size]
		charWidth := istrings.GetRuneWidth(char)
		if nonSpaceCharSeen {
			dotCount += int(charWidth)
			continue
		}
		if char != ' ' {
			nonSpaceCharSeen = true
			dotCount += int(charWidth)
			continue
		}
		spaceCount += int(charWidth)
	}

	var multilinePrefixBuilder strings.Builder

	for i := 0; i < dotCount; i++ {
		multilinePrefixBuilder.WriteByte(multilinePrefixCharacter)
	}
	for i := 0; i < spaceCount; i++ {
		multilinePrefixBuilder.WriteByte(IndentUnit)
	}

	return multilinePrefixBuilder.String()
}

// lex processes the given input with the given lexer
// and writes the result
func (r *Renderer) lex(lexer Lexer, input string) {
	lexer.Init(input)
	s := input

	prefix := r.prefixCallback()
	prefixWidth := istrings.GetWidth(prefix)
	col := r.col - prefixWidth
	multilinePrefix := r.getMultilinePrefix(prefix)
	r.renderPrefix(prefix)
	var lineCharIndex istrings.Width
	for {
		token, ok := lexer.Next()
		if !ok {
			break
		}

		text := strings.SplitAfter(s, token.Lexeme())[0]
		s = strings.TrimPrefix(s, text)

		var lineBuffer strings.Builder

		for _, char := range text {
			if lineCharIndex >= col || char == '\n' {
				if char != '\n' {
					lineBuffer.WriteByte('\n')
				}
				r.writeString(lineBuffer.String(), token.Color())
				r.renderPrefix(multilinePrefix)
				lineCharIndex = 0
				lineBuffer.Reset()
				if char != '\n' {
					lineBuffer.WriteRune(char)
					lineCharIndex += istrings.GetRuneWidth(char)
				}
				continue
			}

			lineBuffer.WriteRune(char)
			lineCharIndex += istrings.GetRuneWidth(char)
		}
		r.writeString(lineBuffer.String(), token.Color())
	}
}

// BreakLine to break line.
func (r *Renderer) BreakLine(buffer *Buffer, lexer Lexer) {
	// Erasing and Renderer
	prefix := r.prefixCallback()
	prefixWidth := istrings.GetWidth(prefix)
	cursor := positionAtEndOfString(buffer.Document().TextBeforeCursor(), r.col-prefixWidth)
	cursor.X += prefixWidth
	r.clear(cursor)

	text := buffer.Document().Text
	r.renderText(lexer, text)
	if _, err := r.out.WriteString("\n"); err != nil {
		panic(err)
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
func (r *Renderer) clear(cursor Position) {
	r.move(cursor, Position{})
	r.out.EraseDown()
}

// backward moves cursor to backward from a current cursor position
// regardless there is a line break.
func (r *Renderer) backward(from Position, n istrings.Width) Position {
	return r.move(from, Position{X: from.X - n, Y: from.Y})
}

// move moves cursor to specified position from the beginning of input
// even if there is a line break.
func (r *Renderer) move(from, to Position) Position {
	newPosition := from.Subtract(to)
	r.out.CursorUp(newPosition.Y)
	r.out.CursorBackward(int(newPosition.X))
	return to
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

func alignNextLine(r *Renderer, col istrings.Width) {
	r.out.CursorDown(1)
	if _, err := r.out.WriteString("\r"); err != nil {
		panic(err)
	}
	r.out.CursorForward(int(col))
}
