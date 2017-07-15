package prompt

type Render struct {
	Prefix         string
	Title          string
	out            *VT100Writer
	row            uint16
	col            uint16 // sigwinchで送られてくる列数を常に見ながら、prefixのlengthとbufferのcursor positionを比べて、completionの表示位置をずらす
	maxCompletions uint8
}

func (r *Render) Setup() {
	if r.Title != "" {
		r.out.SetTitle(r.Title)
	}
	r.out.WriteStr(r.Prefix)
	r.out.Flush()
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

func (r *Render) RenderCompletion(words []string, chosen int) {
	if len(words) == 0 {
		return
	}
	formatted, width := formatCompletions(words)
	l := len(formatted)
	r.prepareArea(l)

	r.out.SetColor("white", "teal")
	for i := 0; i < l; i++ {
		r.out.CursorDown(1)
		if i == chosen {
			r.out.SetColor("white", "turquoise")
		} else {
			r.out.SetColor("black", "cyan")
		}
		r.out.WriteStr(" " + formatted[i] + " ")
		r.out.SetColor("white", "darkGray")
		r.out.Write([]byte(" "))
		r.out.CursorBackward(width + 3)
	}

	r.out.CursorUp(l)
	r.out.SetColor("default", "default")
	return
}

func (r *Render) Erase(buffer *Buffer) {
	r.out.CursorBackward(len(r.Prefix))
	r.out.CursorBackward(buffer.CursorPosition + 100)
	r.out.EraseDown()
	r.out.WriteStr(r.Prefix)
	r.out.Flush()
	return
}

func (r *Render) Render(buffer *Buffer, completions []string, chosen int) {
	line := buffer.Document().CurrentLine()
	r.out.WriteStr(line)
	r.out.CursorBackward(len(line) - buffer.CursorPosition)
	r.RenderCompletion(completions, chosen)
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
	r.out.WriteStr(r.Prefix)
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
