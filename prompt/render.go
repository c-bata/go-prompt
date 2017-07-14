package prompt

type Render struct {
	Prefix string
	Out *VT100Writer
	row uint16
	col uint16 // sigwinchで送られてくる列数を常に見ながら、prefixのlengthとbufferのcursor positionを比べて、completionの表示位置をずらす
}

func (r *Render) UpdateWinSize(ws WinSize) {
	r.row = ws.Row
	r.col = ws.Col
	return
}

func (r *Render) renderCompletion(words []string) {
	return
}

func (r *Render) Render(buffer *Buffer, completions []string) {
}
