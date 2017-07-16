package main

import "github.com/c-bata/go-prompt-toolkit/prompt"

func main() {
	l := 20
	out := prompt.NewVT100Writer()
	for i := 0; i < l; i++ {
		out.CursorGoTo(i, 0)
		out.WriteStr("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
	}


	out.CursorGoTo(5, 10)
	out.EraseLine()

	out.CursorGoTo(l, 0)
	out.Flush()
	return
}
