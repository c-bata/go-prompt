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

	WriteRaw(data []byte)
	Write(data []byte)
	WriteStr(data string)
	WriteRawStr(data string)
	Flush() error

	/* Erasing */

	EraseScreen()
	EraseUp()
	EraseDown()
	EraseStartOfLine()
	EraseEndOfLine()
	EraseLine()

	/* Cursor */

	ShowCursor()
	HideCursor()
	CursorGoTo(row, col int)
	CursorUp(n int)
	CursorDown(n int)
	CursorForward(n int)
	CursorBackward(n int)
	AskForCPR()
	SaveCursor()
	UnSaveCursor()

	/* Scrolling */

	ScrollDown()
	ScrollUp()

	/* Title */

	SetTitle(title string)
	ClearTitle()

	/* Font */

	SetColor(fg, bg Color, bold bool)
}
