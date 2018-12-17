// +build !windows

package main

import (
	"fmt"
	"syscall"

	prompt "github.com/c-bata/go-prompt"
	"github.com/c-bata/go-prompt/internal/term"
)

func main() {
	if err := term.SetRaw(syscall.Stdin); err != nil {
		fmt.Println(err)
		return
	}
	defer term.Restore()

	bufCh := make(chan []byte, 128)
	go readBuffer(bufCh)
	fmt.Print("> ")

	for {
		b := <-bufCh
		if key := prompt.GetKey(b); key == prompt.NotDefined {
			fmt.Printf("Key '%s' data:'%#v'\n", string(b), b)
		} else {
			if key == prompt.ControlC {
				fmt.Println("exit.")
				return
			}
			fmt.Printf("Key '%s' data:'%#v'\n", key, b)
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
