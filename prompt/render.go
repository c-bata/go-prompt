package prompt

import "strings"

const scrollBarWidth = 1

type Render struct {
	out            ConsoleWriter
	prefix         string
	title          string
	row            uint16
	col            uint16
	maxCompletions uint16
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
	}
	r.renderPrefix()
	r.out.Flush()
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

func (r *Render) renderCompletion(buf *Buffer, words []string, chosen int) {
	max := int(r.maxCompletions)
	if r.maxCompletions > r.row {
		max = int(r.row)
	}

	if l := len(words); l == 0 {
		return
	} else if l > max {
		words = words[:max]
	}

	formatted, width := formatCompletions(
		words,
		int(r.col) - len(r.prefix) - scrollBarWidth,
		" ",
		" ",
	)
	l := len(formatted)
	r.prepareArea(l)

	d := (len(r.prefix) + len(buf.Document().TextBeforeCursor())) % int(r.col)
	if d + width + scrollBarWidth > int(r.col) {
		r.out.CursorBackward(d + width + 1 - int(r.col))
	}

	r.out.SetColor(White, Cyan)
	for i := 0; i < l; i++ {
		r.out.CursorDown(1)
		if i == chosen {
			r.out.SetColor(r.selectedSuggestionTextColor, r.selectedSuggestionBGColor)
		} else {
			r.out.SetColor(r.suggestionTextColor, r.suggestionBGColor)
		}
		r.out.WriteStr(formatted[i])
		r.out.SetColor(White, DarkGray)
		r.out.Write([]byte(" "))
		r.out.CursorBackward(width + scrollBarWidth)
	}
	if d + width + scrollBarWidth > int(r.col) {
		r.out.CursorForward(d + width + scrollBarWidth - int(r.col))
	}

	r.out.CursorUp(l)
	r.out.SetColor(DefaultColor, DefaultColor)
	return
}

func (r *Render) Erase(buffer *Buffer) {
	r.out.CursorBackward(int(r.col))
	r.out.EraseDown()
	r.renderPrefix()
	r.out.Flush()
	return
}

func (r *Render) Render(buffer *Buffer, completions []string, chosen int) {
	line := buffer.Document().CurrentLine()
	r.out.SetColor(r.inputTextColor, r.inputBGColor)
	r.out.WriteStr(line)
	r.out.SetColor(DefaultColor, DefaultColor)
	r.out.CursorBackward(len(line) - buffer.CursorPosition)
	r.renderCompletion(buffer, completions, chosen)
	if chosen != -1 {
		c := completions[chosen]
		r.out.CursorBackward(len([]rune(buffer.Document().GetWordBeforeCursor())))
		r.out.SetColor(r.previewSuggestionTextColor, r.previewSuggestionBGColor)
		r.out.WriteStr(c)
		r.out.SetColor(DefaultColor, DefaultColor)
	}
	r.out.Flush()
}

func (r *Render) BreakLine(buffer *Buffer, result string) {
	r.out.SetColor(r.inputTextColor, r.inputBGColor)
	r.out.WriteStr(buffer.Document().Text + "\n")
	r.out.SetColor(r.outputTextColor, r.outputBGColor)
	r.out.WriteStr(result + "\n")
	r.out.SetColor(DefaultColor, DefaultColor)
	r.renderPrefix()
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
