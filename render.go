package prompt

import (
	"math"
)

// Render to render prompt information from state of Buffer.
type Render struct {
	out                ConsoleWriter
	prefix             string
	livePrefixCallback func() (prefix string, useLivePrefix bool)
	title              string
	row                uint16
	col                uint16

	previousBufferSize int

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
		r.out.Flush()
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
	r.out.Flush()
}

func (r *Render) prepareArea(lines int) {
	for i := 0; i < lines; i++ {
		r.out.ScrollDown()
	}
	for i := 0; i < lines; i++ {
		r.out.ScrollUp()
	}
	return
}

// UpdateWinSize called when window size is changed.
func (r *Render) UpdateWinSize(ws *WinSize) {
	r.row = ws.Row
	r.col = ws.Col
	return
}

func (r *Render) renderWindowTooSmall() {
	r.out.CursorGoTo(0, 0)
	r.out.EraseScreen()
	r.out.SetColor(DarkRed, White, false)
	r.out.WriteStr("Your console window is too small...")
	r.out.Flush()
	return
}

func (r *Render) renderCompletion(buf *Buffer, completions *CompletionManager) {
	windowHeight := len(completions.tmp)
	if windowHeight > int(completions.max) {
		windowHeight = int(completions.max)
	}
	contentHeight := len(completions.tmp)

	fractionVisible := float64(windowHeight) / float64(contentHeight)
	fractionAbove := float64(completions.verticalScroll) / float64(contentHeight)

	scrollbarHeight := int(math.Min(float64(windowHeight), math.Max(1, float64(windowHeight)*fractionVisible)))
	scrollbarTop := int(float64(windowHeight) * fractionAbove)

	isScrollThumb := func(row int) bool {
		return scrollbarTop <= row && row <= scrollbarTop+scrollbarHeight
	}

	suggestions := completions.GetSuggestions()
	if l := len(completions.GetSuggestions()); l == 0 {
		return
	}

	prefix := r.getCurrentPrefix()
	formatted, width := formatSuggestions(
		suggestions,
		int(r.col)-len(prefix)-1, // -1 means a width of scrollbar
	)
	formatted = formatted[completions.verticalScroll : completions.verticalScroll+windowHeight]
	l := len(formatted)
	r.prepareArea(windowHeight)

	// +1 means a width of scrollbar.
	d := (len(prefix) + len(buf.Document().TextBeforeCursor()) + 1) % int(r.col)
	if d == 0 { // the cursor is on right end.
		r.out.CursorBackward(width)
	} else if d+width > int(r.col) {
		r.out.CursorBackward(d + width - int(r.col))
	}

	selected := completions.selected - completions.verticalScroll

	r.out.SetColor(White, Cyan, false)
	for i := 0; i < l; i++ {
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
		// +1 means a width of scrollbar.
		r.out.CursorBackward(width + 1)
	}
	if d == 0 && len(prefix)+len(buf.Text()) != 0 { // the cursor is on right end.
		// DON'T CURSOR DOWN HERE. Because the line doesn't erase properly.
		r.out.CursorForward(width + 1)
	} else if d+width > int(r.col) {
		r.out.CursorForward(d + width - int(r.col))
	}

	r.out.CursorUp(l)
	r.out.SetColor(DefaultColor, DefaultColor, false)
	return
}

// Render renders to the console.
func (r *Render) Render(buffer *Buffer, completion *CompletionManager) {
	line := buffer.Text()
	prefix := r.getCurrentPrefix()

<<<<<<< HEAD
	// In situations where a psuedo tty is allocated (e.g. within a docker container),
	// window size via TIOCGWINSZ is not immediately available and will result in 0,0 dimensions.
	if r.col > 0 {
		if len(buffer.Document().Text) == 0 {
			r.out.CursorBackward(int(r.col))
			r.out.SaveCursor()
		}
		// Erasing
		r.out.CursorUp((r.previousBufferSize - 1) / int(r.col))
		r.out.WriteRawStr("\r")
		r.out.EraseDown()

		r.previousBufferSize = len(buffer.Document().TextBeforeCursor()) + len(prefix)

		// prepare area
		h := ((len(prefix) + len(line)) / int(r.col)) + 1 + int(completion.max)
		if h > int(r.row) || completionMargin > int(r.col) {
			r.renderWindowTooSmall()
			return
		}
	}

	// Rendering
	r.renderPrefix()

	r.out.SetColor(r.inputTextColor, r.inputBGColor, false)
	r.out.WriteStr(line)
	r.out.SetColor(DefaultColor, DefaultColor, false)
	r.out.CursorBackward(len([]rune(line)) - buffer.CursorPosition)
	r.renderCompletion(buffer, completion)
	if suggest, ok := completion.GetSelectedSuggestion(); ok {
		r.out.CursorBackward(len([]rune(buffer.Document().GetWordBeforeCursor())))
		r.out.SetColor(r.previewSuggestionTextColor, r.previewSuggestionBGColor, false)
		r.out.WriteStr(suggest.Text)
		r.out.SetColor(DefaultColor, DefaultColor, false)
	}
	r.out.Flush()
}

// BreakLine to break line.
func (r *Render) BreakLine(buffer *Buffer) {
	// Erasing and Render
	r.out.CursorUp((len(buffer.Document().Text) + len(r.getCurrentPrefix()) - 1) / int(r.col))
	r.out.WriteRawStr("\r")
	r.out.EraseDown()
	r.renderPrefix()
	r.out.SetColor(r.inputTextColor, r.inputBGColor, false)
	r.out.WriteStr(buffer.Document().Text + "\n")
	r.out.SetColor(DefaultColor, DefaultColor, false)
	r.out.Flush()

	r.previousBufferSize = 0
}
