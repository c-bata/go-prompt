package prompt

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/daichi-m/go-prompt/internal/debug"
	"github.com/fatih/color"
	fcolor "github.com/fatih/color"
	runewidth "github.com/mattn/go-runewidth"
)

// Render to render prompt information from state of Buffer.
type Render struct {
	out                ConsoleWriter
	prefix             string
	livePrefixCallback func() (prefix string, ok bool)
	breakLineCallback  func(*Document)
	statusBarCallback  func(*Buffer, *CompletionManager) (status string, ok bool)
	title              string
	row                uint16
	col                uint16
	previousCursor     int
	keywords           map[string]bool

	// colors,
	prefixColor              *fcolor.Color
	inputColor               *fcolor.Color
	keywordColor             *fcolor.Color
	previewSuggestionColor   *fcolor.Color
	suggestionColor          *fcolor.Color
	selectedSuggestionColor  *fcolor.Color
	descriptionColor         *fcolor.Color
	selectedDescriptionColor *fcolor.Color
	statusBarColor           *fcolor.Color
	scrollbarColor           *fcolor.Color
	scrollbarThumbColor      *fcolor.Color
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
	r.out.WriteStr(r.getCurrentPrefix(), r.prefixColor)
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
	r.out.WriteStr("Your console window is too small...", color.New(color.FgHiRed, color.BgWhite))
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
	x, y := r.toPos(cursor)
	if x+width >= int(r.col) {
		cursor = r.backward(cursor, x+width-int(r.col))
	}
	debug.Log(fmt.Sprintf("Cursor position for render completion: %d,%d\n", x, y))

	contentHeight := len(completions.tmp)

	fractionVisible := float64(windowHeight) / float64(contentHeight)
	fractionAbove := float64(completions.verticalScroll) / float64(contentHeight)

	scrollbarHeight := int(clamp(float64(windowHeight)*fractionVisible, float64(windowHeight), 1))
	scrollbarTop := int(float64(windowHeight) * fractionAbove)

	isScrollThumb := func(row int) bool {
		return scrollbarTop <= row && row <= scrollbarTop+scrollbarHeight
	}

	selected := completions.selected - completions.verticalScroll
	// r.out.SetColor(White, Cyan, false)
	for i := 0; i < windowHeight; i++ {
		r.out.CursorDown(1)
		var color *fcolor.Color
		if i == selected {
			color = r.selectedSuggestionColor.Add(fcolor.Bold)
		} else {
			color = r.selectedDescriptionColor
		}
		r.out.WriteStr(formatted[i].Text, color)

		if i == selected {
			color = r.selectedDescriptionColor
		} else {
			color = r.descriptionColor
		}
		r.out.WriteStr(formatted[i].Description, color)

		if isScrollThumb(i) {
			color = r.scrollbarThumbColor
		} else {
			color = r.scrollbarColor
		}
		r.out.WriteStr(" ", color)
		// r.out.SetColor(DefaultColor, DefaultColor, false)

		r.lineWrap(cursor + width)
		r.backward(cursor+width, width)
	}

	if x+width >= int(r.col) {
		r.out.CursorForward(x + width - int(r.col))
	}

	r.out.CursorUp(windowHeight)
}

func (r *Render) renderSelectSuggestionInBuffer(buffer *Buffer, suggest Suggest, cursorStart int) (cursor int) {

	cursor = cursorStart
	// r.out.SetColor(r.previewSuggestionTextColor, r.previewSuggestionBGColor, false)
	r.out.WriteStr(suggest.Text, r.previewSuggestionColor)
	// r.out.SetColor(DefaultColor, DefaultColor, false)
	cursor += runewidth.StringWidth(suggest.Text)

	rest := buffer.Document().TextAfterCursor()
	r.out.WriteStr(rest, nil)
	cursor += runewidth.StringWidth(rest)
	r.lineWrap(cursor)

	cursor = r.backward(cursor, runewidth.StringWidth(rest))
	return
}

