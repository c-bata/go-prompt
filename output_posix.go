// +build !windows

package prompt

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"unsafe"
)

const flushMaxRetryCount = 3

// PosixWriter is a ConsoleWriter implementation for POSIX environment.
// To control terminal emulator, this outputs VT100 escape sequences.
type PosixWriter struct {
	VT100Writer
	fd int
}

// Flush to flush buffer
func (w *PosixWriter) Flush() error {
	l := len(w.buffer)
	offset := 0
	retry := 0
	for {
		n, err := syscall.Write(w.fd, w.buffer[offset:])
		if err != nil {
			log.Printf("[DEBUG] flush error: %s", err)
			if retry < flushMaxRetryCount {
				retry++
				continue
			}
			return err
		}
		offset += n
		if offset == l {
			break
		}
	}
	w.buffer = []byte{}
	return nil
}

// winsize is winsize struct got from the ioctl(2) system call.
type ioctlWinsize struct {
	Row uint16
	Col uint16
	X   uint16 // pixel value
	Y   uint16 // pixel value
}

// GetWinSize returns WinSize object to represent width and height of terminal.
func (w *PosixWriter) GetWinSize() WinSize {
	ws := &ioctlWinsize{}
	retCode, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(w.fd),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		panic(errno)
	}
	return WinSize{
		Row: ws.Row,
		Col: ws.Col,
	}
}

// SIGWINCH returns WinSize channel which send WinSize when window size is changed.
func (w *PosixWriter) SIGWINCH(ctx context.Context) chan WinSize {
	sigwinsz := make(chan WinSize, 1)
	sigwinch := make(chan os.Signal, 1)
	signal.Notify(sigwinch, syscall.SIGWINCH)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-sigwinch:
				sigwinsz <- w.GetWinSize()
			}
		}
	}()
	return sigwinsz
}

var _ ConsoleWriter = &PosixWriter{}

// NewStandardOutputWriter returns ConsoleWriter object to write to stdout.
// This generates VT100 escape sequences because almost terminal emulators
// in POSIX OS built on top of a VT100 specification.
func NewStandardOutputWriter() (ConsoleWriter, error) {
	return &PosixWriter{
		fd: syscall.Stdout,
	}, nil
}
