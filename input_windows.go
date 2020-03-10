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
	row int
	col int
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
	if h == 0 {
		h = p.row
	}
	if w == 0 {
		w = p.col
	}
	return &WinSize{
		Row: uint16(h),
		Col: uint16(w),
	}
}

// SetWinSize sets default width and height of terminal when can not be optained automatically.
func (t *PosixParser) SetWinSize(ws *WinSize) {
	t.col = int(ws.Col)
	t.row = int(ws.Row)
}

// NewStandardInputParser returns ConsoleParser object to read from stdin.
func NewStandardInputParser() *WindowsParser {
	return &WindowsParser{}
}
