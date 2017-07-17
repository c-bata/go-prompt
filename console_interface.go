package prompt

type WinSize struct {
	Row uint16
	Col uint16
}

type ConsoleParser interface {
	// Setup
	Setup() error
	// TearDown
	TearDown() error
	// GetSCIICode returns ASCIICode correspond to input byte codes.
	GetASCIICode(b []byte) *ASCIICode
	// GetWinSize returns winsize struct which is the response of ioctl(2).
	GetWinSize() *WinSize
}

type ConsoleWriter interface {
	/* Write */

	Write(data []byte)
	WriteStr(data string)
	WriteRaw(data []byte)
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

	/* colors */

	SetColor(fg, bg Color) (ok bool)
}
