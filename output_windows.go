// +build windows

package prompt

import (
	"context"
	"io"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-tty"
)

// WindowsWriter is a ConsoleWriter implementation for Win32 console.
// Output is converted from VT100 escape sequences by mattn/go-colorable.
type WindowsWriter struct {
	VT100Writer
	tty    *tty.TTY
	writer io.Writer
}

// Flush to flush buffer
func (w *WindowsWriter) Flush() error {
	_, err := w.writer.Write(w.buffer)
	if err != nil {
		return err
	}
	w.buffer = []byte{}
	return nil
}

// GetWinSize returns WinSize object to represent width and height of terminal.
func (w *WindowsWriter) GetWinSize() WinSize {
	col, row, err := w.tty.Size()
	if err != nil {
		panic(err)
	}
	return WinSize{
		Row: uint16(row),
		Col: uint16(col),
	}
}

// SIGWINCH returns WinSize channel which send WinSize when window size is changed.
func (w *WindowsWriter) SIGWINCH(ctx context.Context) chan WinSize {
	sigwinch := make(chan WinSize, 1)
	go func() {
		select {
		case <-ctx.Done():
			return
		case ws := <-w.tty.SIGWINCH():
			sigwinch <- WinSize{Row: uint16(ws.H), Col: uint16(ws.W)}
		}
	}()
	return sigwinch
}

var _ ConsoleWriter = &WindowsWriter{}

// NewStandardOutputWriter returns ConsoleWriter object to write to stdout.
// This generates win32 control sequences.
func NewStandardOutputWriter() (ConsoleWriter, error) {
	t, err := tty.Open()
	if err != nil {
		return nil, err
	}
	w := &WindowsWriter{
		tty:    t,
		writer: colorable.NewColorable(t.Output()),
	}
	return w, nil
}
