package prompt

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/c-bata/go-prompt/internal/debug"
	runewidth "github.com/mattn/go-runewidth"
)

// Render to render prompt information from state of Buffer.
type Render struct {
	out                ConsoleWriter
	prefix             string
	livePrefixCallback func() (prefix string, useLivePrefix bool)
	breakLineCallback  func(*Document)
	title              string
	row                uint16
	col                uint16

	previousCursor int

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
	if prefix, ok := r.livePrefixCallback(); ok {
		return prefix
	}
	return r.prefix
}

func (r *Render) renderPrefix() {
	r.out.SetColor(r.prefixTextColor, r.prefixBGColor, false)
	r.out.WriteStr(r.getCurrentPrefix())
	r.out.SetColor(DefaultColor, DefaultColor, false)
}

// TearDown to clear title and erasing.
func (r *Render) TearDown() {
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
	r.col = ws.Col
}

func (r *Render) renderWindowTooSmall() {
	r.out.CursorGoTo(0, 0)
	r.out.EraseScreen()
	r.out.SetColor(DarkRed, White, false)
	r.out.WriteStr("Your console window is too small...")
}

func (r *Render) renderCompletion(buf *Buffer, completions *CompletionManager) {
	suggestions := completions.GetSuggestions()
	if len(completions.GetSuggestions()) == 0 {
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

	cursor := runewidth.StringWidth(prefix) + runewidth.StringWidth(buf.Document().TextBeforeCursor())
	x, _ := r.toPos(cursor)
	if x+width >= int(r.col) {
		cursor = r.backward(cursor, x+width-int(r.col))
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
	r.out.SetColor(White, Cyan, false)
	for i := 0; i < windowHeight; i++ {
		r.out.CursorDown(1)
		if i == selected {
			r.out.SetColor(r.selectedSuggestionTextColor, r.selectedSuggestionBGColor, true)
		} else {
			r.out.SetColor(r.suggestionTextColor, r.suggestionBGColor, false)
		}
		r.out.WriteStr(formatted[i].Text)

		if i == selected {
			r.out.SetColor(r.selectedDescriptionTextColor, r.selectedDescriptionBGColor, false)
		} else {
			r.out.SetColor(r.descriptionTextColor, r.descriptionBGColor, false)
		}
		r.out.WriteStr(formatted[i].Description)

		if isScrollThumb(i) {
			r.out.SetColor(DefaultColor, r.scrollbarThumbColor, false)
		} else {
			r.out.SetColor(DefaultColor, r.scrollbarBGColor, false)
		}
		r.out.WriteStr(" ")
		r.out.SetColor(DefaultColor, DefaultColor, false)

		r.lineWrap(cursor + width)
		r.backward(cursor+width, width)
	}

	if x+width >= int(r.col) {
		r.out.CursorForward(x + width - int(r.col))
	}

	r.out.CursorUp(windowHeight)
	r.out.SetColor(DefaultColor, DefaultColor, false)
}

// ClearScreen :: Clears the screen and moves the cursor to home
func (r *Render) ClearScreen() {
	r.out.EraseScreen()
	r.out.CursorGoTo(0, 0)
}

// Render renders to the console.
func (r *Render) Render(buffer *Buffer, previousText string, completion *CompletionManager, lexer *Lexer) {
	// In situations where a pseudo tty is allocated (e.g. within a docker container),
	// window size via TIOCGWINSZ is not immediately available and will result in 0,0 dimensions.
	if r.col == 0 {
		return
	}
	defer func() { debug.AssertNoError(r.out.Flush()) }()

	line := buffer.Text()
	traceBackLines := strings.Count(previousText, "\n")
	if len(line) == 0 {
		// if the new buffer is empty, then we shouldn't traceback any
		traceBackLines = 0
	}
	debug.Log(fmt.Sprintln(line))
	debug.Log(fmt.Sprintln(traceBackLines))

	r.move((traceBackLines)*int(r.col)+r.previousCursor, 0)

	prefix := r.getCurrentPrefix()
	cursor := runewidth.StringWidth(prefix) + runewidth.StringWidth(line)

	// prepare area
	_, y := r.toPos((traceBackLines + int(r.col)) + cursor)

	h := y + 1 + int(completion.max)
	if h > int(r.row) || completionMargin > int(r.col) {
		r.renderWindowTooSmall()
		return
	}

	// Rendering
	r.out.HideCursor()

	r.out.EraseLine()
	r.out.EraseDown()
	r.renderPrefix()

	r.out.SetColor(DefaultColor, DefaultColor, false)

	r.lineWrap(cursor)

	if buffer.NewLineCount() > 0 {
		r.renderMultiline(buffer, lexer)
	} else {
		r.renderLine(line, lexer)
		defer r.out.ShowCursor()
	}

	r.lineWrap(cursor)
	r.out.SetColor(DefaultColor, DefaultColor, false)

	cursor = r.backward(cursor, runewidth.StringWidth(line)-buffer.DisplayCursorPosition())

	r.renderCompletion(buffer, completion)
	if suggest, ok := completion.GetSelectedSuggestion(); ok {
		cursor = r.backward(cursor, runewidth.StringWidth(buffer.Document().GetWordBeforeCursorUntilSeparator(completion.wordSeparator)))

		r.out.SetColor(r.previewSuggestionTextColor, r.previewSuggestionBGColor, false)
		r.out.WriteStr(suggest.Text)
		r.out.SetColor(DefaultColor, DefaultColor, false)
		cursor += runewidth.StringWidth(suggest.Text)

		rest := buffer.Document().TextAfterCursor()

		if lexer.IsEnabled {
			processed := lexer.Process(rest)

			var s = rest

			for _, v := range processed {
				a := strings.SplitAfter(s, v.Text)
				s = strings.TrimPrefix(s, a[0])

				r.out.SetColor(v.Color, r.inputBGColor, false)
				r.out.WriteStr(a[0])
			}
		} else {
			r.out.WriteStr(rest)
		}

		r.out.SetColor(DefaultColor, DefaultColor, false)

		cursor += runewidth.StringWidth(rest)
		r.lineWrap(cursor)

		cursor = r.backward(cursor, runewidth.StringWidth(rest))
	}
	r.previousCursor = cursor
}

func (r *Render) renderLine(line string, lexer *Lexer) {
	if lexer.IsEnabled {
		processed := lexer.Process(line)
		var s = line

		for _, v := range processed {
			a := strings.SplitAfter(s, v.Text)
			s = strings.TrimPrefix(s, a[0])

			r.out.SetColor(v.Color, r.inputBGColor, false)
			r.out.WriteStr(a[0])
		}
	} else {
		r.out.SetColor(r.inputTextColor, r.inputBGColor, false)
		r.out.WriteStr(line)
	}
}

func (r *Render) renderMultiline(buffer *Buffer, lexer *Lexer) {
	before := buffer.Document().TextBeforeCursor()
	cursor := ""
	after := ""

	if len(buffer.Document().TextAfterCursor()) == 0 {
		cursor = " "
		after = ""
	} else {
		cursor = string(buffer.Text()[buffer.Document().cursorPosition])
		if cursor == "\n" {
			cursor = " \n"
		}
		after = buffer.Document().TextAfterCursor()[1:]
	}

	r.out.SetColor(r.inputTextColor, r.inputBGColor, false)
	r.renderLine(before, lexer)

	r.out.SetDisplayAttributes(r.inputTextColor, r.inputBGColor, DisplayReverse)
	r.out.WriteRawStr(cursor)

	r.out.SetColor(r.inputTextColor, r.inputBGColor, false)
	r.renderLine(after, lexer)
}

// BreakLine to break line.
func (r *Render) BreakLine(buffer *Buffer, lexer *Lexer) {
	// Erasing and Render
	cursor := (buffer.NewLineCount() * int(r.col)) + runewidth.StringWidth(buffer.Document().TextBeforeCursor()) + runewidth.StringWidth(r.getCurrentPrefix())
	r.clear(cursor)

	r.renderPrefix()

	if lexer.IsEnabled {
		processed := lexer.Process(buffer.Document().Text + "\n")

		var s = buffer.Document().Text + "\n"

		for _, v := range processed {
			a := strings.SplitAfter(s, v.Text)
			s = strings.TrimPrefix(s, a[0])

			r.out.SetColor(v.Color, r.inputBGColor, false)
			r.out.WriteStr(a[0])
		}
	} else {
		r.out.SetColor(r.inputTextColor, r.inputBGColor, false)
		r.out.WriteStr(buffer.Document().Text + "\n")
	}

	r.out.SetColor(DefaultColor, DefaultColor, false)

	debug.AssertNoError(r.out.Flush())
	if r.breakLineCallback != nil {
		r.breakLineCallback(buffer.Document())
	}

	r.previousCursor = 0
}

// clear erases the screen from a beginning of input
// even if there is line break which means input length exceeds a window's width.
func (r *Render) clear(cursor int) {
	r.move(cursor, 0)
	r.out.EraseDown()
}

// backward moves cursor to backward from a current cursor position
// regardless there is a line break.
func (r *Render) backward(from, n int) int {
	return r.move(from, from-n)
}

// move moves cursor to specified position from the beginning of input
// even if there is a line break.
func (r *Render) move(from, to int) int {
	fromX, fromY := r.toPos(from)
	toX, toY := r.toPos(to)

	r.out.CursorUp(fromY - toY)
	r.out.CursorBackward(fromX - toX)
	return to
}

// toPos returns the relative position from the beginning of the string.
func (r *Render) toPos(cursor int) (x, y int) {
	col := int(r.col)
	return cursor % col, cursor / col
}

func (r *Render) lineWrap(cursor int) {
	if runtime.GOOS != "windows" && cursor > 0 && cursor%int(r.col) == 0 {
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
