package prompt

type Render struct {
	Prefix string
	out    *VT100Writer
	row    uint16
	col    uint16 // sigwinchで送られてくる列数を常に見ながら、prefixのlengthとbufferのcursor positionを比べて、completionの表示位置をずらす
}

func (r *Render) PrepareArea(lines int) {
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

func (r *Render) RenderCompletion(words []string) {
	formatted, width := formatCompletions(words)
	l := len(formatted)
	r.PrepareArea(l)

	r.out.SetColor("white", "teal")
	for i := 0; i < l; i++ {
		r.out.CursorDown(1)
		r.out.WriteStr(" " + formatted[i] + " ")
		r.out.SetColor("white", "darkGray")
		r.out.Write([]byte(" "))
		r.out.SetColor("white", "teal")
		r.out.CursorBackward(width + 3)
	}

	r.out.CursorUp(l)
	r.out.SetColor("default", "default")
	return
}

func formatCompletions(words []string) ([]string, int) {
	num := len(words)
	width := 0

	for i := 0; i < num; i++ {
		if width < len([]rune(words[i])) {
			width = len([]rune(words[i]))
		}
	}

	for i := 0; i < num; i++ {
		spaces := width - len([]rune(words[i]))
		for j := 0; j < spaces; j++ {
			words[i] += " "
		}
	}
	return words, width
}

func (r *Render) Render(buffer *Buffer, completions []string) {
	r.out.Flush()
}
