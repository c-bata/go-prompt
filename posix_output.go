// +build !windows

package prompt

import (
	"strconv"
	"syscall"
)

// PosixWriter is a ConsoleWriter implementation for POSIX environment.
// To control terminal emulator, this outputs VT100 escape sequences.
type PosixWriter struct {
	fd     int
	buffer []byte
}

// WriteRaw to write raw byte array
func (w *PosixWriter) WriteRaw(data []byte) {
	w.buffer = append(w.buffer, data...)
	// Flush because sometimes the render is broken when a large amount data in buffer.
	w.Flush()
	return
}

// Write to write byte array without control sequences
func (w *PosixWriter) Write(data []byte) {
	w.WriteRaw(byteFilter(data, writeFilter))
	return
}

// WriteRawStr to write raw string
func (w *PosixWriter) WriteRawStr(data string) {
	w.WriteRaw([]byte(data))
	return
}

// WriteStr to write string without control sequences
func (w *PosixWriter) WriteStr(data string) {
	w.Write([]byte(data))
	return
}

// Flush to flush buffer
func (w *PosixWriter) Flush() error {
	_, err := syscall.Write(w.fd, w.buffer)
	if err != nil {
		return err
	}
	w.buffer = []byte{}
	return nil
}

/* Erase */

// EraseScreen erases the screen with the background colour and moves the cursor to home.
func (w *PosixWriter) EraseScreen() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x32, 0x4a})
	return
}

// EraseUp erases the screen from the current line up to the top of the screen.
func (w *PosixWriter) EraseUp() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x31, 0x4a})
	return
}

// EraseDown erases the screen from the current line down to the bottom of the screen.
func (w *PosixWriter) EraseDown() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x4a})
	return
}

// EraseStartOfLine erases from the current cursor position to the start of the current line.
func (w *PosixWriter) EraseStartOfLine() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x31, 0x4b})
	return
}

// EraseEndOfLine erases from the current cursor position to the end of the current line.
func (w *PosixWriter) EraseEndOfLine() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x4b})
	return
}

// EraseLine erases the entire current line.
func (w *PosixWriter) EraseLine() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x32, 0x4b})
	return
}

/* Cursor */

// ShowCursor stops blinking cursor and show.
func (w *PosixWriter) ShowCursor() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x3f, 0x31, 0x32, 0x6c, 0x1b, 0x5b, 0x3f, 0x32, 0x35, 0x68})
}

// HideCursor hides cursor.
func (w *PosixWriter) HideCursor() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x3f, 0x32, 0x35, 0x6c})
	return
}

// CursorGoTo sets the cursor position where subsequent text will begin.
func (w *PosixWriter) CursorGoTo(row, col int) {
	if row == 0 && col == 0 {
		// If no row/column parameters are provided (ie. <ESC>[H), the cursor will move to the home position.
		w.WriteRaw([]byte{0x1b, 0x5b, 0x3b, 0x48})
		return
	}
	r := strconv.Itoa(row)
	c := strconv.Itoa(col)
	w.WriteRaw([]byte{0x1b, 0x5b})
	w.WriteRaw([]byte(r))
	w.WriteRaw([]byte{0x3b})
	w.WriteRaw([]byte(c))
	w.WriteRaw([]byte{0x48})
	return
}

// CursorUp moves the cursor up by 'n' rows; the default count is 1.
func (w *PosixWriter) CursorUp(n int) {
	if n == 0 {
		return
	} else if n < 0 {
		w.CursorDown(-n)
		return
	}
	s := strconv.Itoa(n)
	w.WriteRaw([]byte{0x1b, 0x5b})
	w.WriteRaw([]byte(s))
	w.WriteRaw([]byte{0x41})
	return
}

// CursorDown moves the cursor down by 'n' rows; the default count is 1.
func (w *PosixWriter) CursorDown(n int) {
	if n == 0 {
		return
	} else if n < 0 {
		w.CursorUp(-n)
		return
	}
	s := strconv.Itoa(n)
	w.WriteRaw([]byte{0x1b, 0x5b})
	w.WriteRaw([]byte(s))
	w.WriteRaw([]byte{0x42})
	return
}

// CursorForward moves the cursor forward by 'n' columns; the default count is 1.
func (w *PosixWriter) CursorForward(n int) {
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

// CursorBackward moves the cursor backward by 'n' columns; the default count is 1.
func (w *PosixWriter) CursorBackward(n int) {
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

// AskForCPR asks for a cursor position report (CPR).
func (w *PosixWriter) AskForCPR() {
	// CPR: Cursor Position Request.
	w.WriteRaw([]byte{0x1b, 0x5b, 0x36, 0x6e})
	w.Flush()
	return
}

// SaveCursor saves current cursor position.
func (w *PosixWriter) SaveCursor() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x73})
	return
}

// UnSaveCursor restores cursor position after a Save Cursor.
func (w *PosixWriter) UnSaveCursor() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x75})
	return
}

/* Scrolling */

// ScrollDown scrolls display down one line.
func (w *PosixWriter) ScrollDown() {
	w.WriteRaw([]byte{0x1b, 0x44})
	return
}

// ScrollUp scroll display up one line.
func (w *PosixWriter) ScrollUp() {
	w.WriteRaw([]byte{0x1b, 0x4d})
	return
}

/* Title */

// SetTitle sets a title of terminal window.
func (w *PosixWriter) SetTitle(title string) {
	w.WriteRaw([]byte{0x1b, 0x5d, 0x32, 0x3b})
	w.WriteRaw(byteFilter([]byte(title), setTextFilter))
	w.WriteRaw([]byte{0x07})
	return
}

// ClearTitle clears a title of terminal window.
func (w *PosixWriter) ClearTitle() {
	w.WriteRaw([]byte{0x1b, 0x5d, 0x32, 0x3b, 0x07})
	return
}

/* Font */

// SetColor sets text and background colors. and specify whether text is bold.
func (w *PosixWriter) SetColor(fg, bg Color, bold bool) {
	f, ok := foregroundANSIColors[fg]
	if !ok {
		f = foregroundANSIColors[DefaultColor]
	}
	b, ok := backgroundANSIColors[bg]
	if !ok {
		b = backgroundANSIColors[DefaultColor]
	}
	syscall.Write(syscall.Stdout, []byte{0x1b, 0x5b, 0x33, 0x39, 0x3b, 0x34, 0x39, 0x6d})
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

var _ ConsoleWriter = &PosixWriter{}

// NewStandardOutputWriter returns ConsoleWriter object to write to stdout.
func NewStandardOutputWriter() *PosixWriter {
	return &PosixWriter{
		fd: syscall.Stdout,
	}
}
