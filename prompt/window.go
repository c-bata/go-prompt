package prompt

import (
	"syscall"
	"unsafe"
)

// GetWinSize returns winsize struct which is the response of ioctl(2).
var GetWinSize = getWinSize

// Winsize is winsize struct got from the ioctl(2) system call.
type Winsize struct {
	Row uint16
	Col uint16
	X   uint16 // pixel value
	Y   uint16 // pixel value
}

func getWinSize(fileno int) (ws *Winsize) {
	ws = &Winsize{}
	retCode, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(fileno),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		panic(errno)
	}
	return ws
}
