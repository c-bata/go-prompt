package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"unsafe"
)

// Winsize is winsize struct got from the ioctl(2) system call.
type Winsize struct {
	Row uint16
	Col uint16
	X   uint16 // pixel value
	Y   uint16 // pixel value
}

// GetWinSize returns winsize struct which is the response of ioctl(2).
func GetWinSize(fd int) *Winsize {
	ws := &Winsize{}
	retCode, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(fd),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		panic(errno)
	}
	return ws
}

func main() {
	var fd = syscall.Stdout
	sigwinch := make(chan os.Signal, 1)
	signal.Notify(sigwinch, syscall.SIGWINCH)

	var wg sync.WaitGroup
	defer func() {
		wg.Wait()
	}()
	wg.Add(1)
	ws := GetWinSize(fd)
	fmt.Printf("Row %d : Col %d\n", ws.Row, ws.Col)

	go func() {
		for {
			select {
			case <-sigwinch:
				ws := GetWinSize(fd)
				fmt.Printf("Row %d : Col %d\n", ws.Row, ws.Col)
			}
		}
		wg.Done()
	}()
}
