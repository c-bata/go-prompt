package prompt

// Render to render prompt information from state of Buffer.
type Render struct {
	out                ConsoleWriter
	prefix             string
	livePrefixCallback func() (prefix string, useLivePrefix bool)
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
	suggestions := completions.GetSuggestions()
	if len(completions.GetSuggestions()) == 0 {
		return
	}
	prefix := r.getCurrentPrefix()
	formatted, width := formatSuggestions(
		suggestions,
		int(r.col)-len(prefix)-1, // -1 means a width of scrollbar
	)
	// +1 means a width of scrollbar.
	width++

	windowHeight := len(formatted)
	if windowHeight > int(completions.max) {
		windowHeight = int(completions.max)
	}
	formatted = formatted[completions.verticalScroll : completions.verticalScroll+windowHeight]
	r.prepareArea(windowHeight)

	cursor := len(prefix) + len(buf.Document().TextBeforeCursor())
	x, _ := r.toPos(cursor)
	if x+width >= int(r.col) {
		r.out.CursorBackward(x + width - int(r.col))
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
		r.out.CursorBackward(width)
	}

	if x+width >= int(r.col) {
		r.out.CursorForward(x + width - int(r.col))
	}

	r.out.CursorUp(windowHeight)
	r.out.SetColor(DefaultColor, DefaultColor, false)
	return
}

// Render renders to the console.
func (r *Render) Render(buffer *Buffer, completion *CompletionManager) {
	line := buffer.Text()
	prefix := r.getCurrentPrefix()
	cursor := len(prefix) + len(line)

	// In situations where a psuedo tty is allocated (e.g. within a docker container),
	// window size via TIOCGWINSZ is not immediately available and will result in 0,0 dimensions.
	if r.col > 0 {
		// Erasing
		r.clear(r.previousCursor)

		// prepare area
		_, y := r.toPos(cursor)

		h := y + 1 + int(completion.max)
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

	cursor = r.backward(cursor, len(line)-buffer.CursorPosition)

	r.renderCompletion(buffer, completion)
	if suggest, ok := completion.GetSelectedSuggestion(); ok {
		cursor = r.backward(cursor, len(buffer.Document().GetWordBeforeCursor()))

		r.out.SetColor(r.previewSuggestionTextColor, r.previewSuggestionBGColor, false)
		r.out.WriteStr(suggest.Text)
		r.out.SetColor(DefaultColor, DefaultColor, false)

		cursor += len(suggest.Text)
	}
	r.out.Flush()

	r.previousCursor = cursor
}

// BreakLine to break line.
func (r *Render) BreakLine(buffer *Buffer) {
	// Erasing and Render
	cursor := len(buffer.Document().TextBeforeCursor()) + len(r.getCurrentPrefix())
	r.clear(cursor)
	r.renderPrefix()
	r.out.SetColor(r.inputTextColor, r.inputBGColor, false)
	r.out.WriteStr(buffer.Document().Text + "\n")
	r.out.SetColor(DefaultColor, DefaultColor, false)
	r.out.Flush()

	r.previousCursor = 0
}

func (r *Render) clear(cursor int) {
	r.backward(cursor, cursor)
	r.out.EraseDown()
}

func (r *Render) backward(from, n int) int {
	return r.move(from, from-n)
}

func (r *Render) move(from, to int) int {
	_, fromY := r.toPos(from)
	toX, toY := r.toPos(to)

	r.out.CursorUp(fromY - toY)
	r.out.WriteRaw([]byte{'\r'})
	r.out.CursorForward(toX)
	return to
}

// toPos returns the relative position from the beginning of the string.
// the coordinate system with the beginning of the string as (0,0) and the width as r.col.
// the cursor points to the next character, but it points to that character only at the right end (x == r.col - 1).
// x will not return 0 except for the first row.
func (r *Render) toPos(cursor int) (x, y int) {
	col := int(r.col)

	if cursor > 0 && cursor%col == 0 {
		return col - 1, cursor/col - 1
	}

	return cursor % col, cursor / col
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
