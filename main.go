package main

import (
	"fmt"
	"syscall"

	"github.com/c-bata/go-prompt-toolkit/prompt"
)

func ClearTitle(fd int) {
	seq := []byte{0x1b, 0x5d, 0x02, 0x3b, 0x07}
	syscall.Write(fd, seq)
	return
}

func SetTitle(fd int, title string) {
	var t []byte
	for i := range title {
		if title[i] != 0x0b && title[i] != 0x07 {
			t = append(t, title[i])
		}
	}
	seqStart := []byte{0x1b, 0x5d, 0x02, 0x3b}
	seqEnd := []byte{0x07}
	syscall.Write(fd, seqStart)
	syscall.Write(fd, t)
	syscall.Write(fd, seqEnd)
}

func HideCursor(fd int) {
	syscall.Write(fd, []byte{0x1b, 0x5b, 0x3f, 0x25, 0x6c})
}

func Clear(fd int) {
	syscall.Write(fd, []byte{0x1b, 0x5b, 0x02, 0x6a, 0x1b, 0x63})
}

func CursorBackward(fd int) {
	syscall.Write(fd, []byte{0x1b, 0x5b, 0x01, 0x44})
}

func CursorPosition(fd, row, col int) {
	//syscall.Write(fd, []byte{0x1b, 0x5b, byte(row), 0x3b, byte(col), 0x0f})
	syscall.Write(fd, []byte{0x1b, 0x5b, 0x03, 0x3b, 0x04, 0x48})
}

func enterAlternateScreen(fd int) {
	syscall.Write(fd, []byte{0x1b, 0x5b, 0x3f, 0x01, 0x00, 0x04, 0x09, 0x68, 0x1b, 0x5b, 0x48})
}

func main() {
	t := prompt.NewVT100Parser()
	t.Setup()
	defer t.TearDown()
	defer fmt.Println("exited!")

	fd := syscall.Stdout
	ClearTitle(fd)

	bufCh := make(chan []byte, 128)
	go readBuffer(bufCh)

	buffer := prompt.NewBuffer()

	for {
		b := <-bufCh
		if ac := t.GetASCIICode(b); ac == nil {
			buffer.InsertText(string(b), false, true)
		} else if ac.Key == prompt.Enter || ac.Key == prompt.ControlJ {
			buffer.InsertText("\n", false, true)
		} else if ac.Key == prompt.Left {
			buffer.CursorLeft(1)
		} else if ac.Key == prompt.Right {
			buffer.CursorRight(1)
		} else if ac.Key == prompt.Up {
			buffer.CursorUp(1)
		} else if ac.Key == prompt.Down {
			buffer.CursorDown(1)
		} else if ac.Key == prompt.ControlT {
			enterAlternateScreen(fd)
		} else {
			if ac.Key == prompt.ControlC {
				return
			}
			buffer.InsertText(ac.Key.String(), false, true)
		}
		Clear(fd)
		fmt.Print(buffer.Text())
		CursorBackward(fd)
		CursorBackward(fd)
		CursorPosition(fd, 3, 3)
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
