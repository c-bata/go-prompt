package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/c-bata/go-prompt"
	"github.com/pkg/term/termios"
)

const fd = 0

var (
	orig syscall.Termios
)

func init() {
	if err := termios.Tcgetattr(uintptr(fd), &orig); err != nil {
		fmt.Println("Failed to get attribute")
		return
	}
}

func setRawMode() {
	n := orig
	n.Iflag &^= syscall.IGNBRK | syscall.BRKINT | syscall.PARMRK |
		syscall.ISTRIP | syscall.INLCR | syscall.IGNCR |
		syscall.ICRNL | syscall.IXON
	n.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.IEXTEN | syscall.ISIG | syscall.ECHONL
	n.Cflag &^= syscall.CSIZE | syscall.PARENB
	n.Cc[syscall.VMIN] = 1
	n.Cc[syscall.VTIME] = 0
	termios.Tcsetattr(uintptr(fd), termios.TCSANOW, (*syscall.Termios)(&n))
}

func restoreTermios() {
	termios.Tcsetattr(uintptr(fd), termios.TCSANOW, &orig)
}

func main() {
	_, _, e := syscall.Syscall(syscall.SYS_FCNTL, uintptr(fd), uintptr(syscall.F_SETFL),
		uintptr(syscall.O_ASYNC|syscall.O_NONBLOCK))
	if e != 0 {
		fmt.Printf("[ERROR] Cannot set non-blocking mode: %d\n", e)
		return
	}
	_, _, e = syscall.Syscall(syscall.SYS_FCNTL, uintptr(fd), uintptr(syscall.F_SETOWN),
		uintptr(syscall.Getpid()))
	if runtime.GOOS != "darwin" && e != 0 {
		fmt.Printf("[ERROR] Cannot set F_SETOWN: %d\n", e)
		return
	}
	defer func() {
		err := syscall.SetNonblock(fd, false)
		if err != nil {
			fmt.Printf("[ERROR] Cannot unset non-blocking mode: %s\n", err)
		}
	}()

	// Set raw mode
	setRawMode()
	defer restoreTermios()

	// Set signal handler
	sigio := make(chan os.Signal, 1)
	signal.Notify(sigio, syscall.SIGIO)
	sigquit := make(chan os.Signal, 1)
	signal.Notify(sigquit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	buf := make([]byte, 128)
	for {
		fmt.Print("> ")
		select {
		case <-sigio:
			n, err := syscall.Read(syscall.Stdin, buf)
			if err == syscall.EAGAIN || err == syscall.EWOULDBLOCK {
				continue
			} else if err != nil {
				fmt.Printf("[ERROR] cannot read %s", err)
				break
			}
			b := buf[:n]
			if key := prompt.GetKey(b); key == prompt.NotDefined {
				fmt.Printf("Key '%s' data:'%#v'\n", string(b), b)
			} else {
				fmt.Printf("Key '%s' data:'%#v'\n", key, b)
				if key == prompt.ControlC {
					return
				}
			}
		case <-sigquit:
			return
		}
	}
}
