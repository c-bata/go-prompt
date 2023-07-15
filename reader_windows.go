//go:build windows
// +build windows

package prompt

import (
	"errors"
	"syscall"
	"unicode/utf8"
	"unsafe"

	tty "github.com/mattn/go-tty"
)

var kernel32 = syscall.NewLazyDLL("kernel32.dll")

var procGetNumberOfConsoleInputEvents = kernel32.NewProc("GetNumberOfConsoleInputEvents")

// WindowsReader is a Reader implementation for Win32 console.
type WindowsReader struct {
	tty *tty.TTY
}

// Open should be called before starting input
func (p *WindowsReader) Open() error {
	t, err := tty.Open()
	if err != nil {
		return err
	}
	p.tty = t
	return nil
}

// Close should be called after stopping input
func (p *WindowsReader) Close() error {
	return p.tty.Close()
}

// Read returns byte array.
func (p *WindowsReader) Read(buff []byte) (int, error) {
	var ev uint32
	r0, _, err := procGetNumberOfConsoleInputEvents.Call(p.tty.Input().Fd(), uintptr(unsafe.Pointer(&ev)))
	if r0 == 0 {
		return 0, err
	}
	if ev == 0 {
		return 0, errors.New("EAGAIN")
	}

	r, err := p.tty.ReadRune()
	if err != nil {
		return 0, err
	}

	n := utf8.EncodeRune(buff[:], r)
	for p.tty.Buffered() && n < len(buff) {
		r, err := p.tty.ReadRune()
		if err != nil {
			break
		}
		n += utf8.EncodeRune(buff[n:], r)
	}
	return n, nil
}

// GetWinSize returns WinSize object to represent width and height of terminal.
func (p *WindowsReader) GetWinSize() *WinSize {
	w, h, err := p.tty.Size()
	if err != nil {
		panic(err)
	}
	return &WinSize{
		Row: uint16(h),
		Col: uint16(w),
	}
}

var _ Reader = &WindowsReader{}

// NewStdinReader returns Reader object to read from stdin.
func NewStdinReader() *WindowsReader {
	return &WindowsReader{}
}
