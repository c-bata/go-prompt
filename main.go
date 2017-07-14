package main

import (
	"fmt"
	"syscall"
	"os"
	"os/signal"

	"github.com/c-bata/go-prompt-toolkit/prompt"
)

func scroll(out *prompt.VT100Writer, lines int) {
	for i := 0; i < lines; i++ {
		out.ScrollDown()
	}
	for i := 0; i < lines; i++ {
		out.ScrollUp()
	}
	return
}

func main() {
	in := prompt.NewVT100Parser()
	in.Setup()
	defer in.TearDown()
	defer fmt.Println("\nGoodbye!")
	out := prompt.NewVT100Writer()

	renderer := prompt.Render{
		Prefix: ">>> ",
		Out:    out,
	}

	out.SetTitle("はろー")
	scroll(out, 7)
	out.Flush()

	bufCh := make(chan []byte, 128)
	go readBuffer(bufCh)

	winSizeCh := make(chan *prompt.WinSize, 128)
	go updateWindowSize(in, winSizeCh)

	buffer := prompt.NewBuffer()

	for {
		b := <-bufCh
		ac := in.GetASCIICode(b)
		if ac == nil {
			buffer.InsertText(string(b), false, true)
			out.EraseDown()
			out.WriteRaw(b)
			after := buffer.Document().TextAfterCursor()
			out.WriteStr(after)
		} else if ac.Key == prompt.ControlC {
			out.EraseDown()
			out.ClearTitle()
			out.Flush()
			return
		} else {
			prompt.InputHandler(ac, buffer, out)
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

		completions := []string{"select", "insert", "update", "where"}
		renderer.Render(buffer, completions)
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

func updateWindowSize(in *prompt.VT100Parser, winSizeCh chan *prompt.WinSize) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(
		sigCh,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGWINCH,
	)

	for {
		s := <-sigCh
		switch s {
		// kill -SIGHUP XXXX
		case syscall.SIGHUP:

			// kill -SIGINT XXXX or Ctrl+c
		case syscall.SIGINT:

			// kill -SIGTERM XXXX
		case syscall.SIGTERM:

			// kill -SIGQUIT XXXX
		case syscall.SIGQUIT:

		case syscall.SIGWINCH:
			winSizeCh <- in.GetWinSize()
		default:
		}
	}
}
