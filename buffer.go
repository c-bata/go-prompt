package prompt

import (
	"strings"

	"github.com/elk-language/go-prompt/debug"
	istrings "github.com/elk-language/go-prompt/strings"
	"golang.org/x/exp/utf8string"
)

// Buffer emulates the console buffer.
type Buffer struct {
	workingLines    []string // The working lines. Similar to history
	workingIndex    int      // index of the current line
	startLine       int      // Line number of the first visible line in the terminal (0-indexed)
	cursorPosition  istrings.RuneNumber
	cacheDocument   *Document
	preferredColumn istrings.Width // Remember the original column for the next up/down movement.
	lastKeyStroke   Key
}

// Text returns string of the current line.
func (b *Buffer) Text() string {
	return b.workingLines[b.workingIndex]
}

func (b *Buffer) resetPreferredColumn() {
	b.preferredColumn = -1
}

func (b *Buffer) updatePreferredColumn() {
	b.preferredColumn = b.Document().CursorPositionCol()
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

// Insert string into the buffer and move the cursor.
func (b *Buffer) InsertTextMoveCursor(text string, columns istrings.Width, rows int, overwrite bool) {
	b.insertText(text, columns, rows, overwrite, true)
}

// Insert string into the buffer without moving the cursor.
func (b *Buffer) InsertText(text string, overwrite bool) {
	b.insertText(text, 0, 0, overwrite, false)
}

// insertText insert string from current line.
func (b *Buffer) insertText(text string, columns istrings.Width, rows int, overwrite bool, moveCursor bool) {
	currentTextRunes := []rune(b.Text())
	cursor := b.cursorPosition

	if overwrite {
		overwritten := string(currentTextRunes[cursor:])
		if len(overwritten) >= int(cursor)+len(text) {
			overwritten = string(currentTextRunes[cursor : cursor+istrings.RuneCountInString(text)])
		}
		if i := strings.IndexAny(overwritten, "\n"); i != -1 {
			overwritten = overwritten[:i]
		}
		b.setText(
			string(currentTextRunes[:cursor])+text+string(currentTextRunes[cursor+istrings.RuneCountInString(overwritten):]),
			columns,
			rows,
		)
	} else {
		b.setText(
			string(currentTextRunes[:cursor])+text+string(currentTextRunes[cursor:]),
			columns,
			rows,
		)
	}

	if moveCursor {
		b.cursorPosition += istrings.RuneCountInString(text)
		b.recalculateStartLine(columns, rows)
		b.updatePreferredColumn()
	}
}

func (b *Buffer) resetStartLine() {
	b.startLine = 0
}

// Calculates the startLine once again and returns true when it's been changed.
func (b *Buffer) recalculateStartLine(columns istrings.Width, rows int) bool {
	origStartLine := b.startLine
	pos := b.DisplayCursorPosition(columns)
	if pos.Y > b.startLine+rows-1 {
		b.startLine = pos.Y - rows + 1
	} else if pos.Y < b.startLine {
		b.startLine = pos.Y
	}

	if b.startLine < 0 {
		b.startLine = 0
	}
	return origStartLine != b.startLine
}

// SetText method to set text and update cursorPosition.
// (When doing this, make sure that the cursor_position is valid for this text.
// text/cursor_position should be consistent at any time, otherwise set a Document instead.)
func (b *Buffer) setText(text string, col istrings.Width, row int) {
	debug.Assert(b.cursorPosition <= istrings.RuneCountInString(text), "length of input should be shorter than cursor position")
	b.workingLines[b.workingIndex] = text
	b.recalculateStartLine(col, row)
	b.resetPreferredColumn()
}

// Set cursor position. Return whether it changed.
func (b *Buffer) setCursorPosition(p istrings.RuneNumber) {
	if p > 0 {
		b.cursorPosition = p
	} else {
		b.cursorPosition = 0
	}
}

func (b *Buffer) setDocument(d *Document, columns istrings.Width, rows int) {
	b.cacheDocument = d
	b.setCursorPosition(d.cursorPosition) // Call before setText because setText check the relation between cursorPosition and line length.
	b.setText(d.Text, columns, rows)
	b.recalculateStartLine(columns, rows)
	b.resetPreferredColumn()
}

// Move to the left on the current line by the given amount of graphemes.
// Returns true when the view should be rerendered.
func (b *Buffer) CursorLeft(count istrings.GraphemeNumber, columns istrings.Width, rows int) bool {
	return b.cursorHorizontalMove(b.Document().GetCursorLeftPosition(count), columns, rows)
}

// Move to the left on the current line by the given amount of runes.
// Returns true when the view should be rerendered.
func (b *Buffer) CursorLeftRunes(count istrings.RuneNumber, columns istrings.Width, rows int) bool {
	return b.cursorHorizontalMove(b.Document().GetCursorLeftPositionRunes(count), columns, rows)
}

// Move to the left on the current line by the given amount of bytes.
// Returns true when the view should be rerendered.
func (b *Buffer) CursorLeftBytes(count istrings.ByteNumber, columns istrings.Width, rows int) bool {
	return b.cursorHorizontalMove(b.Document().GetCursorLeftPositionBytes(count), columns, rows)
}

// Move to the right on the current line by the given amount of graphemes.
// Returns true when the view should be rerendered.
func (b *Buffer) CursorRight(count istrings.GraphemeNumber, columns istrings.Width, rows int) bool {
	return b.cursorHorizontalMove(b.Document().GetCursorRightPosition(count), columns, rows)
}

// Move to the right on the current line by the given amount of runes.
// Returns true when the view should be rerendered.
func (b *Buffer) CursorRightRunes(count istrings.RuneNumber, columns istrings.Width, rows int) bool {
	return b.cursorHorizontalMove(b.Document().GetCursorRightPositionRunes(count), columns, rows)
}

// Move to the right on the current line by the given amount of bytes.
// Returns true when the view should be rerendered.
func (b *Buffer) CursorRightBytes(count istrings.ByteNumber, columns istrings.Width, rows int) bool {
	return b.cursorHorizontalMove(b.Document().GetCursorRightPositionBytes(count), columns, rows)
}

func (b *Buffer) cursorHorizontalMove(count istrings.RuneNumber, columns istrings.Width, rows int) bool {
	b.cursorPosition += count
	b.updatePreferredColumn()
	return b.recalculateStartLine(columns, rows)
}

// CursorUp move cursor to the previous line.
// (for multi-line edit).
// Returns true when the view should be rerendered.
func (b *Buffer) CursorUp(count int, columns istrings.Width, rows int) bool {
	b.cursorPosition += b.Document().GetCursorUpPosition(count, b.preferredColumn)
	return b.recalculateStartLine(columns, rows)
}

// CursorDown move cursor to the next line.
// (for multi-line edit).
// Returns true when the view should be rerendered.
func (b *Buffer) CursorDown(count int, columns istrings.Width, rows int) bool {
	b.cursorPosition += b.Document().GetCursorDownPosition(count, b.preferredColumn)
	return b.recalculateStartLine(columns, rows)
}

// Deletes the specified number of graphemes before the cursor and returns the deleted text.
func (b *Buffer) DeleteBeforeCursor(count istrings.GraphemeNumber, columns istrings.Width, rows int) string {
	debug.Assert(count >= 0, "count should be positive")
	if b.cursorPosition < 0 {
		return ""
	}

	var deleted string

	textUtf8 := utf8string.NewString(b.Text())
	textBeforeCursor := textUtf8.Slice(0, int(b.cursorPosition))
	graphemeLength := istrings.GraphemeCountInString(textBeforeCursor)

	start := istrings.RuneIndexNthGrapheme(textBeforeCursor, graphemeLength-count)
	if start < 0 {
		start = 0
	}
	deleted = textUtf8.Slice(int(start), int(b.cursorPosition))
	b.setDocument(
		&Document{
			Text:           textUtf8.Slice(0, int(start)) + textUtf8.Slice(int(b.cursorPosition), textUtf8.RuneCount()),
			cursorPosition: b.cursorPosition - istrings.RuneCountInString(deleted),
		},
		columns,
		rows,
	)

	b.recalculateStartLine(columns, rows)
	b.updatePreferredColumn()
	return deleted
}

// Deletes the specified number of runes before the cursor and returns the deleted text.
func (b *Buffer) DeleteBeforeCursorRunes(count istrings.RuneNumber, columns istrings.Width, rows int) (deleted string) {
	debug.Assert(count >= 0, "count should be positive")
	if b.cursorPosition <= 0 {
		return ""
	}
	r := []rune(b.Text())

	start := b.cursorPosition - count
	if start < 0 {
		start = 0
	}
	deleted = string(r[start:b.cursorPosition])
	b.setDocument(
		&Document{
			Text:           string(r[:start]) + string(r[b.cursorPosition:]),
			cursorPosition: b.cursorPosition - istrings.RuneNumber(len([]rune(deleted))),
		},
		columns,
		rows,
	)
	b.recalculateStartLine(columns, rows)
	b.updatePreferredColumn()
	return
}

// NewLine means CR.
func (b *Buffer) NewLine(columns istrings.Width, rows int, copyMargin bool) {
	if copyMargin {
		b.InsertTextMoveCursor("\n"+b.Document().leadingWhitespaceInCurrentLine(), columns, rows, false)
	} else {
		b.InsertTextMoveCursor("\n", columns, rows, false)
	}
}

// Deletes the specified number of graphemes and returns the deleted text.
func (b *Buffer) Delete(count istrings.GraphemeNumber, col istrings.Width, row int) string {
	textUtf8 := utf8string.NewString(b.Text())
	if b.cursorPosition >= istrings.RuneCountInString(b.Text()) {
		return ""
	}

	textAfterCursor := b.Document().TextAfterCursor()
	textAfterCursorUtf8 := utf8string.NewString(textAfterCursor)

	deletedRunes := textAfterCursorUtf8.Slice(0, int(istrings.RuneIndexNthGrapheme(textAfterCursor, count)))

	b.setText(
		textUtf8.Slice(0, int(b.cursorPosition))+textUtf8.Slice(int(b.cursorPosition)+int(istrings.RuneCountInString(deletedRunes)), textUtf8.RuneCount()),
		col,
		row,
	)

	deleted := string(deletedRunes)
	return deleted
}

// Deletes the specified number of runes and returns the deleted text.
func (b *Buffer) DeleteRunes(count istrings.RuneNumber, col istrings.Width, row int) string {
	r := []rune(b.Text())
	if b.cursorPosition < istrings.RuneNumber(len(r)) {
		textAfterCursor := b.Document().TextAfterCursor()
		textAfterCursorRunes := []rune(textAfterCursor)
		deletedRunes := textAfterCursorRunes[:count]
		b.setText(
			string(r[:b.cursorPosition])+string(r[b.cursorPosition+istrings.RuneNumber(len(deletedRunes)):]),
			col,
			row,
		)

		deleted := string(deletedRunes)
		return deleted
	}

	return ""
}

// JoinNextLine joins the next line to the current one by deleting the line ending after the current line.
func (b *Buffer) JoinNextLine(separator string, col istrings.Width, row int) {
	if !b.Document().OnLastLine() {
		b.cursorPosition += b.Document().GetEndOfLinePosition()
		b.Delete(1, col, row)
		// Remove spaces
		b.setText(
			b.Document().TextBeforeCursor()+separator+strings.TrimLeft(b.Document().TextAfterCursor(), " "),
			col,
			row,
		)
	}
}

// SwapCharactersBeforeCursor swaps the last two characters before the cursor.
func (b *Buffer) SwapCharactersBeforeCursor(col istrings.Width, row int) {
	if b.cursorPosition >= 2 {
		x := b.Text()[b.cursorPosition-2 : b.cursorPosition-1]
		y := b.Text()[b.cursorPosition-1 : b.cursorPosition]
		b.setText(
			b.Text()[:b.cursorPosition-2]+y+x+b.Text()[b.cursorPosition:],
			col,
			row,
		)
	}
}

// NewBuffer is constructor of Buffer struct.
func NewBuffer() (b *Buffer) {
	b = &Buffer{
		workingLines:    []string{""},
		workingIndex:    0,
		startLine:       0,
		preferredColumn: -1,
	}
	return
}
