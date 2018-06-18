// +build windows

package prompt

import (
	"bytes"
	"io"
	"strconv"

	"github.com/mattn/go-colorable"
)

// WindowsWriter is a ConsoleWriter implementation for Win32 console.
// Output is converted from VT100 escape sequences by mattn/go-colorable.
type WindowsWriter struct {
	out    io.Writer
	buffer []byte
}

// WriteRaw to write raw byte array
func (w *WindowsWriter) WriteRaw(data []byte) {
	w.buffer = append(w.buffer, data...)
	return
}

// Write to write safety byte array by removing control sequences.
func (w *WindowsWriter) Write(data []byte) {
	w.WriteRaw(bytes.Replace(data, []byte{0x1b}, []byte{'?'}, -1))
	return
}

// WriteRawStr to write raw string
func (w *WindowsWriter) WriteRawStr(data string) {
	w.WriteRaw([]byte(data))
	return
}

// WriteStr to write safety string by removing control sequences.
func (w *WindowsWriter) WriteStr(data string) {
	w.Write([]byte(data))
	return
}

// Flush to flush buffer
func (w *WindowsWriter) Flush() error {
	_, err := w.out.Write(w.buffer)
	if err != nil {
		return err
	}
	w.buffer = []byte{}
	return nil
}

/* Erase */

// EraseScreen erases the screen with the background colour and moves the cursor to home.
func (w *WindowsWriter) EraseScreen() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x32, 0x4a})
	return
}

// EraseUp erases the screen from the current line up to the top of the screen.
func (w *WindowsWriter) EraseUp() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x31, 0x4a})
	return
}

// EraseDown erases the screen from the current line down to the bottom of the screen.
func (w *WindowsWriter) EraseDown() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x4a})
	return
}

// EraseStartOfLine erases from the current cursor position to the start of the current line.
func (w *WindowsWriter) EraseStartOfLine() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x31, 0x4b})
	return
}

// EraseEndOfLine erases from the current cursor position to the end of the current line.
func (w *WindowsWriter) EraseEndOfLine() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x4b})
	return
}

// EraseLine erases the entire current line.
func (w *WindowsWriter) EraseLine() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x32, 0x4b})
	return
}

/* Cursor */

// ShowCursor stops blinking cursor and show.
func (w *WindowsWriter) ShowCursor() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x3f, 0x31, 0x32, 0x6c, 0x1b, 0x5b, 0x3f, 0x32, 0x35, 0x68})
}

// HideCursor hides cursor.
func (w *WindowsWriter) HideCursor() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x3f, 0x32, 0x35, 0x6c})
	return
}

// CursorGoTo sets the cursor position where subsequent text will begin.
func (w *WindowsWriter) CursorGoTo(row, col int) {
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

// CursorDown moves the cursor down by 'n' rows; the default count is 1.
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

// CursorForward moves the cursor forward by 'n' columns; the default count is 1.
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

// CursorBackward moves the cursor backward by 'n' columns; the default count is 1.
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

// AskForCPR asks for a cursor position report (CPR).
func (w *WindowsWriter) AskForCPR() {
	// CPR: Cursor Position Request.
	w.WriteRaw([]byte{0x1b, 0x5b, 0x36, 0x6e})
	w.Flush()
	return
}

// SaveCursor saves current cursor position.
func (w *WindowsWriter) SaveCursor() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x73})
	return
}

// UnSaveCursor restores cursor position after a Save Cursor.
func (w *WindowsWriter) UnSaveCursor() {
	w.WriteRaw([]byte{0x1b, 0x5b, 0x75})
	return
}

/* Scrolling */

// ScrollDown scrolls display down one line.
func (w *WindowsWriter) ScrollDown() {
	w.WriteRaw([]byte{0x1b, 0x44})
	return
}

// ScrollUp scroll display up one line.
func (w *WindowsWriter) ScrollUp() {
	w.WriteRaw([]byte{0x1b, 0x4d})
	return
}

/* Title */

// SetTitle sets a title of terminal window.
func (w *WindowsWriter) SetTitle(title string) {
	titleBytes := []byte(title)
	patterns := []struct {
		from []byte
		to   []byte
	}{
		{
			from: []byte{0x13},
			to:   []byte{},
		},
		{
			from: []byte{0x07},
			to:   []byte{},
		},
	}
	for i := range patterns {
		titleBytes = bytes.Replace(titleBytes, patterns[i].from, patterns[i].to, -1)
	}

	w.WriteRaw([]byte{0x1b, 0x5d, 0x32, 0x3b})
	w.WriteRaw(titleBytes)
	w.WriteRaw([]byte{0x07})
	return
}

// ClearTitle clears a title of terminal window.
func (w *WindowsWriter) ClearTitle() {
	w.WriteRaw([]byte{0x1b, 0x5d, 0x32, 0x3b, 0x07})
	return
}

/* Font */

// SetColor sets text and background colors. and specify whether text is bold.
func (w *WindowsWriter) SetColor(fg, bg Color, bold bool) {
	f, ok := foregroundANSIColors[fg]
	if !ok {
		f, _ = foregroundANSIColors[DefaultColor]
	}
	b, ok := backgroundANSIColors[bg]
	if !ok {
		b, _ = backgroundANSIColors[DefaultColor]
	}
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

var _ ConsoleWriter = &WindowsWriter{}

// NewStandardOutputWriter returns ConsoleWriter object to write to stdout.
func NewStandardOutputWriter() *WindowsWriter {
	return &WindowsWriter{
		out: colorable.NewColorableStdout(),
	}
}
