package prompt

type Render struct {
	Prefix         string
	Title          string
	out            ConsoleWriter
	row            uint16
	col            uint16
	maxCompletions uint16
}

func (r *Render) Setup() {
	if r.Title != "" {
		r.out.SetTitle(r.Title)
	}
	r.renderPrefix()
	r.out.Flush()
}

func (r *Render) renderPrefix() {
	r.out.SetColor("green", "default")
	r.out.WriteStr(r.Prefix)
	r.out.SetColor("default", "default")
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

	formatted, width := formatCompletions(words)
	l := len(formatted)
	r.prepareArea(l)

	d := (len(r.Prefix) + len(buf.Document().TextBeforeCursor())) % int(r.col)
	if d + width + 3 > int(r.col) {
		r.out.CursorBackward(d + width + 3 - int(r.col))
	}

	r.out.SetColor("white", "teal")
	for i := 0; i < l; i++ {
		r.out.CursorDown(1)
		if i == chosen {
			r.out.SetColor("black", "turquoise")
		} else {
			r.out.SetColor("white", "cyan")
		}
		r.out.WriteStr(" " + formatted[i] + " ")
		r.out.SetColor("white", "darkGray")
		r.out.Write([]byte(" "))
		r.out.CursorBackward(width + 3)
	}
	if d + width + 3 > int(r.col) {
		r.out.CursorForward(d + width + 3 - int(r.col))
	}

	r.out.CursorUp(l)
	r.out.SetColor("default", "default")
	return
}

func (r *Render) Erase(buffer *Buffer) {
	r.out.CursorBackward(len(r.Prefix))
	r.out.CursorBackward(buffer.CursorPosition + 100)
	r.out.EraseDown()
	r.renderPrefix()
	r.out.Flush()
	return
}

func (r *Render) Render(buffer *Buffer, completions []string, chosen int) {
	line := buffer.Document().CurrentLine()
	r.out.WriteStr(line)
	r.out.CursorBackward(len(line) - buffer.CursorPosition)
	r.renderCompletion(buffer, completions, chosen)
	if chosen != -1 {
		c := completions[chosen]
		r.out.CursorBackward(len([]rune(buffer.Document().GetWordBeforeCursor())))
		r.out.WriteStr(c)
	}
	r.out.Flush()
}

func (r *Render) BreakLine(buffer *Buffer, result string) {
	r.out.WriteStr(buffer.Document().Text)
	r.out.WriteStr("\n")
	r.out.WriteStr(result)
	r.out.WriteStr("\n")
	r.renderPrefix()
}

func formatCompletions(words []string) (new []string, width int) {
	num := len(words)
	new = make([]string, num)
	width = 0

	for i := 0; i < num; i++ {
		if width < len([]rune(words[i])) {
			width = len([]rune(words[i]))
		}
	}

	for i := 0; i < num; i++ {
		spaces := width - len([]rune(words[i]))
		new[i] = words[i]
		for j := 0; j < spaces; j++ {
			new[i] += " "
		}
	}
	return
}
