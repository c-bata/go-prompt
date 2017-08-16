package prompt

import (
	"strings"
	"unicode/utf8"
)

// Document has text displayed in terminal and cursor position.
type Document struct {
	Text           string
	CursorPosition int
}

// NewDocument return the new empty document.
func NewDocument() *Document {
	return &Document{
		Text:           "",
		CursorPosition: 0,
	}
}

// GetCharRelativeToCursor return character relative to cursor position, or empty string
func (d *Document) GetCharRelativeToCursor(offset int) (r rune) {
	s := d.Text
	cnt := 0

	for len(s) > 0 {
		cnt++
		r, size := utf8.DecodeRuneInString(s)
		if cnt == d.CursorPosition+offset {
			return r
		}
		s = s[size:]
	}
	return 0
}

// TextBeforeCursor returns the text before the cursor.
func (d *Document) TextBeforeCursor() string {
	r := []rune(d.Text)
	return string(r[:d.CursorPosition])
}

// TextAfterCursor returns the text after the cursor.
func (d *Document) TextAfterCursor() string {
	r := []rune(d.Text)
	return string(r[d.CursorPosition:])
}

// GetWordBeforeCursor returns the word before the cursor.
// If we have whitespace before the cursor this returns an empty string.
func (d *Document) GetWordBeforeCursor() string {
	x := d.TextBeforeCursor()
	return x[d.FindStartOfPreviousWord():]
}

// GetWordBeforeCursorWithSpace returns the word before the cursor.
// Unlike GetWordBeforeCursor, it returns string containing space
func (d *Document) GetWordBeforeCursorWithSpace() string {
	x := d.TextBeforeCursor()
	return x[d.FindStartOfPreviousWordWithSpace():]
}

// FindStartOfPreviousWord returns an index relative to the cursor position
// pointing to the start of the previous word. Return `None` if nothing was found.
func (d *Document) FindStartOfPreviousWord() int {
	// Reverse the text before the cursor, in order to do an efficient backwards search.
	x := d.TextBeforeCursor()
	l := len(x)
	for i := l; i > 0; i-- {
		if x[i-1:i] == " " {
			return i
		}
	}
	return 0
}

// FindStartOfPreviousWordWithSpace is almost the same as FindStartOfPreviousWord.
// The only difference is to ignore contiguous spaces.
func (d *Document) FindStartOfPreviousWordWithSpace() int {
	// Reverse the text before the cursor, in order to do an efficient backwards search.
	x := d.TextBeforeCursor()
	l := len(x)
	appear := false
	for i := l; i > 0; i-- {
		if x[i-1:i] != " " {
			appear = true
		}
		if x[i-1:i] == " " && appear {
			return i
		}
	}
	return 0
}

// CurrentLineBeforeCursor returns the text from the start of the line until the cursor.
func (d *Document) CurrentLineBeforeCursor() string {
	s := strings.Split(d.TextBeforeCursor(), "\n")
	return s[len(s)-1]
}

// CurrentLineAfterCursor returns the text from the cursor until the end of the line.
func (d *Document) CurrentLineAfterCursor() string {
	return strings.Split(d.TextAfterCursor(), "\n")[0]
}

// CurrentLine return the text on the line where the cursor is. (when the input
// consists of just one line, it equals `text`.
func (d *Document) CurrentLine() string {
	return d.CurrentLineBeforeCursor() + d.CurrentLineAfterCursor()
}

// Array pointing to the start indexes of all the lines.
func (d *Document) lineStartIndexes() []int {
	// TODO: Cache, because this is often reused.
	// (If it is used, it's often used many times.
	// And this has to be fast for editing big documents!)
	lc := d.LineCount()
	lengths := make([]int, lc)
	for i, l := range d.Lines() {
		lengths[i] = len(l)
	}

	// Calculate cumulative sums.
	indexes := make([]int, lc+1)
	indexes[0] = 0 // https://github.com/jonathanslenders/python-prompt-toolkit/blob/master/prompt_toolkit/document.py#L189
	pos := 0
	for i, l := range lengths {
		pos += l + 1
		indexes[i+1] = pos
	}
	if lc > 1 {
		// Pop the last item. (This is not a new line.)
		indexes = indexes[:lc]
	}
	return indexes
}

// For the index of a character at a certain line, calculate the index of
// the first character on that line.
func (d *Document) findLineStartIndex(index int) (pos int, lineStartIndex int) {
	indexes := d.lineStartIndexes()
	pos = BisectRight(indexes, index) - 1
	lineStartIndex = indexes[pos]
	return
}

