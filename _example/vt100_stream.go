package main

import (
	"bytes"
	"fmt"
	"syscall"

	"github.com/c-bata/go-prompt-toolkit/prompt"
	"github.com/pkg/term/termios"
)

const fd = 0

var orig syscall.Termios

func SetRawMode() {
	var t syscall.Termios
	if err := termios.Tcgetattr(uintptr(fd), &orig); err != nil {
		fmt.Println("Failed to get attribute")
		return
	}
	t = orig
	// "&^=" used like: https://play.golang.org/p/8eJw3JxS4O
	t.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.IEXTEN | syscall.ISIG
	t.Cc[syscall.VMIN] = 1
	t.Cc[syscall.VTIME] = 0
	termios.Tcsetattr(uintptr(fd), termios.TCSANOW, (*syscall.Termios)(&t))
}

func Restore() {
	termios.Tcsetattr(uintptr(fd), termios.TCSANOW, &orig)
}

func GetASCIICode(b []byte) (ac *prompt.ASCIICode) {
	for _, k := range prompt.ASCII_SEQUENCES {
		if bytes.Compare(k.ASCIICode, b) == 0 {
			ac = k
			return
		}
	}
	return
}

func main() {
	SetRawMode()
	defer Restore()
	defer fmt.Println("exited!")

	bufCh := make(chan []byte, 128)
	go readBuffer(bufCh)
	fmt.Print("> ")

	for {
		b := <-bufCh
		if ac := GetASCIICode(b); ac == nil {
			fmt.Println(string(b))
		} else {
			if ac.Key == prompt.ControlC {
				fmt.Println("exit.")
				return
			}
			fmt.Println(ac.Key)
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
