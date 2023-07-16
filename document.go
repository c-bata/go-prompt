package prompt

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/elk-language/go-prompt/bisect"
	istrings "github.com/elk-language/go-prompt/strings"
	"golang.org/x/exp/utf8string"
)

// Document has text displayed in terminal and cursor position.
type Document struct {
	Text string
	// This represents a index in a rune array of Document.Text.
	// So if Document is "日本(cursor)語", cursorPosition is 2.
	// But DisplayedCursorPosition returns 4 because '日' and '本' are double width characters.
	cursorPosition istrings.RuneNumber
	lastKey        Key
}

// NewDocument return the new empty document.
func NewDocument() *Document {
	return &Document{
		Text:           "",
		cursorPosition: 0,
	}
}

// LastKeyStroke return the last key pressed in this document.
func (d *Document) LastKeyStroke() Key {
	return d.lastKey
}

// DisplayCursorPosition returns the cursor position on rendered text on terminal emulators.
// So if Document is "日本(cursor)語", DisplayedCursorPosition returns 4 because '日' and '本' are double width characters.
func (d *Document) DisplayCursorPosition(columns istrings.Width) Position {
	str := utf8string.NewString(d.Text).Slice(0, int(d.cursorPosition))
	return positionAtEndOfString(str, columns)
}

// GetCharRelativeToCursor return character relative to cursor position, or empty string
func (d *Document) GetCharRelativeToCursor(offset istrings.RuneNumber) (r rune) {
	s := d.Text
	var cnt istrings.RuneNumber

	for len(s) > 0 {
		cnt++
		r, size := utf8.DecodeRuneInString(s)
		if cnt == d.cursorPosition+istrings.RuneNumber(offset) {
			return r
		}
		s = s[size:]
	}
	return 0
}

// TextBeforeCursor returns the text before the cursor.
func (d *Document) TextBeforeCursor() string {
	r := []rune(d.Text)
	return string(r[:d.cursorPosition])
}

// TextAfterCursor returns the text after the cursor.
func (d *Document) TextAfterCursor() string {
	r := []rune(d.Text)
	return string(r[d.cursorPosition:])
}

// GetWordBeforeCursor returns the word before the cursor.
// If we have whitespace before the cursor this returns an empty string.
func (d *Document) GetWordBeforeCursor() string {
	x := d.TextBeforeCursor()
	return x[d.FindStartOfPreviousWord():]
}

// GetWordAfterCursor returns the word after the cursor.
// If we have whitespace after the cursor this returns an empty string.
func (d *Document) GetWordAfterCursor() string {
	x := d.TextAfterCursor()
	return x[:d.FindEndOfCurrentWord()]
}

// GetWordBeforeCursorWithSpace returns the word before the cursor.
// Unlike GetWordBeforeCursor, it returns string containing space
func (d *Document) GetWordBeforeCursorWithSpace() string {
	x := d.TextBeforeCursor()
	return x[d.FindStartOfPreviousWordWithSpace():]
}

// GetWordAfterCursorWithSpace returns the word after the cursor.
// Unlike GetWordAfterCursor, it returns string containing space
func (d *Document) GetWordAfterCursorWithSpace() string {
	x := d.TextAfterCursor()
	return x[:d.FindEndOfCurrentWordWithSpace()]
}

// GetWordBeforeCursorUntilSeparator returns the text before the cursor until next separator.
func (d *Document) GetWordBeforeCursorUntilSeparator(sep string) string {
	x := d.TextBeforeCursor()
	return x[d.FindStartOfPreviousWordUntilSeparator(sep):]
}

// GetWordAfterCursorUntilSeparator returns the text after the cursor until next separator.
func (d *Document) GetWordAfterCursorUntilSeparator(sep string) string {
	x := d.TextAfterCursor()
	return x[:d.FindEndOfCurrentWordUntilSeparator(sep)]
}

// GetWordBeforeCursorUntilSeparatorIgnoreNextToCursor returns the word before the cursor.
// Unlike GetWordBeforeCursor, it returns string containing space
func (d *Document) GetWordBeforeCursorUntilSeparatorIgnoreNextToCursor(sep string) string {
	x := d.TextBeforeCursor()
	return x[d.FindStartOfPreviousWordUntilSeparatorIgnoreNextToCursor(sep):]
}

