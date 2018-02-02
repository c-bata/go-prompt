package main

import (
	"fmt"
	"os"
	"os/signal"
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
	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGWINCH,
	)
	ws := GetWinSize(syscall.Stdin)
	fmt.Printf("Row %d : Col %d\n", ws.Row, ws.Col)

	exitChan := make(chan int)
	go func() {
		for {
			s := <-signalChan
			switch s {
			// kill -SIGHUP XXXX
			case syscall.SIGHUP:
				exitChan <- 0

			// kill -SIGINT XXXX or Ctrl+c
			case syscall.SIGINT:
				exitChan <- 0

			// kill -SIGTERM XXXX
			case syscall.SIGTERM:
				exitChan <- 0

			// kill -SIGQUIT XXXX
			case syscall.SIGQUIT:
				exitChan <- 0

			case syscall.SIGWINCH:
				ws := GetWinSize(syscall.Stdin)
				fmt.Printf("Row %d : Col %d\n", ws.Row, ws.Col)

			default:
				exitChan <- 1
			}
		}
	}()

	code := <-exitChan
	os.Exit(code)
}
