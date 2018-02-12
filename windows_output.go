// +build windows

package prompt

import (
	"io"
	"strconv"

	"github.com/mattn/go-colorable"
)

type WindowsWriter struct {
	out    io.Writer
	buffer []byte
}

func (w *WindowsWriter) WriteRaw(data []byte) {
	w.buffer = append(w.buffer, data...)
	// Flush because sometimes the render is broken when a large amount data in buffer.
	w.Flush()
	return
}

func (w *WindowsWriter) Write(data []byte) {
	w.WriteRaw(byteFilter(data, writeFilter))
	return
}

func (w *WindowsWriter) WriteRawStr(data string) {
	w.WriteRaw([]byte(data))
	return
}

func (w *WindowsWriter) WriteStr(data string) {
	w.Write([]byte(data))
	return
}

func (w *WindowsWriter) Flush() error {
	_, err := w.out.Write(w.buffer)
	if err != nil {
		return err
	}
	w.buffer = []byte{}
	return nil
}

/* Erase */

func (w *WindowsWriter) EraseScreen() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x32, 0x4a})
	return
}

func (w *WindowsWriter) EraseUp() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x31, 0x4a})
	return
}

func (w *WindowsWriter) EraseDown() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x4a})
	return
}

func (w *WindowsWriter) EraseStartOfLine() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x31, 0x4b})
	return
}

func (w *WindowsWriter) EraseEndOfLine() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x4b})
	return
}

func (w *WindowsWriter) EraseLine() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x32, 0x4b})
	return
}

/* Cursor */

func (w *WindowsWriter) ShowCursor() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x3f, 0x31, 0x32, 0x6c, 0x1b, 0x5b, 0x3f, 0x32, 0x35, 0x68})
}

func (w *WindowsWriter) HideCursor() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x3f, 0x32, 0x35, 0x6c})
	return
}

func (w *WindowsWriter) CursorGoTo(row, col int) {
	r := strconv.Itoa(row)
	c := strconv.Itoa(col)
	w.WriteRaw([]byte{0x1b, 0x5b})
	w.WriteRaw([]byte(r))
	w.WriteRaw([]byte{0x3b})
	w.WriteRaw([]byte(c))
	w.WriteRaw([]byte{0x48})
	return
}

func (w *WindowsWriter) CursorUp(n int) {
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

func (w *WindowsWriter) CursorDown(n int) {
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

func (w *WindowsWriter) CursorForward(n int) {
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

func (w *WindowsWriter) CursorBackward(n int) {
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

func (w *WindowsWriter) AskForCPR() {
	// CPR: Cursor Position Request.
	w.WriteRaw([]byte{0x1b, 0x5b, 0x36, 0x6e})
	w.Flush()
	return
}

func (w *WindowsWriter) SaveCursor() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x73})
	return
}

func (w *WindowsWriter) UnSaveCursor() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x75})
	return
}

/* Scrolling */

func (w *WindowsWriter) ScrollDown() {
	w.WriteRaw([]byte{0x1b, 0x44})
	return
}

func (w *WindowsWriter) ScrollUp() {
	w.WriteRaw([]byte{0x1b, 0x4d})
	return
}

/* Title */

func (w *WindowsWriter) SetTitle(title string) {
	w.WriteRaw([]byte{0x1b, 0x5d, 0x32, 0x3b})
	w.WriteRaw(byteFilter([]byte(title), setTextFilter))
	w.WriteRaw([]byte{0x07})
	return
}

func (w *WindowsWriter) ClearTitle() {
	w.WriteRaw([]byte{0x1b, 0x5d, 0x32, 0x3b, 0x07})
	return
}

/* Font */

func (w *WindowsWriter) SetColor(fg, bg Color, bold bool) {
	f, ok := foregroundANSIColors[fg]
	if !ok {
		f, _ = foregroundANSIColors[DefaultColor]
	}
	b, ok := backgroundANSIColors[bg]
	if !ok {
		b, _ = backgroundANSIColors[DefaultColor]
	}
	w.out.Write([]byte{0x1b, 0x5b, 0x33, 0x39, 0x3b, 0x34, 0x39, 0x6d})
	w.WriteRaw([]byte{0x1b, 0x5b})
	if !bold {
		w.WriteRaw([]byte{0x30, 0x3b})
	}
	w.WriteRaw(f)
	w.WriteRaw([]byte{0x3b})
	w.WriteRaw(b)
	if bold {
		w.WriteRaw([]byte{0x3b, 0x31})
	}
	w.WriteRaw([]byte{0x6d})
	return
}

var foregroundANSIColors = map[Color][]byte{
	DefaultColor: {0x33, 0x39}, // 39

	// Low intensity.
	Black:     {0x33, 0x30}, // 30
	DarkRed:   {0x33, 0x31}, // 31
	DarkGreen: {0x33, 0x32}, // 32
	Brown:     {0x33, 0x33}, // 33
	DarkBlue:  {0x33, 0x34}, // 34
	Purple:    {0x33, 0x35}, // 35
	Cyan:      {0x33, 0x36}, //36
	LightGray: {0x33, 0x37}, //37

	// High intensity.
	DarkGray:  {0x39, 0x30}, // 90
	Red:       {0x39, 0x31}, // 91
	Green:     {0x39, 0x32}, // 92
	Yellow:    {0x39, 0x33}, // 93
	Blue:      {0x39, 0x34}, // 94
	Fuchsia:   {0x39, 0x35}, // 95
	Turquoise: {0x39, 0x36}, // 96
	White:     {0x39, 0x37}, // 97
}

var backgroundANSIColors = map[Color][]byte{
	DefaultColor: {0x34, 0x39}, // 49

	// Low intensity.
	Black:     {0x34, 0x30}, // 40
	DarkRed:   {0x34, 0x31}, // 41
	DarkGreen: {0x34, 0x32}, // 42
	Brown:     {0x34, 0x33}, // 43
	DarkBlue:  {0x34, 0x34}, // 44
	Purple:    {0x34, 0x35}, // 45
	Cyan:      {0x34, 0x36}, // 46
	LightGray: {0x34, 0x37}, // 47

	// High intensity
	DarkGray:  {0x31, 0x30, 0x30}, // 100
	Red:       {0x31, 0x30, 0x31}, // 101
	Green:     {0x31, 0x30, 0x32}, // 102
	Yellow:    {0x31, 0x30, 0x33}, // 103
	Blue:      {0x31, 0x30, 0x34}, // 104
	Fuchsia:   {0x31, 0x30, 0x35}, // 105
	Turquoise: {0x31, 0x30, 0x36}, // 106
	White:     {0x31, 0x30, 0x37}, // 107
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

var _ ConsoleWriter = &WindowsWriter{}

func NewStandardOutputWriter() *WindowsWriter {
	return &WindowsWriter{
		out: colorable.NewColorableStdout(),
	}
}
