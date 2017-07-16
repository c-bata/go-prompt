package prompt

type Render struct {
	out            ConsoleWriter
	prefix         string
	title          string
	row            uint16
	col            uint16
	maxCompletions uint16
	// colors
	prefixColor                 Color
	textColor                   Color
	outputTextColor             Color
	completedTextColor          Color
	completionTextColor         Color
	completionBGColor           Color
	selectedCompletionTextColor Color
	selectedCompletionBGColor   Color
}

func (r *Render) Setup() {
	if r.title != "" {
		r.out.SetTitle(r.title)
	}
	r.renderPrefix()
	r.out.Flush()
}

func (r *Render) renderPrefix() {
	r.out.SetColor(r.prefixColor, DefaultColor)
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
	if l := len(words); l == 0 {
		return
	} else if l > int(r.maxCompletions) - 2 || l >= int(r.row) - 2 {
		if r.maxCompletions > r.row {
			words = words[:int(r.row) - 2]
		} else {
			words = words[:int(r.maxCompletions) - 2]
		}
	}

	formatted, width := formatCompletions(words, int(r.col) - len(r.prefix) - 3)
	l := len(formatted)
	r.prepareArea(l)

	d := (len(r.prefix) + len(buf.Document().TextBeforeCursor())) % int(r.col)
	if d + width + 3 > int(r.col) {
		r.out.CursorBackward(d + width + 3 - int(r.col))
	}

	r.out.SetColor(White, Cyan)
	for i := 0; i < l; i++ {
		r.out.CursorDown(1)
		if i == chosen {
			r.out.SetColor(r.selectedCompletionTextColor, r.selectedCompletionBGColor)
		} else {
			r.out.SetColor(r.completionTextColor, r.completionBGColor)
		}
		r.out.WriteStr(" " + formatted[i] + " ")
		r.out.SetColor(White, DarkGray)
		r.out.Write([]byte(" "))
		r.out.CursorBackward(width + 3)
	}
	if d + width + 3 > int(r.col) {
		r.out.CursorForward(d + width + 3 - int(r.col))
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
	r.out.SetColor(r.textColor, DefaultColor)
	r.out.WriteStr(line)
	r.out.SetColor(DefaultColor, DefaultColor)
	r.out.CursorBackward(len(line) - buffer.CursorPosition)
	r.renderCompletion(buffer, completions, chosen)
	if chosen != -1 {
		c := completions[chosen]
		r.out.CursorBackward(len([]rune(buffer.Document().GetWordBeforeCursor())))
		r.out.SetColor(r.completedTextColor, DefaultColor)
		r.out.WriteStr(c)
		r.out.SetColor(DefaultColor, DefaultColor)
	}
	r.out.Flush()
}

func (r *Render) BreakLine(buffer *Buffer, result string) {
	r.out.SetColor(r.textColor, DefaultColor)
	r.out.WriteStr(buffer.Document().Text + "\n")
	r.out.SetColor(r.outputTextColor, DefaultColor)
	r.out.WriteStr(result + "\n")
	r.out.SetColor(DefaultColor, DefaultColor)
	r.renderPrefix()
}

func formatCompletions(words []string, max int) (new []string, width int) {
	num := len(words)
	new = make([]string, num)
	width = 0

	for i := 0; i < num; i++ {
		if width < len([]rune(words[i])) {
			width = len([]rune(words[i]))
		}
	}

	if width > max {
		width = max
	}

	for i := 0; i < num; i++ {
		if l := len(words[i]); l > width {
			new[i] = words[i][:width - 3] + "..."
		} else if l < width  {
			spaces := width - len([]rune(words[i]))
			new[i] = words[i]
			for j := 0; j < spaces; j++ {
				new[i] += " "
			}
		} else {
			new[i] = words[i]
		}
	}
	return
}
