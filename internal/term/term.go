// +build !windows

package term

import (
	"sync"
	"syscall"

	"github.com/pkg/term/termios"
)

var (
	saveTermios     syscall.Termios
	saveTermiosFD   int
	saveTermiosOnce sync.Once
)

func getOriginalTermios(fd int) (syscall.Termios, error) {
	var err error
	saveTermiosOnce.Do(func() {
		saveTermiosFD = fd
		err = termios.Tcgetattr(uintptr(fd), &saveTermios)
	})
	return saveTermios, err
}

// Restore terminal's mode.
func Restore() error {
	o, err := getOriginalTermios(saveTermiosFD)
	if err != nil {
		return err
	}
	return termios.Tcsetattr(uintptr(saveTermiosFD), termios.TCSANOW, &o)
}
