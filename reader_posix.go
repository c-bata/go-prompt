//go:build !windows
// +build !windows

package prompt

import (
	"os"
	"syscall"

	"github.com/elk-language/go-prompt/internal/term"
	"golang.org/x/sys/unix"
)

// PosixReader is a Reader implementation for the POSIX environment.
type PosixReader struct {
	fd int
}

// Open should be called before starting input
func (t *PosixReader) Open() error {
	in, err := syscall.Open("/dev/tty", syscall.O_RDONLY, 0)
	if os.IsNotExist(err) {
		in = syscall.Stdin
	} else if err != nil {
		panic(err)
	}
	t.fd = in
	// Set NonBlocking mode because if syscall.Read block this goroutine, it cannot receive data from stopCh.
	if err := syscall.SetNonblock(t.fd, true); err != nil {
		return err
	}
	if err := term.SetRaw(t.fd); err != nil {
		return err
	}
	return nil
}

// Close should be called after stopping input
func (t *PosixReader) Close() error {
	if err := syscall.Close(t.fd); err != nil {
		return err
	}
	if err := term.Restore(); err != nil {
		return err
	}
	return nil
}

// Read returns byte array.
func (t *PosixReader) Read(buff []byte) (int, error) {
	return syscall.Read(t.fd, buff)
}

// GetWinSize returns WinSize object to represent width and height of terminal.
func (t *PosixReader) GetWinSize() *WinSize {
	ws, err := unix.IoctlGetWinsize(t.fd, unix.TIOCGWINSZ)
	if err != nil {
		// If this errors, we simply return the default window size as
		// it's our best guess.
		return &WinSize{
			Row: 25,
			Col: 80,
		}
	}
	return &WinSize{
		Row: ws.Row,
		Col: ws.Col,
	}
}

var _ Reader = &PosixReader{}

// NewStdinReader returns Reader object to read from stdin.
func NewStdinReader() *PosixReader {
	return &PosixReader{}
}
