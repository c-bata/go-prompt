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

func (r *Render) UpdateWinSize(ws WinSize) {
	r.row = ws.Row
	r.col = ws.Col
	return
}

func (r *Render) RenderCompletion(words []string) {
	r.PrepareArea(4)
	r.out.SetColor("white", "teal")

	r.out.CursorDown(1)
	r.out.Write([]byte(" select "))
	r.out.SetColor("white", "darkGray")
	r.out.Write([]byte(" "))
	r.out.SetColor("white", "teal")
	r.out.CursorBackward(len("select") + 3)

	r.out.CursorDown(1)
	r.out.Write([]byte(" insert "))
	r.out.SetColor("white", "darkGray")
	r.out.Write([]byte(" "))
	r.out.SetColor("white", "teal")
	r.out.CursorBackward(len("insert") + 3)

	r.out.CursorDown(1)
	r.out.Write([]byte(" update "))
	r.out.SetColor("white", "darkGray")
	r.out.Write([]byte(" "))
	r.out.SetColor("white", "teal")
	r.out.CursorBackward(len("update") + 3)

	r.out.CursorDown(1)
	r.out.Write([]byte(" where  "))
	r.out.SetColor("white", "darkGray")
	r.out.Write([]byte(" "))
	r.out.SetColor("white", "teal")
	r.out.CursorBackward(len("where ") + 3)

	r.out.CursorUp(4)
	r.out.SetColor("default", "default")
	return
}

func (r *Render) Render(buffer *Buffer, completions []string) {
	r.out.Flush()
}