// CursorPositionRow returns the current row. (0-based.)
func (d *Document) CursorPositionRow() (row int) {
	row, _ = d.findLineStartIndex(d.CursorPosition)
	return
}

// CursorPositionCol returns the current column. (0-based.)
func (d *Document) CursorPositionCol() (col int) {
	// Don't use self.text_before_cursor to calculate this. Creating substrings
	// and splitting is too expensive for getting the cursor position.
	_, index := d.findLineStartIndex(d.CursorPosition)
	col = d.CursorPosition - index
	return
}

// GetCursorLeftPosition returns the relative position for cursor left.
func (d *Document) GetCursorLeftPosition(count int) int {
	if count < 0 {
		return d.GetCursorRightPosition(-count)
	}
	if d.CursorPositionCol() > count {
		return -count
	}
	return -d.CursorPositionCol()
}

// GetCursorRightPosition returns relative position for cursor right.
func (d *Document) GetCursorRightPosition(count int) int {
	if count < 0 {
		return d.GetCursorLeftPosition(-count)
	}
	if len(d.CurrentLineAfterCursor()) > count {
		return count
	}
	return len(d.CurrentLineAfterCursor())
}

// GetCursorUpPosition return the relative cursor position (character index) where we would be
// if the user pressed the arrow-up button.
func (d *Document) GetCursorUpPosition(count int, preferredColumn int) int {
	var col int
	if preferredColumn == -1 { // -1 means nil
		col = d.CursorPositionCol()
	} else {
		col = preferredColumn
	}

	row := d.CursorPositionRow() - count
	if row < 0 {
		row = 0
	}
	return d.TranslateRowColToIndex(row, col) - d.CursorPosition
}

// GetCursorDownPosition return the relative cursor position (character index) where we would be if the
// user pressed the arrow-down button.
func (d *Document) GetCursorDownPosition(count int, preferredColumn int) int {
	var col int
	if preferredColumn == -1 { // -1 means nil
		col = d.CursorPositionCol()
	} else {
		col = preferredColumn
	}
	row := d.CursorPositionRow() + count
	return d.TranslateRowColToIndex(row, col) - d.CursorPosition
}

// Lines returns the array of all the lines.
func (d *Document) Lines() []string {
	// TODO: Cache, because this one is reused very often.
	return strings.Split(d.Text, "\n")
}

// LineCount return the number of lines in this document. If the document ends
// with a trailing \n, that counts as the beginning of a new line.
func (d *Document) LineCount() int {
	return len(d.Lines())
}

// TranslateIndexToPosition given an index for the text, return the corresponding (row, col) tuple.
// (0-based. Returns (0, 0) for index=0.)
func (d *Document) TranslateIndexToPosition(index int) (row int, col int) {
	row, rowIndex := d.findLineStartIndex(index)
	col = index - rowIndex
	return
}

// TranslateRowColToIndex given a (row, col), return the corresponding index.
// (Row and col params are 0-based.)
func (d *Document) TranslateRowColToIndex(row int, column int) (index int) {
	indexes := d.lineStartIndexes()
	if row < 0 {
		row = 0
	} else if row > len(indexes) {
		row = len(indexes) - 1
	}
	index = indexes[row]
	line := d.Lines()[row]

	// python) result += max(0, min(col, len(line)))
	if column > 0 || len(line) > 0 {
		if column > len(line) {
			index += len(line)
		} else {
			index += column
		}
	}

	// Keep in range. (len(self.text) is included, because the cursor can be
	// right after the end of the text as well.)
	// python) result = max(0, min(result, len(self.text)))
	if index > len(d.Text) {
		index = len(d.Text)
	}
	if index < 0 {
		index = 0
	}
	return index
}

// OnLastLine returns true when we are at the last line.
func (d *Document) OnLastLine() bool {
	return d.CursorPositionRow() == (d.LineCount() - 1)
}

// GetEndOfLinePosition returns relative position for the end of this line.
func (d *Document) GetEndOfLinePosition() int {
	return len([]rune(d.CurrentLineAfterCursor()))
}

func (d *Document) leadingWhitespaceInCurrentLine() (margin string) {
	trimmed := strings.TrimSpace(d.CurrentLine())
	margin = d.CurrentLine()[:len(d.CurrentLine())-len(trimmed)]
	return
}
