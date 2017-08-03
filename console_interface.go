package prompt

type WinSize struct {
	Row uint16
	Col uint16
}

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

type ConsoleParser interface {
	// Setup
	Setup() error
	// TearDown
	TearDown() error
	// GetSCIICode returns ASCIICode correspond to input byte codes.
	GetKey(b []byte) Key
	// GetWinSize returns winsize struct which is the response of ioctl(2).
	GetWinSize() *WinSize
}

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
