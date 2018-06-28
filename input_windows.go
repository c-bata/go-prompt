// +build windows

package prompt

import (
	"context"
	"errors"
	"log"
	"syscall"
	"time"
	"unicode/utf8"
	"unsafe"

	"github.com/mattn/go-tty"
)

const maxReadBytes = 1024

var kernel32 = syscall.NewLazyDLL("kernel32.dll")

var procGetNumberOfConsoleInputEvents = kernel32.NewProc("GetNumberOfConsoleInputEvents")

// WindowsParser is a ConsoleParser implementation for Win32 console.
type WindowsParser struct {
	tty *tty.TTY
}

// SetUp should be called before starting input
func (p *WindowsParser) SetUp() error {
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

// NewStandardInputParser returns ConsoleParser object to read from stdin.
func NewStandardInputParser() ConsoleParser {
	return &WindowsParser{}
}

// Run start to worker goroutine.
func (ip *InputProcessor) Run(ctx context.Context) (err error) {
	log.Printf("[INFO] InputProcessor: Start running input processor")
	defer log.Print("[INFO] InputProcessor: Stop input processor")

	ip.in.SetUp()
	defer ip.in.TearDown()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			b, err := ip.in.Read()
			if err != nil {
				log.Printf("[ERROR] cannot read %s", err)
				return err
			}
			if !(len(b) == 1 && b[0] == 0) {
				ip.UserInput <- b
			}
			continue
		}
		time.Sleep(10 * time.Millisecond)
	}
}
