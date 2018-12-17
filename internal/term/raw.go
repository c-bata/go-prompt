// +build !windows

package term

import (
	"syscall"

	"github.com/pkg/term/termios"
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
	n.Cc[syscall.VMIN] = 1
	n.Cc[syscall.VTIME] = 0
	return termios.Tcsetattr(uintptr(fd), termios.TCSANOW, (*syscall.Termios)(&n))
}