// GetWordAfterCursorUntilSeparatorIgnoreNextToCursor returns the word after the cursor.
// Unlike GetWordAfterCursor, it returns string containing space
func (d *Document) GetWordAfterCursorUntilSeparatorIgnoreNextToCursor(sep string) string {
	x := d.TextAfterCursor()
	return x[:d.FindEndOfCurrentWordUntilSeparatorIgnoreNextToCursor(sep)]
}

// FindStartOfPreviousWord returns an index relative to the cursor position
// pointing to the start of the previous word. Return 0 if nothing was found.
func (d *Document) FindStartOfPreviousWord() istrings.ByteNumber {
	x := d.TextBeforeCursor()
	i := istrings.ByteNumber(strings.LastIndexAny(x, " \n"))
	if i != -1 {
		return i + 1
	}
	return 0
}

// Returns the rune count
// of the text before the cursor until the start of the previous word.
func (d *Document) FindRuneNumberUntilStartOfPreviousWord() istrings.RuneNumber {
	x := d.TextBeforeCursor()
	return istrings.RuneCount(x[d.FindStartOfPreviousWordWithSpace():])
}

// FindStartOfPreviousWordWithSpace is almost the same as FindStartOfPreviousWord.
// The only difference is to ignore contiguous spaces.
func (d *Document) FindStartOfPreviousWordWithSpace() istrings.ByteNumber {
	x := d.TextBeforeCursor()
	end := istrings.LastIndexNotByte(x, ' ')
	if end == -1 {
		return 0
	}

	start := istrings.ByteNumber(strings.LastIndexByte(x[:end], ' '))
	if start == -1 {
		return 0
	}
	return start + 1
}

// FindStartOfPreviousWordUntilSeparator is almost the same as FindStartOfPreviousWord.
// But this can specify Separator. Return 0 if nothing was found.
func (d *Document) FindStartOfPreviousWordUntilSeparator(sep string) istrings.ByteNumber {
	if sep == "" {
		return d.FindStartOfPreviousWord()
	}

	x := d.TextBeforeCursor()
	i := istrings.ByteNumber(strings.LastIndexAny(x, sep))
	if i != -1 {
		return i + 1
	}
	return 0
}

// FindStartOfPreviousWordUntilSeparatorIgnoreNextToCursor is almost the same as FindStartOfPreviousWordWithSpace.
// But this can specify Separator. Return 0 if nothing was found.
func (d *Document) FindStartOfPreviousWordUntilSeparatorIgnoreNextToCursor(sep string) istrings.ByteNumber {
	if sep == "" {
		return d.FindStartOfPreviousWordWithSpace()
	}

	x := d.TextBeforeCursor()
	end := istrings.LastIndexNotAny(x, sep)
	if end == -1 {
		return 0
	}
	start := istrings.ByteNumber(strings.LastIndexAny(x[:end], sep))
	if start == -1 {
		return 0
	}
	return start + 1
}

// FindEndOfCurrentWord returns a byte index relative to the cursor position.
// pointing to the end of the current word. Return 0 if nothing was found.
func (d *Document) FindEndOfCurrentWord() istrings.ByteNumber {
	x := d.TextAfterCursor()
	i := istrings.ByteNumber(strings.IndexByte(x, ' '))
	if i != -1 {
		return i
	}
	return istrings.ByteNumber(len(x))
}

// FindEndOfCurrentWordWithSpace is almost the same as FindEndOfCurrentWord.
// The only difference is to ignore contiguous spaces.
func (d *Document) FindEndOfCurrentWordWithSpace() istrings.ByteNumber {
	x := d.TextAfterCursor()

	start := istrings.IndexNotByte(x, ' ')
	if start == -1 {
		return istrings.ByteNumber(len(x))
	}

	end := istrings.ByteNumber(strings.IndexByte(x[start:], ' '))
	if end == -1 {
		return istrings.ByteNumber(len(x))
	}

	return start + end
}

// Returns the number of runes
// of the text after the cursor until the end of the current word.
func (d *Document) FindRuneNumberUntilEndOfCurrentWord() istrings.RuneNumber {
	t := d.TextAfterCursor()
	var count istrings.RuneNumber
	nonSpaceCharSeen := false
	for _, char := range t {
		if !nonSpaceCharSeen && char == ' ' {
			count += 1
			continue
		}

		if nonSpaceCharSeen && char == ' ' {
			break
		}

		nonSpaceCharSeen = true
		count += 1
	}

	return count
}

