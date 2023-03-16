package prompt

import (
	"io"

	"github.com/c-bata/go-prompt/internal/debug"
)

type ioWriter struct {
	VT100Writer
	w io.Writer
}

// Flush to flush buffer.
func (w *ioWriter) Flush() error {
	//_log.Infow("before flush", "message", string(w.buffer))
	_, err := w.w.Write(w.buffer)
	if err != nil {
		debug.Logf("unable to flush, error=%v\n", err)
		return err
	}
	w.buffer = []byte{}
	return nil
}

var _ ConsoleWriter = (*ioWriter)(nil)

// NewIOWriter returns new console writer which writes to io.Writer
func NewIOWriter(w io.Writer) ConsoleWriter {
	return &ioWriter{w: w}
}
