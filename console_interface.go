package prompt

// WinSize represents the width and height of terminal.
type WinSize struct {
	Row uint16
	Col uint16
}

// Color represents color on terminal.
type Color int

const (
	DefaultColor Color = iota

	// Low intensity
	Black
	DarkRed
	DarkGreen
	Brown
	DarkBlue
	Purple
	Cyan
	LightGray

	// High intensity
	DarkGray
	Red
	Green
	Yellow
	Blue
	Fuchsia
	Turquoise
	White
)

// ConsoleParser is an interface to abstract input layer.
type ConsoleParser interface {
	// Setup should be called before starting input
	Setup() error
	// TearDown should be called after stopping input
	TearDown() error
	// GetKey returns Key correspond to input byte codes.
	GetKey(b []byte) Key
	// GetWinSize returns WinSize object to represent width and height of terminal.
	GetWinSize() *WinSize
	// Read returns byte array.
	Read() ([]byte, error)
}

// ConsoleWriter is an interface to abstract output layer.
type ConsoleWriter interface {
	/* Write */

	// WriteRaw to write raw byte array.
	WriteRaw(data []byte)
	// Write to write byte array without control sequences.
	Write(data []byte)
	// WriteStr to write raw string.
	WriteRawStr(data string)
	// WriteStr to write string without control sequences.
	WriteStr(data string)
	// Flush to flush buffer.
	Flush() error

	/* Erasing */

	// EraseScreen erases the screen with the background colour and moves the cursor to home.
	EraseScreen()
	// EraseUp erases the screen from the current line up to the top of the screen.
	EraseUp()
	// EraseDown erases the screen from the current line down to the bottom of the screen.
	EraseDown()
	// EraseStartOfLine erases from the current cursor position to the start of the current line.
	EraseStartOfLine()
	// EraseEndOfLine erases from the current cursor position to the end of the current line.
	EraseEndOfLine()
	// EraseLine erases the entire current line.
	EraseLine()

	/* Cursor */

	// ShowCursor stops blinking cursor and show.
	ShowCursor()
	// HideCursor hides cursor.
	HideCursor()
	// CursorGoTo sets the cursor position where subsequent text will begin.
	CursorGoTo(row, col int)
	// CursorUp moves the cursor up by 'n' rows; the default count is 1.
	CursorUp(n int)
	// CursorDown moves the cursor down by 'n' rows; the default count is 1.
	CursorDown(n int)
	// CursorForward moves the cursor forward by 'n' columns; the default count is 1.
	CursorForward(n int)
	// CursorBackward moves the cursor backward by 'n' columns; the default count is 1.
	CursorBackward(n int)
	// AskForCPR asks for a cursor position report (CPR).
	AskForCPR()
	// SaveCursor saves current cursor position.
	SaveCursor()
	// UnSaveCursor restores cursor position after a Save Cursor.
	UnSaveCursor()

	/* Scrolling */

	// ScrollDown scrolls display down one line.
	ScrollDown()
	// ScrollUp scroll display up one line.
	ScrollUp()

	/* Title */

	// SetTitle sets a title of terminal window.
	SetTitle(title string)
	// ClearTitle clears a title of terminal window.
	ClearTitle()

	/* Font */

	// SetColor sets text and background colors. and specify whether text is bold.
	SetColor(fg, bg Color, bold bool)
}
