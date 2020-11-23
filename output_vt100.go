package prompt

import (
	"bytes"
	"strconv"

	fcolor "github.com/fatih/color"
)

// VT100Writer generates VT100 escape sequences. Thread unsafe
type VT100Writer struct {
	buffer []byte
}

// writeRaw to write raw byte array
func (w *VT100Writer) writeRaw(data []byte) {
	w.buffer = append(w.buffer, data...)
}

// Write to write safety byte array by removing control sequences.
func (w *VT100Writer) write(data []byte, color *fcolor.Color) {
	safeData := bytes.Replace(data, []byte{0x1b}, []byte{'?'}, -1)
	if color != nil {
		out := color.Sprint(string(safeData))
		safeData = []byte(out)
	}
	w.buffer = append(w.buffer, safeData...)
}

// writeRawStr to write raw string
func (w *VT100Writer) writeRawStr(data string) {
	w.writeRaw([]byte(data))
}

// WriteStr to write safety string by removing control sequences.
func (w *VT100Writer) WriteStr(data string, color *fcolor.Color) {
	w.write([]byte(data), color)
}

/* Erase */

// EraseScreen erases the screen with the background colour and moves the cursor to home.
func (w *VT100Writer) EraseScreen() {
	w.writeRaw([]byte{0x1b, '[', '2', 'J'})
}

// EraseUp erases the screen from the current line up to the top of the screen.
func (w *VT100Writer) EraseUp() {
	w.writeRaw([]byte{0x1b, '[', '1', 'J'})
}

// EraseDown erases the screen from the current line down to the bottom of the screen.
func (w *VT100Writer) EraseDown() {
	w.writeRaw([]byte{0x1b, '[', 'J'})
}

// EraseStartOfLine erases from the current cursor position to the start of the current line.
func (w *VT100Writer) EraseStartOfLine() {
	w.writeRaw([]byte{0x1b, '[', '1', 'K'})
}

// EraseEndOfLine erases from the current cursor position to the end of the current line.
func (w *VT100Writer) EraseEndOfLine() {
	w.writeRaw([]byte{0x1b, '[', 'K'})
}

// EraseLine erases the entire current line.
func (w *VT100Writer) EraseLine() {
	w.writeRaw([]byte{0x1b, '[', '2', 'K'})
}

/* Cursor */

// ShowCursor stops blinking cursor and show.
func (w *VT100Writer) ShowCursor() {
	w.writeRaw([]byte{0x1b, '[', '?', '1', '2', 'l', 0x1b, '[', '?', '2', '5', 'h'})
}

// HideCursor hides cursor.
func (w *VT100Writer) HideCursor() {
	w.writeRaw([]byte{0x1b, '[', '?', '2', '5', 'l'})
}

// CursorGoTo sets the cursor position where subsequent text will begin.
func (w *VT100Writer) CursorGoTo(row, col int) {
	if row == 0 && col == 0 {
		// If no row/column parameters are provided (ie. <ESC>[H), the cursor will move to the home position.
		w.writeRaw([]byte{0x1b, '[', 'H'})
		return
	}
	r := strconv.Itoa(row)
	c := strconv.Itoa(col)
	w.writeRaw([]byte{0x1b, '['})
	w.writeRaw([]byte(r))
	w.writeRaw([]byte{';'})
	w.writeRaw([]byte(c))
	w.writeRaw([]byte{'H'})
}

// CursorUp moves the cursor up by 'n' rows; the default count is 1.
func (w *VT100Writer) CursorUp(n int) {
	if n == 0 {
		return
	} else if n < 0 {
		w.CursorDown(-n)
		return
	}
	s := strconv.Itoa(n)
	w.writeRaw([]byte{0x1b, '['})
	w.writeRaw([]byte(s))
	w.writeRaw([]byte{'A'})
}

// CursorDown moves the cursor down by 'n' rows; the default count is 1.
func (w *VT100Writer) CursorDown(n int) {
	if n == 0 {
		return
	} else if n < 0 {
		w.CursorUp(-n)
		return
	}
	s := strconv.Itoa(n)
	w.writeRaw([]byte{0x1b, '['})
	w.writeRaw([]byte(s))
	w.writeRaw([]byte{'B'})
}

// CursorForward moves the cursor forward by 'n' columns; the default count is 1.
func (w *VT100Writer) CursorForward(n int) {
	if n == 0 {
		return
	} else if n < 0 {
		w.CursorBackward(-n)
		return
	}
	s := strconv.Itoa(n)
	w.writeRaw([]byte{0x1b, '['})
	w.writeRaw([]byte(s))
	w.writeRaw([]byte{'C'})
}

// CursorBackward moves the cursor backward by 'n' columns; the default count is 1.
func (w *VT100Writer) CursorBackward(n int) {
	if n == 0 {
		return
	} else if n < 0 {
		w.CursorForward(-n)
		return
	}
	s := strconv.Itoa(n)
	w.writeRaw([]byte{0x1b, '['})
	w.writeRaw([]byte(s))
	w.writeRaw([]byte{'D'})
}

