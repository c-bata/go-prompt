// +build !windows

package prompt

import (
	"syscall"

	"github.com/c-bata/go-prompt/internal/term"
	"golang.org/x/sys/unix"
)

const maxReadBytes = 1024

// PosixParser is a ConsoleParser implementation for POSIX environment.
type PosixParser struct {
	fd          int
	row         uint16
	col         uint16
	origTermios syscall.Termios
}

// Setup should be called before starting input
func (t *PosixParser) Setup() error {
	// Set NonBlocking mode because if syscall.Read block this goroutine, it cannot receive data from stopCh.
	if err := syscall.SetNonblock(t.fd, true); err != nil {
		return err
	}
	if err := term.SetRaw(t.fd); err != nil {
		return err
	}
	return nil
}

// TearDown should be called after stopping input
func (t *PosixParser) TearDown() error {
	if err := syscall.SetNonblock(t.fd, false); err != nil {
		return err
	}
	if err := term.Restore(); err != nil {
		return err
	}
	return nil
}

// Read returns byte array.
func (t *PosixParser) Read() ([]byte, error) {
	buf := make([]byte, maxReadBytes)
	n, err := syscall.Read(t.fd, buf)
	if err != nil {
		return []byte{}, err
	}
	return buf[:n], nil
}

// GetWinSize returns WinSize object to represent width and height of terminal.
func (t *PosixParser) GetWinSize() *WinSize {
	ws, err := unix.IoctlGetWinsize(t.fd, unix.TIOCGWINSZ)
	if err != nil {
		panic(err)
	}
	if ws.Row == 0 {
		ws.Row = t.row
	}
	if ws.Col == 0 {
		ws.Col = t.col
	}
	return &WinSize{
		Row: ws.Row,
		Col: ws.Col,
	}
}

// SetWinSize sets default width and height of terminal when can not be optained automatically.
func (t *PosixParser) SetWinSize(ws *WinSize) {
	t.col = ws.Col
	t.row = ws.Row
}

var _ ConsoleParser = &PosixParser{}

// NewStandardInputParser returns ConsoleParser object to read from stdin.
func NewStandardInputParser() *PosixParser {
	in, err := syscall.Open("/dev/tty", syscall.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}

	return &PosixParser{
		fd: in,
	}
}
