package main

import "github.com/c-bata/go-prompt"

func main() {
	l := 20
	out := prompt.NewVT100StandardOutputWriter()
	out.EraseScreen()
	for i := 0; i < l; i++ {
		out.CursorGoTo(i, 0)
		out.WriteStr("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
	}

	out.CursorGoTo(5, 10)
	out.EraseLine()

	out.CursorGoTo(l, 0)
	out.Flush()
	return
}
