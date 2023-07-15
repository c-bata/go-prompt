// +build !windows

package term

import (
	"sync"

	"github.com/pkg/term/termios"
	"golang.org/x/sys/unix"
)

var (
	saveTermios     *unix.Termios
	saveTermiosFD   int
	saveTermiosOnce sync.Once
)

func getOriginalTermios(fd int) (unix.Termios, error) {
	var err error
	saveTermiosOnce.Do(func() {
		saveTermiosFD = fd
		saveTermios, err = termios.Tcgetattr(uintptr(fd))
	})
	return *saveTermios, err
}

// Restore terminal's mode.
func Restore() error {
	o, err := getOriginalTermios(saveTermiosFD)
	if err != nil {
		return err
	}
	return termios.Tcsetattr(uintptr(saveTermiosFD), termios.TCSANOW, &o)
}
