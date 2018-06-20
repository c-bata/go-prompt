// +build !windows

package prompt

import (
	"syscall"
)

// PosixWriter is a ConsoleWriter implementation for POSIX environment.
// To control terminal emulator, this outputs VT100 escape sequences.
type PosixWriter struct {
	VT100Writer
	fd int
}

// Flush to flush buffer
func (w *PosixWriter) Flush() error {
	_, err := syscall.Write(w.fd, w.buffer)
	if err != nil {
		return err
	}
	w.buffer = []byte{}
	return nil
}

var _ ConsoleWriter = &PosixWriter{}

// NewStandardOutputWriter returns ConsoleWriter object to write to stdout.
// This generates VT100 escape sequences because almost terminal emulators
// in POSIX OS built on top of a VT100 specification.
func NewStandardOutputWriter() *PosixWriter {
	return &PosixWriter{
		fd: syscall.Stdout,
	}
}