// FindEndOfCurrentWordUntilSeparator is almost the same as FindEndOfCurrentWord.
// But this can specify Separator. Return 0 if nothing was found.
func (d *Document) FindEndOfCurrentWordUntilSeparator(sep string) istrings.ByteNumber {
	if sep == "" {
		return d.FindEndOfCurrentWord()
	}

	x := d.TextAfterCursor()
	i := istrings.ByteNumber(strings.IndexAny(x, sep))
	if i != -1 {
		return i
	}
	return istrings.ByteNumber(len(x))
}

// FindEndOfCurrentWordUntilSeparatorIgnoreNextToCursor is almost the same as FindEndOfCurrentWordWithSpace.
// But this can specify Separator. Return 0 if nothing was found.
func (d *Document) FindEndOfCurrentWordUntilSeparatorIgnoreNextToCursor(sep string) istrings.ByteNumber {
	if sep == "" {
		return d.FindEndOfCurrentWordWithSpace()
	}

	x := d.TextAfterCursor()

	start := istrings.IndexNotAny(x, sep)
	if start == -1 {
		return istrings.ByteNumber(len(x))
	}

	end := istrings.ByteNumber(strings.IndexAny(x[start:], sep))
	if end == -1 {
		return istrings.ByteNumber(len(x))
	}

	return start + end
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

// Array pointing to the start indices of all the lines.
func (d *Document) lineStartIndices() []istrings.RuneNumber {
	// TODO: Cache, because this is often reused.
	// (If it is used, it's often used many times.
	// And this has to be fast for editing big documents!)
	lc := d.LineCount()
	lengths := make([]istrings.RuneNumber, lc)
	for i, l := range d.Lines() {
		lengths[i] = istrings.RuneNumber(len([]rune(l)))
	}

	// Calculate cumulative sums.
	indices := make([]istrings.RuneNumber, lc+1)
	indices[0] = 0 // https://github.com/jonathanslenders/python-prompt-toolkit/blob/master/prompt_toolkit/document.py#L189
	var pos istrings.RuneNumber
	for i, l := range lengths {
		pos += l + 1
		indices[i+1] = istrings.RuneNumber(pos)
	}
	if lc > 1 {
		// Pop the last item. (This is not a new line.)
		indices = indices[:lc]
	}
	return indices
}

// For the index of a character at a certain line, calculate the index of
// the first character on that line.
func (d *Document) findLineStartIndex(index istrings.RuneNumber) (pos, lineStartIndex istrings.RuneNumber) {
	indices := d.lineStartIndices()
	pos = bisect.Right(indices, index) - 1
	lineStartIndex = indices[pos]
	return
}

// CursorPositionRow returns the current row. (0-based.)
func (d *Document) CursorPositionRow() (row istrings.RuneNumber) {
	row, _ = d.findLineStartIndex(d.cursorPosition)
	return
}

// TextEndPositionRow returns the row of the end of the current text. (0-based.)
func (d *Document) TextEndPositionRow() (row istrings.RuneNumber) {
	textLength := istrings.RuneCount(d.Text)
	if textLength == 0 {
		return 0
	}
	row, _ = d.findLineStartIndex(textLength - 1)
	return
}

// CursorPositionCol returns the current column. (0-based.)
func (d *Document) CursorPositionCol() (col istrings.RuneNumber) {
	_, index := d.findLineStartIndex(d.cursorPosition)
	col = d.cursorPosition - index
	return
}

// GetCursorLeftPosition returns the relative position for cursor left.
func (d *Document) GetCursorLeftPosition(count istrings.RuneNumber) istrings.RuneNumber {
	if count < 0 {
		return d.GetCursorRightPosition(-count)
	}
	runeSlice := []rune(d.Text)
	var counter istrings.RuneNumber
	targetPosition := d.cursorPosition - count
	if targetPosition < 0 {
		targetPosition = 0
	}
	for range runeSlice[targetPosition:d.cursorPosition] {
		counter--
	}

	return counter
}

// GetCursorRightPosition returns relative position for cursor right.
func (d *Document) GetCursorRightPosition(count istrings.RuneNumber) istrings.RuneNumber {
	if count < 0 {
		return d.GetCursorLeftPosition(-count)
	}
	runeSlice := []rune(d.Text)
	var counter istrings.RuneNumber
	targetPosition := d.cursorPosition + count
	if targetPosition > istrings.RuneNumber(len(runeSlice)) {
		targetPosition = istrings.RuneNumber(len(runeSlice))
	}
	for range runeSlice[d.cursorPosition:targetPosition] {
		counter++
	}

	return counter
}

// Get the current cursor position.
func (d *Document) GetCursorPosition(columns istrings.Width) Position {
	return positionAtEndOfString(d.TextBeforeCursor(), columns)
}

// Get the position of the end of the current text.
func (d *Document) GetEndOfTextPosition(columns istrings.Width) Position {
	return positionAtEndOfString(d.Text, columns)
}

// GetCursorUpPosition return the relative cursor position (character index) where we would be
// if the user pressed the arrow-up button.
func (d *Document) GetCursorUpPosition(count int, preferredColumn istrings.RuneNumber) istrings.RuneNumber {
	var col istrings.RuneNumber
	if preferredColumn == -1 { // -1 means nil
		col = d.CursorPositionCol()
	} else {
		col = preferredColumn
	}

	row := int(d.CursorPositionRow()) - count
	if row < 0 {
		row = 0
	}
	return d.TranslateRowColToIndex(row, col) - d.cursorPosition
}

// GetCursorDownPosition return the relative cursor position (character index) where we would be if the
// user pressed the arrow-down button.
func (d *Document) GetCursorDownPosition(count int, preferredColumn istrings.RuneNumber) istrings.RuneNumber {
	var col istrings.RuneNumber
	if preferredColumn == -1 { // -1 means nil
		col = d.CursorPositionCol()
	} else {
		col = preferredColumn
	}
	row := int(d.CursorPositionRow()) + count
	return d.TranslateRowColToIndex(row, col) - d.cursorPosition
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
func (d *Document) TranslateIndexToPosition(index istrings.RuneNumber) (int, int) {
	r, rowIndex := d.findLineStartIndex(index)
	c := index - rowIndex
	return int(r), int(c)
}

// TranslateRowColToIndex given a (row, col), return the corresponding index.
// (Row and col params are 0-based.)
func (d *Document) TranslateRowColToIndex(row int, column istrings.RuneNumber) (index istrings.RuneNumber) {
	indices := d.lineStartIndices()
	if row < 0 {
		row = 0
	} else if row > len(indices) {
		row = len(indices) - 1
	}
	index = indices[row]
	line := []rune(d.Lines()[row])

	// python) result += max(0, min(col, len(line)))
	if column > 0 || len(line) > 0 {
		if column > istrings.RuneNumber(len(line)) {
			index += istrings.RuneNumber(len(line))
		} else {
			index += istrings.RuneNumber(column)
		}
	}

	text := []rune(d.Text)
	// Keep in range. (len(self.text) is included, because the cursor can be
	// right after the end of the text as well.)
	// python) result = max(0, min(result, len(self.text)))
	if index > istrings.RuneNumber(len(text)) {
		index = istrings.RuneNumber(len(text))
	}
	if index < 0 {
		index = 0
	}
	return index
}

// OnLastLine returns true when we are at the last line.
func (d *Document) OnLastLine() bool {
	return d.CursorPositionRow() == istrings.RuneNumber(d.LineCount()-1)
}

// GetEndOfLinePosition returns relative position for the end of this line.
func (d *Document) GetEndOfLinePosition() istrings.RuneNumber {
	return istrings.RuneCount(d.CurrentLineAfterCursor())
}

// GetStartOfLinePosition returns relative position for the start of this line.
func (d *Document) GetStartOfLinePosition() istrings.RuneNumber {
	return istrings.RuneCount(d.CurrentLineBeforeCursor())
}

// GetStartOfLinePosition returns relative position for the start of this line.
func (d *Document) FindStartOfFirstWordOfLine() istrings.RuneNumber {
	line := d.CurrentLineBeforeCursor()
	var counter istrings.RuneNumber
	var nonSpaceCharSeen bool
	for _, char := range line {
		if !nonSpaceCharSeen && unicode.IsSpace(char) {
			continue
		}

		if !nonSpaceCharSeen {
			nonSpaceCharSeen = true
		}
		counter++
	}

	if counter == 0 {
		return istrings.RuneCount(line)
	}

	return counter
}

func (d *Document) leadingWhitespaceInCurrentLine() (margin string) {
	trimmed := strings.TrimSpace(d.CurrentLine())
	margin = d.CurrentLine()[:len(d.CurrentLine())-len(trimmed)]
	return
}
