// +build windows

package prompt

import (
	"errors"
	"syscall"
	"unicode/utf8"
	"unsafe"

	tty "github.com/mattn/go-tty"
)

const maxReadBytes = 1024

var kernel32 = syscall.NewLazyDLL("kernel32.dll")

var procGetNumberOfConsoleInputEvents = kernel32.NewProc("GetNumberOfConsoleInputEvents")

// WindowsParser is a ConsoleParser implementation for Win32 console.
type WindowsParser struct {
	tty *tty.TTY
}

// Setup should be called before starting input
func (p *WindowsParser) Setup() error {
	t, err := tty.Open()
	if err != nil {
		return err
	}
	p.tty = t
	return nil
}

// TearDown should be called after stopping input
func (p *WindowsParser) TearDown() error {
	return p.tty.Close()
}

// Read returns byte array.
func (p *WindowsParser) Read() ([]byte, error) {
	var ev uint32
	// #nosec G103
	r0, _, err := procGetNumberOfConsoleInputEvents.Call(p.tty.Input().Fd(), uintptr(unsafe.Pointer(&ev)))
	if r0 == 0 {
		return nil, err
	}
	if ev == 0 {
		return nil, errors.New("EAGAIN")
	}

	r, err := p.tty.ReadRune()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, maxReadBytes)
	n := utf8.EncodeRune(buf[:], r)
	for p.tty.Buffered() && n < maxReadBytes {
		r, err := p.tty.ReadRune()
		if err != nil {
			break
		}
		n += utf8.EncodeRune(buf[n:], r)
	}
	return buf[:n], nil
}

// GetWinSize returns WinSize object to represent width and height of terminal.
func (p *WindowsParser) GetWinSize() *WinSize {
	w, h, err := p.tty.Size()
	if err != nil {
		panic(err)
	}
	return &WinSize{
		Row: uint16(h),
		Col: uint16(w),
	}
}

// NewStandardInputParser returns ConsoleParser object to read from stdin.
func NewStandardInputParser() *WindowsParser {
	return &WindowsParser{}
}
