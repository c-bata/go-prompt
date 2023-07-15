package prompt

import (
	"strings"

	"github.com/elk-language/go-prompt/debug"
	istrings "github.com/elk-language/go-prompt/strings"
)

// Buffer emulates the console buffer.
type Buffer struct {
	workingLines   []string // The working lines. Similar to history
	workingIndex   int      // index of the current line
	cursorPosition istrings.RuneNumber
	cacheDocument  *Document
	lastKeyStroke  Key
}

// Text returns string of the current line.
func (b *Buffer) Text() string {
	return b.workingLines[b.workingIndex]
}

// Document method to return document instance from the current text and cursor position.
func (b *Buffer) Document() (d *Document) {
	if b.cacheDocument == nil ||
		b.cacheDocument.Text != b.Text() ||
		b.cacheDocument.cursorPosition != b.cursorPosition {
		b.cacheDocument = &Document{
			Text:           b.Text(),
			cursorPosition: b.cursorPosition,
		}
	}
	b.cacheDocument.lastKey = b.lastKeyStroke
	return b.cacheDocument
}

// DisplayCursorPosition returns the cursor position on rendered text on terminal emulators.
// So if Document is "日本(cursor)語", DisplayedCursorPosition returns 4 because '日' and '本' are double width characters.
func (b *Buffer) DisplayCursorPosition(columns istrings.Width) Position {
	return b.Document().DisplayCursorPosition(columns)
}

// InsertText insert string from current line.
func (b *Buffer) InsertText(text string, overwrite bool, moveCursor bool) {
	currentTextRunes := []rune(b.Text())
	cursor := b.cursorPosition

	if overwrite {
		overwritten := string(currentTextRunes[cursor:])
		if len(overwritten) >= int(cursor)+len(text) {
			overwritten = string(currentTextRunes[cursor : cursor+istrings.RuneCount(text)])
		}
		if i := strings.IndexAny(overwritten, "\n"); i != -1 {
			overwritten = overwritten[:i]
		}
		b.setText(string(currentTextRunes[:cursor]) + text + string(currentTextRunes[cursor+istrings.RuneCount(overwritten):]))
	} else {
		b.setText(string(currentTextRunes[:cursor]) + text + string(currentTextRunes[cursor:]))
	}

	if moveCursor {
		b.cursorPosition += istrings.RuneCount(text)
	}
}

// SetText method to set text and update cursorPosition.
// (When doing this, make sure that the cursor_position is valid for this text.
// text/cursor_position should be consistent at any time, otherwise set a Document instead.)
func (b *Buffer) setText(text string) {
	debug.Assert(b.cursorPosition <= istrings.RuneCount(text), "length of input should be shorter than cursor position")
	b.workingLines[b.workingIndex] = text
}

// Set cursor position. Return whether it changed.
func (b *Buffer) setCursorPosition(p istrings.RuneNumber) {
	if p > 0 {
		b.cursorPosition = p
	} else {
		b.cursorPosition = 0
	}
}

func (b *Buffer) setDocument(d *Document) {
	b.cacheDocument = d
	b.setCursorPosition(d.cursorPosition) // Call before setText because setText check the relation between cursorPosition and line length.
	b.setText(d.Text)
}

// CursorLeft move to left on the current line.
func (b *Buffer) CursorLeft(count istrings.RuneNumber) {
	l := b.Document().GetCursorLeftPosition(count)
	b.cursorPosition += l
}

// CursorRight move to right on the current line.
func (b *Buffer) CursorRight(count istrings.RuneNumber) {
	l := b.Document().GetCursorRightPosition(count)
	b.cursorPosition += l
}

// CursorUp move cursor to the previous line.
// (for multi-line edit).
func (b *Buffer) CursorUp(count int) {
	orig := b.Document().CursorPositionCol()
	b.cursorPosition += b.Document().GetCursorUpPosition(count, orig)
}

// CursorDown move cursor to the next line.
// (for multi-line edit).
func (b *Buffer) CursorDown(count int) {
	orig := b.Document().CursorPositionCol()
	b.cursorPosition += b.Document().GetCursorDownPosition(count, orig)
}

// DeleteBeforeCursor delete specified number of characters before cursor and return the deleted text.
func (b *Buffer) DeleteBeforeCursor(count istrings.RuneNumber) (deleted string) {
	debug.Assert(count >= 0, "count should be positive")
	r := []rune(b.Text())

	if b.cursorPosition > 0 {
		start := b.cursorPosition - count
		if start < 0 {
			start = 0
		}
		deleted = string(r[start:b.cursorPosition])
		b.setDocument(&Document{
			Text:           string(r[:start]) + string(r[b.cursorPosition:]),
			cursorPosition: b.cursorPosition - istrings.RuneNumber(len([]rune(deleted))),
		})
	}
	return
}

// NewLine means CR.
func (b *Buffer) NewLine(copyMargin bool) {
	if copyMargin {
		b.InsertText("\n"+b.Document().leadingWhitespaceInCurrentLine(), false, true)
	} else {
		b.InsertText("\n", false, true)
	}
}

// Delete specified number of characters and Return the deleted text.
func (b *Buffer) Delete(count istrings.RuneNumber) string {
	r := []rune(b.Text())
	if b.cursorPosition < istrings.RuneNumber(len(r)) {
		textAfterCursor := b.Document().TextAfterCursor()
		textAfterCursorRunes := []rune(textAfterCursor)
		deletedRunes := textAfterCursorRunes[:count]
		b.setText(string(r[:b.cursorPosition]) + string(r[b.cursorPosition+istrings.RuneNumber(len(deletedRunes)):]))

		deleted := string(deletedRunes)
		return deleted
	}

	return ""
}

// JoinNextLine joins the next line to the current one by deleting the line ending after the current line.
func (b *Buffer) JoinNextLine(separator string) {
	if !b.Document().OnLastLine() {
		b.cursorPosition += b.Document().GetEndOfLinePosition()
		b.Delete(1)
		// Remove spaces
		b.setText(b.Document().TextBeforeCursor() + separator + strings.TrimLeft(b.Document().TextAfterCursor(), " "))
	}
}

// SwapCharactersBeforeCursor swaps the last two characters before the cursor.
func (b *Buffer) SwapCharactersBeforeCursor() {
	if b.cursorPosition >= 2 {
		x := b.Text()[b.cursorPosition-2 : b.cursorPosition-1]
		y := b.Text()[b.cursorPosition-1 : b.cursorPosition]
		b.setText(b.Text()[:b.cursorPosition-2] + y + x + b.Text()[b.cursorPosition:])
	}
}

// NewBuffer is constructor of Buffer struct.
func NewBuffer() (b *Buffer) {
	b = &Buffer{
		workingLines: []string{""},
		workingIndex: 0,
	}
	return
}
