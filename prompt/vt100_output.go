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

/* Erase */

func (w *VT100Writer) EraseScreen() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x32, 0x4a})
	return
}

func (w *VT100Writer) EraseUp() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x31, 0x4a})
	return
}

func (w *VT100Writer) EraseDown() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x4a})
	return
}

func (w *VT100Writer) EraseStartOfLine() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x31, 0x4b})
	return
}

func (w *VT100Writer) EraseEndOfLine() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x4b})
	return
}

func (w *VT100Writer) EraseLine() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x32, 0x4b})
	return
}

/* Cursor */

func (w *VT100Writer) ShowCursor() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x3f, 0x31, 0x32, 0x6c, 0x1b, 0x5b, 0x3f, 0x32, 0x35, 0x68})
}

func (w *VT100Writer) HideCursor() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x3f, 0x32, 0x35, 0x6c})
	return
}

func (w *VT100Writer) CursorGoTo(row, col int) {
	r := strconv.Itoa(row)
	c := strconv.Itoa(col)
	w.WriteRaw([]byte{0x1b, 0x5b})
	w.WriteRaw([]byte(r))
	w.WriteRaw([]byte{0x3b})
	w.WriteRaw([]byte(c))
	w.WriteRaw([]byte{0x48})
	return
}

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
	if n == 0 {
		return
	} else if n < 0 {
		w.CursorBackward(-n)
		return
	}
	s := strconv.Itoa(n)
	w.WriteRaw([]byte{0x1b, 0x5b})
	w.WriteRaw([]byte(s))
	w.WriteRaw([]byte{0x43})
	return
}

func (w *VT100Writer) CursorBackward(n int) {
	if n == 0 {
		return
	} else if n < 0 {
		w.CursorForward(-n)
		return
	}
	s := strconv.Itoa(n)
	w.WriteRaw([]byte{0x1b, 0x5b})
	w.WriteRaw([]byte(s))
	w.WriteRaw([]byte{0x44})
	return
}

func (w *VT100Writer) AskForCPR() {
	// CPR: Cursor Position Request.
	w.WriteRaw([]byte{0x1b, 0x5b, 0x36, 0x6e})
	w.Flush()
	return
}

func (w *VT100Writer) SaveCursor() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x73})
	return
}

func (w *VT100Writer) UnSaveCursor() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x75})
	return
}

/* Scrolling */

func (w *VT100Writer) ScrollDown() {
	w.WriteRaw([]byte{0x1b, 0x44})
	return
}

func (w *VT100Writer) ScrollUp() {
	w.WriteRaw([]byte{0x1b, 0x4d})
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

/* utility */

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

/* colors */

func (w *VT100Writer) SetColor(fg, bg string) (ok bool) {
	f, ok := ForegroundANSIColors[fg]
	if !ok {
		return
	}
	b, ok := BackgroundANSIColors[bg]
	if !ok {
		return
	}
	syscall.Write(syscall.Stdout, []byte{0x1b, 0x5b, 0x33, 0x39, 0x3b, 0x34, 0x39, 0x6d})
	w.WriteRaw([]byte{0x1b, 0x5b})
	w.WriteRaw(f)
	w.WriteRaw([]byte{0x3b})
	w.WriteRaw(b)
	w.WriteRaw([]byte{0x6d})
	return
}

var ForegroundANSIColors = map[string][]byte{
	"default": []byte{0x33, 0x39}, // 39

	// Low intensity.
	"black":     []byte{0x33, 0x30}, // 30
	"darkRed":   []byte{0x33, 0x31}, // 31
	"darkGreen": []byte{0x33, 0x32}, // 32
	"brown":     []byte{0x33, 0x33}, // 33
	"darkBlue":  []byte{0x33, 0x34}, // 34
	"purple":    []byte{0x33, 0x35}, // 35
	"teal":      []byte{0x33, 0x36}, //36
	"lightGray": []byte{0x33, 0x37}, //37

	// High intensity.
	"darkGray":  []byte{0x39, 0x30}, // 90
	"red":       []byte{0x39, 0x31}, // 91
	"green":     []byte{0x39, 0x32}, // 92
	"yellow":    []byte{0x39, 0x33}, // 93
	"blue":      []byte{0x39, 0x34}, // 94
	"fuchsia":   []byte{0x39, 0x35}, // 95
	"turquoise": []byte{0x39, 0x36}, // 96
	"white":     []byte{0x39, 0x37}, // 97
}

var BackgroundANSIColors = map[string][]byte{
	"default": []byte{0x34, 0x39}, // 49

	// Low intensity.
	"black":     []byte{0x34, 0x30}, // 40
	"darkRed":   []byte{0x34, 0x31}, // 41
	"darkGreen": []byte{0x34, 0x32}, // 42
	"brown":     []byte{0x34, 0x33}, // 43
	"darkBlue":  []byte{0x34, 0x34}, // 44
	"purple":    []byte{0x34, 0x35}, // 45
	"teal":      []byte{0x34, 0x36}, // 46
	"lightGray": []byte{0x34, 0x37}, // 47

	// High intensity
	"darkGray":  []byte{0x31, 0x30, 0x30}, // 100
	"red":       []byte{0x31, 0x30, 0x31}, // 101
	"green":     []byte{0x31, 0x30, 0x32}, // 102
	"yellow":    []byte{0x31, 0x30, 0x33}, // 103
	"blue":      []byte{0x31, 0x30, 0x34}, // 104
	"fuchsia":   []byte{0x31, 0x30, 0x35}, // 105
	"turquoise": []byte{0x31, 0x30, 0x36}, // 106
	"white":     []byte{0x31, 0x30, 0x37}, // 107
}
