// +build !windows

package term

import (
	"syscall"

	"github.com/pkg/term/termios"
	"golang.org/x/sys/unix"
)

// SetRaw put terminal into a raw mode
func SetRaw(fd int) error {
	n, err := getOriginalTermios(fd)
	if err != nil {
		return err
	}

	n.Iflag &^= syscall.IGNBRK | syscall.BRKINT | syscall.PARMRK |
		syscall.ISTRIP | syscall.INLCR | syscall.IGNCR |
		syscall.ICRNL | syscall.IXON
	n.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.IEXTEN | syscall.ISIG | syscall.ECHONL
	n.Cflag &^= syscall.CSIZE | syscall.PARENB
	n.Cflag |= syscall.CS8 // Set to 8-bit wide.  Typical value for displaying characters.
	n.Cc[syscall.VMIN] = 1
	n.Cc[syscall.VTIME] = 0

	return termios.Tcsetattr(uintptr(fd), termios.TCSANOW, (*unix.Termios)(&n))
}
