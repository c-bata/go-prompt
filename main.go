package main

import (
	"fmt"
	"syscall"

	"github.com/c-bata/go-prompt-toolkit/prompt"
)

func enterAlternateScreen(fd int) {
	syscall.Write(fd, []byte{0x1b, 0x5b, 0x3f, 0x01, 0x00, 0x04, 0x09, 0x68, 0x1b, 0x5b, 0x48})
}

func scroll(out *prompt.VT100Writer, lines int) {
	for i := 0; i < lines; i++ {
		out.ScrollDown()
		defer out.ScrollUp()
	}
	return
}

func main() {
	in := prompt.NewVT100Parser()
	in.Setup()
	defer in.TearDown()
	defer fmt.Println("\nexited!")
	out := prompt.NewVT100Writer()
	out.SetTitle("はろー")
	scroll(out, 7)
	out.Flush()

	bufCh := make(chan []byte, 128)
	go readBuffer(bufCh)

	buffer := prompt.NewBuffer()

	for {
		b := <-bufCh
		if ac := in.GetASCIICode(b); ac == nil {
			out.EraseDown()
			out.WriteRaw(b)
			buffer.InsertText(string(b), false, true)
		} else if ac.Key == prompt.Enter || ac.Key == prompt.ControlJ {
			out.EraseDown()
			out.WriteStr(buffer.Document().TextAfterCursor())

			out.WriteStr("\n>>> Your input: '")
			out.WriteStr(buffer.Text())
			out.WriteStr("' <<<\n")
			buffer = prompt.NewBuffer()
		} else if ac.Key == prompt.Left {
			l := buffer.CursorLeft(1)
			if l == 0 {
				continue
			}
			out.EraseLine()
			out.EraseDown()
			after := buffer.Document().CurrentLine()
			out.WriteStr(after)
			out.CursorBackward(len(after) - buffer.CursorPosition)
		} else if ac.Key == prompt.Right {
			l := buffer.CursorRight(1)
			if l == 0 {
				continue
			}

			out.CursorForward(l)
			out.WriteRaw(b)
			out.EraseDown()
			after := buffer.Document().TextAfterCursor()
			out.WriteStr(after)
		} else if ac.Key == prompt.Backspace {
			deleted := buffer.DeleteBeforeCursor(1)
			if deleted == "" {
				continue
			}
			out.CursorBackward(1)
			out.EraseDown()

			after := buffer.Document().TextAfterCursor()
			out.WriteStr(after)
		} else if ac.Key == prompt.Tab || ac.Key == prompt.ControlI {
		} else if ac.Key == prompt.BackTab {
		} else if ac.Key == prompt.Right {
			buffer.CursorRight(1)
		} else if ac.Key == prompt.ControlT {
			enterAlternateScreen(syscall.Stdout)
		} else if ac.Key == prompt.ControlC {
			out.EraseDown()
			out.ClearTitle()
			out.Flush()
			return
		} else if ac.Key == prompt.Up || ac.Key == prompt.Down {
		} else {
			out.WriteRaw(b)
			//buffer.InsertText(ac.Key.String(), false, true)
		}

		// Display completions
		if w := buffer.Document().GetWordBeforeCursor(); w != "" {
			out.SetColor("white", "teal")

			out.CursorDown(1)
			out.Write([]byte(" select "))
			out.SetColor("white", "darkGray")
			out.Write([]byte(" "))
			out.SetColor("white", "teal")
			out.CursorBackward(len("select") + 3)

			out.CursorDown(1)
			out.Write([]byte(" insert "))
			out.SetColor("white", "darkGray")
			out.Write([]byte(" "))
			out.SetColor("white", "teal")
			out.CursorBackward(len("insert") + 3)

			out.CursorDown(1)
			out.Write([]byte(" update "))
			out.SetColor("white", "darkGray")
			out.Write([]byte(" "))
			out.SetColor("white", "teal")
			out.CursorBackward(len("update") + 3)

			out.CursorDown(1)
			out.Write([]byte(" where  "))
			out.SetColor("white", "darkGray")
			out.Write([]byte(" "))
			out.SetColor("white", "teal")
			out.CursorBackward(len("where ") + 3)

			out.CursorUp(4)
			out.SetColor("default", "default")
		}

		scroll(out, 4)
		out.Flush()
	}
}

func readBuffer(bufCh chan []byte) {
	buf := make([]byte, 1024)

	for {
		if n, err := syscall.Read(syscall.Stdin, buf); err == nil {
			bufCh <- buf[:n]
		}
	}
}
