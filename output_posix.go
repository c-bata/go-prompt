// +build !windows

package prompt

import (
	"syscall"
)

type stdoutWriter struct{}

func (w *stdoutWriter) Write(b []byte) (n int, err error) {
	return syscall.Write(syscall.Stdout, b)
}

// NewStandardOutputWriter returns ConsoleWriter object to write to stdout.
func NewStandardOutputWriter() *PosixWriter {
	return &PosixWriter{
		w: &stdoutWriter{},
	}
}
