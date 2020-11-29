package prompt

import (
	"sync"

	fcolor "github.com/fatih/color"
)

var (
	consoleWriterMu sync.Mutex
	consoleWriter   ConsoleWriter
)

func registerConsoleWriter(f ConsoleWriter) {
	consoleWriterMu.Lock()
	defer consoleWriterMu.Unlock()
	consoleWriter = f
}

// ConsoleWriter is an interface to abstract output layer.
type ConsoleWriter interface {
	// WriteStr to write safety string by removing control sequences.
	WriteStr(data string, color *fcolor.Color)
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
}
