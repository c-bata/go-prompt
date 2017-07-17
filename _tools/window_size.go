package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"unsafe"
)

// winsize is winsize struct got from the ioctl(2) system call.
type winsize struct {
	Row uint16
	Col uint16
	X   uint16 // pixel value
	Y   uint16 // pixel value
}

// GetWinSize returns winsize struct which is the response of ioctl(2).
func GetWinSize(fd int) (row, col uint16) {
	ws := &winsize{}
	retCode, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(fd),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		panic(errno)
	}
	return ws.Row, ws.Col
}

func main() {
	signal_chan := make(chan os.Signal, 1)
	signal.Notify(
		signal_chan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGWINCH,
	)

	exit_chan := make(chan int)
	go func() {
		for {
			s := <-signal_chan
			switch s {
			// kill -SIGHUP XXXX
			case syscall.SIGHUP:
				exit_chan <- 0

			// kill -SIGINT XXXX or Ctrl+c
			case syscall.SIGINT:
				exit_chan <- 0

			// kill -SIGTERM XXXX
			case syscall.SIGTERM:
				exit_chan <- 0

			// kill -SIGQUIT XXXX
			case syscall.SIGQUIT:
				exit_chan <- 0

			case syscall.SIGWINCH:
				r, c := GetWinSize(syscall.Stdin)
				fmt.Printf("Row %d : Col %d\n", r, c)

			default:
				exit_chan <- 1
			}
		}
	}()

	code := <-exit_chan
	os.Exit(code)
}
