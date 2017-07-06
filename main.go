package main

import (
	"fmt"
	"syscall"

	"github.com/c-bata/go-prompt-toolkit/prompt"
)

func HideCursor(fd int) {
	syscall.Write(fd, []byte{0x1b, 0x5b, 0x3f, 0x25, 0x6c})
}

func CursorPosition(fd, row, col int) {
	//syscall.Write(fd, []byte{0x1b, 0x5b, byte(row), 0x3b, byte(col), 0x0f})
	syscall.Write(fd, []byte{0x1b, 0x5b, 0x03, 0x3b, 0x04, 0x48})
}

func enterAlternateScreen(fd int) {
	syscall.Write(fd, []byte{0x1b, 0x5b, 0x3f, 0x01, 0x00, 0x04, 0x09, 0x68, 0x1b, 0x5b, 0x48})
}

func main() {
	in := prompt.NewVT100Parser()
	out := prompt.NewVT100Writer()
	in.Setup()
	defer in.TearDown()
	defer fmt.Println("exited!")
	out.SetTitle("はろー")
	defer out.ClearTitle()

	bufCh := make(chan []byte, 128)
	go readBuffer(bufCh)

	buffer := prompt.NewBuffer()

	for {
		b := <-bufCh
		if ac := in.GetASCIICode(b); ac == nil {
			buffer.InsertText(string(b), false, true)
		} else if ac.Key == prompt.Enter || ac.Key == prompt.ControlJ {
			buffer.InsertText("\n", false, true)
		} else if ac.Key == prompt.Left {
			buffer.CursorLeft(1)
			out.CursorBackward(1)
		} else if ac.Key == prompt.Right {
			buffer.CursorRight(1)
		} else if ac.Key == prompt.Up {
			buffer.CursorUp(1)
		} else if ac.Key == prompt.Down {
			buffer.CursorDown(1)
		} else if ac.Key == prompt.ControlT {
			enterAlternateScreen(syscall.Stdout)
		} else {
			if ac.Key == prompt.ControlC {
				return
			}
			buffer.InsertText(ac.Key.String(), false, true)
		}
		out.Clear()
		out.WriteStr(buffer.Text())
		out.CursorBackward(1)
		out.CursorBackward(1)
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