func (r *Render) renderStatusBar(buffer *Buffer, completion *CompletionManager) {

	r.out.SaveCursor()
	defer func() {
		r.out.UnSaveCursor()
		r.out.CursorUp(0)
	}()

	if status, ok := r.statusBarCallback(buffer, completion); ok {
		r.out.CursorDown(int(r.row))
		r.out.CursorBackward(int(r.col))
		fs, _ := formatTexts([]string{status}, int(r.col-2), "", "")
		if len(fs) == 0 {
			return
		}
		fs = padTexts(fs, " ", int(r.col))
		r.out.WriteStr(fs[0], r.statusBarColor)
	}

}

func padTexts(orig []string, pad string, length int) []string {

	pl := len(pad)
	if pl <= 0 {
		return orig
	}
	if len(orig) == 0 {
		tot, mod := length/pl, length%pl
		return []string{strings.Repeat(pad, tot) + pad[0:mod]}
	}
	padded := make([]string, 0, len(orig))

	for _, o := range orig {
		fillLen := length - len(o)
		tot, mod := fillLen/pl, fillLen%pl
		p := strings.Repeat(pad, tot) + pad[0:mod]
		padded = append(padded, o+p)
	}
	return padded

}

func (r *Render) renderLine(line string) {
	words := strings.Split(line, " ")
	l := len(words)
	for i, w := range words {
		lw := strings.ToLower(w)
		if _, ok := r.keywords[lw]; ok {
			r.out.WriteStr(w, r.keywordColor)
		} else if len(w) > 0 {
			r.out.WriteStr(w, r.inputColor)
		}
		if i < l-1 {
			r.out.WriteStr(" ", r.inputColor)
		}
	}
}

// Render renders to the console.
func (r *Render) Render(buffer *Buffer, completion *CompletionManager) {
	// In situations where a pseudo tty is allocated (e.g. within a docker container),
	// window size via TIOCGWINSZ is not immediately available and will result in 0,0 dimensions.
	debug.Log(fmt.Sprintf("Window size in which to render: (%d x %d)", r.row, r.col))
	if r.col == 0 {
		return
	}
	defer func() { debug.AssertNoError(r.out.Flush()) }()
	r.prepareArea(2)
	r.move(r.previousCursor, 0)

	line := buffer.Text()
	prefix := r.getCurrentPrefix()
	cursor := runewidth.StringWidth(prefix) + runewidth.StringWidth(line)

	// prepare area
	_, y := r.toPos(cursor)

	h := y + 1 + int(completion.max)
	if h > int(r.row) || completionMargin > int(r.col) {
		r.renderWindowTooSmall()
		return
	}

	// Rendering
	r.out.HideCursor()
	defer r.out.ShowCursor()

	r.renderPrefix()
	r.renderLine(line)
	r.lineWrap(cursor)

	r.out.EraseDown()

	cursor = r.backward(cursor, runewidth.StringWidth(line)-buffer.DisplayCursorPosition())

	r.renderCompletion(buffer, completion)
	r.renderStatusBar(buffer, completion)
	if suggest, ok := completion.GetSelectedSuggestion(); ok {
		cursor = r.backward(cursor, runewidth.StringWidth(
			buffer.Document().GetWordBeforeCursorUntilSeparator(completion.wordSeparator)))
		cursor = r.renderSelectSuggestionInBuffer(buffer, suggest, cursor)
	}
	r.previousCursor = cursor
}

// BreakLine to break line.
func (r *Render) BreakLine(buffer *Buffer) {
	// Erasing and Render
	cursor := runewidth.StringWidth(buffer.Document().TextBeforeCursor()) + runewidth.StringWidth(r.getCurrentPrefix())
	r.clear(cursor)
	r.renderPrefix()
	r.renderLine(buffer.Document().Text)
	r.out.WriteStr("\n", r.inputColor)
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
		r.out.WriteStr("\n", nil)
	}
}

// clamp ensures that x is within the limit of low and high.
func clamp(x, high, low float64) float64 {
	switch {
	case high < x:
		return high
	case x < low:
		return low
	default:
		return x
	}
}
