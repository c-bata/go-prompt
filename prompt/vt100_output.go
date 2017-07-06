package prompt

import (
	"strconv"
	"syscall"
)

type VT100Writer struct {
	fd     int
	buffer []byte
}

func (w *VT100Writer) Write(data []byte) {
	w.WriteRaw(byteFilter(data, writeFilter))
	return
}

func (w *VT100Writer) WriteStr(data string) {
	w.Write([]byte(data))
	return
}

func (w *VT100Writer) WriteRaw(data []byte) {
	w.buffer = append(w.buffer, data...)
	return
}

func (w *VT100Writer) Flush() error {
	_, err := syscall.Write(w.fd, w.buffer)
	if err != nil {
		return err
	}
	w.buffer = []byte{}
	return nil
}

func (w *VT100Writer) Clear() {
	syscall.Write(w.fd, []byte{0x1b, 0x5b, 0x02, 0x6a, 0x1b, 0x63})
	return
}

/* Cursor */

func (w *VT100Writer) CursorUp(n int) {
	if n < 0 {
		w.CursorDown(n)
		return
	}
	s := strconv.Itoa(n)
	w.WriteRaw([]byte{0x1b, 0x5b})
	w.WriteRaw([]byte(s))
	w.WriteRaw([]byte{0x41})
	return
}

func (w *VT100Writer) CursorDown(n int) {
	if n < 0 {
		w.CursorUp(n)
		return
	}
	s := strconv.Itoa(n)
	w.WriteRaw([]byte{0x1b, 0x5b})
	w.WriteRaw([]byte(s))
	w.WriteRaw([]byte{0x42})
	return
}

func (w *VT100Writer) CursorForward(n int) {
	if n < 0 {
		w.CursorBackward(n)
		return
	}
	s := strconv.Itoa(n)
	w.WriteRaw([]byte{0x1b, 0x5b})
	w.WriteRaw([]byte(s))
	w.WriteRaw([]byte{0x43})
	return
}

func (w *VT100Writer) CursorBackward(n int) {
	if n < 0 {
		w.CursorForward(n)
		return
	}
	s := strconv.Itoa(n)
	w.WriteRaw([]byte{0x1b, 0x5b})
	w.WriteRaw([]byte(s))
	w.WriteRaw([]byte{0x44})
	return
}

/* Title */

func (w *VT100Writer) SetTitle(title string) {
	w.WriteRaw([]byte{0x1b, 0x5d, 0x32, 0x3b})
	w.WriteRaw(byteFilter([]byte(title), setTextFilter))
	w.WriteRaw([]byte{0x07})
	return
}

func (w *VT100Writer) ClearTitle() {
	w.WriteRaw([]byte{0x1b, 0x5d, 0x32, 0x3b, 0x07})
	return
}

func NewVT100Writer() *VT100Writer {
	return &VT100Writer{
		fd: syscall.Stdout,
	}
}

func writeFilter(buf byte) bool {
	return buf != 0x1b && buf != 0x3f
}

func setTextFilter(buf byte) bool {
	return buf != 0x1b && buf != 0x07
}

func byteFilter(buf []byte, fn ...func(b byte) bool) []byte {
	if len(fn) == 0 {
		return buf
	}
	ret := make([]byte, 0, len(buf))
	f := fn[0]
	for i, n := range buf {
		if f(n) {
			ret = append(ret, buf[i])
		}
	}
	return byteFilter(ret, fn[1:]...)
}
