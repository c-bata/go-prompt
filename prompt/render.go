package prompt

import "strings"

type Render struct {
	out            ConsoleWriter
	prefix         string
	title          string
	row            uint16
	col            uint16
	// colors
	prefixTextColor             Color
	prefixBGColor               Color
	inputTextColor              Color
	inputBGColor                Color
	outputTextColor             Color
	outputBGColor               Color
	previewSuggestionTextColor  Color
	previewSuggestionBGColor    Color
	suggestionTextColor         Color
	suggestionBGColor           Color
	selectedSuggestionTextColor Color
	selectedSuggestionBGColor   Color
}

func (r *Render) Setup() {
	if r.title != "" {
		r.out.SetTitle(r.title)
		r.out.Flush()
	}
}

func (r *Render) renderPrefix() {
	r.out.SetColor(r.prefixTextColor, r.prefixBGColor)
	r.out.WriteStr(r.prefix)
	r.out.SetColor(DefaultColor, DefaultColor)
}

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

func (r *Render) UpdateWinSize(ws *WinSize) {
	r.row = ws.Row
	r.col = ws.Col
	return
}

func (r *Render) renderWindowTooSmall() {
	r.out.CursorGoTo(0, 0)
	r.out.EraseScreen()
	r.out.SetColor(DarkRed, White)
	r.out.WriteStr("Your console window is too small...")
	r.out.Flush()
	return
}

func (r *Render) renderCompletion(buf *Buffer, words []string, max uint16, selected int) {
	if max > r.row {
		max = r.row
	}

	if l := len(words); l == 0 {
		return
	} else if l > int(max) {
		words = words[:max]
	}

	formatted, width := formatCompletions(
		words,
		int(r.col) - len(r.prefix),
		" ",
		" ",
	)
	l := len(formatted)

	d := (len(r.prefix) + len(buf.Document().TextBeforeCursor())) % int(r.col)
	if d + width > int(r.col) {
		r.out.CursorBackward(d + width - int(r.col))
	}

	r.out.SetColor(White, Cyan)
	for i := 0; i < l; i++ {
		r.out.CursorDown(1)
		if i == selected {
			r.out.SetColor(r.selectedSuggestionTextColor, r.selectedSuggestionBGColor)
		} else {
			r.out.SetColor(r.suggestionTextColor, r.suggestionBGColor)
		}
		r.out.WriteStr(formatted[i])
		r.out.CursorBackward(width)
	}
	if d + width > int(r.col) {
		r.out.CursorForward(d + width - int(r.col))
	}

	r.out.CursorUp(l)
	r.out.SetColor(DefaultColor, DefaultColor)
	return
}

func (r *Render) Render(buffer *Buffer, completions []string, maxCompletions uint16, selected int) {
	// Erasing
	r.out.CursorBackward(int(r.col) + len(buffer.Text()) + len(r.prefix))
	r.out.EraseDown()

	// prepare area
	line := buffer.Text()
	h := ((len(r.prefix) + len(line)) / int(r.col)) + 1 + int(maxCompletions)
	if h > int(r.row) {
		r.renderWindowTooSmall()
		return
	}
	r.prepareArea(h)

	// Rendering
	r.renderPrefix()

	r.out.SetColor(r.inputTextColor, r.inputBGColor)
	r.out.WriteStr(line)
	r.out.SetColor(DefaultColor, DefaultColor)
	r.out.CursorBackward(len(line) - buffer.CursorPosition)
	r.renderCompletion(buffer, completions, maxCompletions, selected)
	if selected != -1 {
		c := completions[selected]
		r.out.CursorBackward(len([]rune(buffer.Document().GetWordBeforeCursor())))
		r.out.SetColor(r.previewSuggestionTextColor, r.previewSuggestionBGColor)
		r.out.WriteStr(c)
		r.out.SetColor(DefaultColor, DefaultColor)
	}
	r.out.Flush()
}

func (r *Render) BreakLine(buffer *Buffer, result string) {
	// Erasing
	r.out.CursorBackward(int(r.col) + len(buffer.Text()) + len(r.prefix))
	r.out.EraseDown()
	r.renderPrefix()

	// Render Line Break
	r.out.SetColor(r.inputTextColor, r.inputBGColor)
	r.out.WriteStr(buffer.Document().Text + "\n")

	// Render Result
	if result != "" {
		r.out.SetColor(r.outputTextColor, r.outputBGColor)
		r.out.WriteStr(result + "\n")
	}
	r.out.SetColor(DefaultColor, DefaultColor)
}

func formatCompletions(words []string, max int, prefix string, suffix string) (new []string, width int) {
	num := len(words)
	new = make([]string, num)
	width = 0

	for i := 0; i < num; i++ {
		if width < len([]rune(words[i])) {
			width = len([]rune(words[i]))
		}
	}

	if len(prefix) + width + len(suffix) > max {
		width = max - len(prefix) - len(suffix)
	}

	for i := 0; i < num; i++ {
		if l := len(words[i]); l > width {
			new[i] = prefix + words[i][:width - len("...")] + "..." + suffix
		} else if l < width  {
			spaces := strings.Repeat(" ", width - len([]rune(words[i])))
			new[i] = prefix + words[i] + spaces + suffix
		} else {
			new[i] = prefix + words[i] + suffix
		}
	}
	width += len(prefix) + len(suffix)
	return
}
