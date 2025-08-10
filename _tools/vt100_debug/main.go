//go:build !windows
// +build !windows

package main

import (
	"fmt"
	"syscall"

	prompt "github.com/c-bata/go-prompt"
	"github.com/c-bata/go-prompt/internal/term"
	"golang.org/x/sys/unix"
)

func main() {
	if err := term.SetRaw(syscall.Stdin); err != nil {
		fmt.Println(err)
		return
	}
	defer term.Restore()

	matcher := prompt.NewSequenceMatcher()
	bufCh := make(chan []byte, 128)
	go readBuffer(bufCh)
	fmt.Print("> ")

	inputBuffer := make([]byte, 0, 32)
	for {
		b := <-bufCh
		inputBuffer = append(inputBuffer, b...)
		for len(inputBuffer) > 0 {
			result, key := matcher.MatchSequence(inputBuffer)
			switch result {
			case prompt.Exact:
				if key != nil && *key == prompt.ControlC {
					fmt.Println("exit.")
					return
				}
				fmt.Printf("Key '%s' data:'%#v'\n", key, inputBuffer[:])
				inputBuffer = inputBuffer[:0]
			case prompt.Prefix:
				// Wait for more bytes
				// Wait for more bytes, exit inner loop
				return
			case prompt.NoMatch:
				fmt.Printf("Key '%s' data:'%#v'\n", string(inputBuffer[0]), inputBuffer[:1])
				inputBuffer = inputBuffer[1:]
			}
			fmt.Print("> ")
			// If Prefix, break to wait for more input
			if result == prompt.Prefix {
				break
			}
		}
	}
}

func readBuffer(bufCh chan []byte) {
	buf := make([]byte, 1024)
	fd := int(syscall.Stdin)
	for {
		var readfds unix.FdSet
		readfds.Set(fd)
		n, err := unix.Select(fd+1, &readfds, nil, nil, nil)
		if err != nil {
			continue
		}
		if n > 0 && readfds.IsSet(fd) {
			if rn, err := syscall.Read(syscall.Stdin, buf); err == nil && rn > 0 {
				bufCh <- buf[:rn]
			}
		}
	}
}
