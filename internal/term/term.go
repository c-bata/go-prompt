// +build !windows

package term

import (
	"sync"

	"github.com/pkg/term/termios"
	"golang.org/x/sys/unix"
)

var (
	saveTermios     *unix.Termios
	saveTermiosErr  error
	saveTermiosFD   int
	saveTermiosOnce sync.Once
)

func getOriginalTermios(fd int) (*unix.Termios, error) {
	saveTermiosOnce.Do(func() {
		saveTermiosFD = fd
		saveTermios, saveTermiosErr = termios.Tcgetattr(uintptr(fd))
	})
	return saveTermios, saveTermiosErr
}

// Restore terminal's mode.
func Restore() error {
	o, err := getOriginalTermios(saveTermiosFD)
	if err != nil {
		return err
	}
	return termios.Tcsetattr(uintptr(saveTermiosFD), termios.TCSANOW, o)
}
