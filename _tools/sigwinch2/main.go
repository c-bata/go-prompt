// +build windows

package main

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/mattn/go-tty"
)

const maxReadByteLen = 1024

var kernel32 = syscall.NewLazyDLL("kernel32.dll")

var procGetNumberOfConsoleInputEvents = kernel32.NewProc("GetNumberOfConsoleInputEvents")

func main() {
	t, err := tty.Open()
	if err != nil {
		return
	}
	sigwinch := t.SIGWINCH()

	go func() {
		for {
			select {
			default:
				var ev uint32
				r0, _, err := procGetNumberOfConsoleInputEvents.Call(t.Input().Fd(), uintptr(unsafe.Pointer(&ev)))
				if r0 == 0 {
					fmt.Println(err)
					return
				}
				if ev == 0 {
					fmt.Println("EAGAIN")
				}
				t.ReadRune()
			}
			t.ReadRune()
		}
	}()

	for {
		ws := <-sigwinch
		fmt.Printf("Row %d : Col %d\n", ws.H, ws.W)
	}
}
