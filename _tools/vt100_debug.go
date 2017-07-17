package main

import (
	"fmt"
	"syscall"

	"github.com/c-bata/go-prompt-toolkit"
	"github.com/pkg/term/termios"
)

const fd = 0

var orig syscall.Termios

func SetRawMode() {
	var new syscall.Termios
	if err := termios.Tcgetattr(uintptr(fd), &orig); err != nil {
		fmt.Println("Failed to get attribute")
		return
	}
	new = orig
	// "&^=" used like: https://play.golang.org/p/8eJw3JxS4O
	new.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.IEXTEN | syscall.ISIG
	new.Cc[syscall.VMIN] = 1
	new.Cc[syscall.VTIME] = 0
	termios.Tcsetattr(uintptr(fd), termios.TCSANOW, (*syscall.Termios)(&new))
}

func Restore() {
	termios.Tcsetattr(uintptr(fd), termios.TCSANOW, &orig)
}

func main() {
	SetRawMode()
	defer Restore()
	defer fmt.Println("exited!")

	bufCh := make(chan []byte, 128)
	go readBuffer(bufCh)
	fmt.Print("> ")
	parser := prompt.NewVT100StandardInputParser()

	for {
		b := <-bufCh
		if ac := parser.GetASCIICode(b); ac == nil {
			fmt.Printf("Key '%s' data:'%#v'\n", string(b), b)
		} else {
			if ac.Key == prompt.ControlC {
				fmt.Println("exit.")
				return
			}
			fmt.Printf("Key '%s' data:'%#v'\n", ac.Key, b)
		}
		fmt.Print("> ")
	}
}

func readBuffer(bufCh chan []byte) {
	buf := make([]byte, 1024)

	for {
		if n, err := syscall.Read(syscall.Stdin, buf); err == nil {
			bufCh <- buf[:n]
		}
	}
}