// AskForCPR asks for a cursor position report (CPR).
func (w *VT100Writer) AskForCPR() {
	// CPR: Cursor Position Request.
	w.writeRaw([]byte{0x1b, '[', '6', 'n'})
}

// SaveCursor saves current cursor position.
func (w *VT100Writer) SaveCursor() {
	w.writeRaw([]byte{0x1b, '[', 's'})
}

// UnSaveCursor restores cursor position after a Save Cursor.
func (w *VT100Writer) UnSaveCursor() {
	w.writeRaw([]byte{0x1b, '[', 'u'})
}

/* Scrolling */

// ScrollDown scrolls display down one line.
func (w *VT100Writer) ScrollDown() {
	w.writeRaw([]byte{0x1b, 'D'})
}

// ScrollUp scroll display up one line.
func (w *VT100Writer) ScrollUp() {
	w.writeRaw([]byte{0x1b, 'M'})
}

/* Title */

// SetTitle sets a title of terminal window.
func (w *VT100Writer) SetTitle(title string) {
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

	w.writeRaw([]byte{0x1b, ']', '2', ';'})
	w.writeRaw(titleBytes)
	w.writeRaw([]byte{0x07})
}

// ClearTitle clears a title of terminal window.
func (w *VT100Writer) ClearTitle() {
	w.writeRaw([]byte{0x1b, ']', '2', ';', 0x07})
}

/* Font */

/*
// SetColor sets text and background colors. and specify whether text is bold.
func (w *VT100Writer) SetColor(currColor *color.Color) {
	w.currColor = currColor
}

// ResetColor resets the text and background color to default.
func (w *VT100Writer) ResetColor() {
	w.currColor = nil
}

/*
// SetDisplayAttributes to set VT100 display attributes.
func (w *VT100Writer) SetDisplayAttributes(fg, bg Color, attrs ...DisplayAttribute) {
	w.writeRaw([]byte{0x1b, '['}) // control sequence introducer
	defer w.writeRaw([]byte{'m'}) // final character

	var separator byte = ';'
	for i := range attrs {
		p, ok := displayAttributeParameters[attrs[i]]
		if !ok {
			continue
		}
		w.writeRaw(p)
		w.writeRaw([]byte{separator})
	}

	f, ok := foregroundANSIColors[fg]
	if !ok {
		f = foregroundANSIColors[DefaultColor]
	}
	w.writeRaw(f)
	w.writeRaw([]byte{separator})
	b, ok := backgroundANSIColors[bg]
	if !ok {
		b = backgroundANSIColors[DefaultColor]
	}
	w.writeRaw(b)
}

var displayAttributeParameters = map[DisplayAttribute][]byte{
	DisplayReset:        {'0'},
	DisplayBold:         {'1'},
	DisplayLowIntensity: {'2'},
	DisplayItalic:       {'3'},
	DisplayUnderline:    {'4'},
	DisplayBlink:        {'5'},
	DisplayRapidBlink:   {'6'},
	DisplayReverse:      {'7'},
	DisplayInvisible:    {'8'},
	DisplayCrossedOut:   {'9'},
	DisplayDefaultFont:  {'1', '0'},
}

var foregroundANSIColors = map[Color][]byte{
	DefaultColor: {'3', '9'},

	// Low intensity.
	Black:     {'3', '0'},
	DarkRed:   {'3', '1'},
	DarkGreen: {'3', '2'},
	Brown:     {'3', '3'},
	DarkBlue:  {'3', '4'},
	Purple:    {'3', '5'},
	Cyan:      {'3', '6'},
	LightGray: {'3', '7'},

	// High intensity.
	DarkGray:  {'9', '0'},
	Red:       {'9', '1'},
	Green:     {'9', '2'},
	Yellow:    {'9', '3'},
	Blue:      {'9', '4'},
	Fuchsia:   {'9', '5'},
	Turquoise: {'9', '6'},
	White:     {'9', '7'},
}

var backgroundANSIColors = map[Color][]byte{
	DefaultColor: {'4', '9'},

	// Low intensity.
	Black:     {'4', '0'},
	DarkRed:   {'4', '1'},
	DarkGreen: {'4', '2'},
	Brown:     {'4', '3'},
	DarkBlue:  {'4', '4'},
	Purple:    {'4', '5'},
	Cyan:      {'4', '6'},
	LightGray: {'4', '7'},

	// High intensity
	DarkGray:  {'1', '0', '0'},
	Red:       {'1', '0', '1'},
	Green:     {'1', '0', '2'},
	Yellow:    {'1', '0', '3'},
	Blue:      {'1', '0', '4'},
	Fuchsia:   {'1', '0', '5'},
	Turquoise: {'1', '0', '6'},
	White:     {'1', '0', '7'},
} */
