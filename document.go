package prompt

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/elk-language/go-prompt/bisect"
	istrings "github.com/elk-language/go-prompt/strings"
	"github.com/rivo/uniseg"
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

// Returns the index of the rune that's under the cursor.
func (d *Document) CurrentRuneIndex() istrings.RuneNumber {
	return d.cursorPosition
}

// Returns the amount of spaces that the last line of input
// of the given text is indented with.
func (d *Document) IndentSpaces(input string) int {
	lastNewline := strings.LastIndexByte(input, '\n')
	var spaces int
	for i := lastNewline + 1; i < len(input); i++ {
		b := input[i]
		if b != ' ' {
			break
		}

		spaces++
	}

	return spaces
}

// Returns the indentation level of the last line of the given text.
func (d *Document) IndentLevel(input string, indentSize int) int {
	if indentSize == 0 {
		return 0
	}
	return d.IndentSpaces(input) / indentSize
}

// Returns the amount of spaces that the last line of input
// is indented with.
func (d *Document) LastLineIndentSpaces() int {
	return d.IndentSpaces(d.Text)
}

// Returns the indentation level of the last line of input.
func (d *Document) LastLineIndentLevel(indentSize int) int {
	return d.IndentLevel(d.Text, indentSize)
}

// Returns the amount of spaces that the current line the cursor is on
// is indented with.
func (d *Document) CurrentLineIndentSpaces() int {
	return d.IndentSpaces(d.TextBeforeCursor())
}

// Returns the indentation level of the current line the cursor is on.
func (d *Document) CurrentLineIndentLevel(indentSize int) int {
	return d.IndentLevel(d.TextBeforeCursor(), indentSize)
}

// Returns the amount of spaces that the previous line (relative to the cursor)
// is indented with.
func (d *Document) PreviousLineIndentSpaces() int {
	line, ok := d.PreviousLine()
	if !ok {
		return 0
	}
	return d.IndentSpaces(line)
}

// Returns the indentation level of the previous line (relative to the cursor).
func (d *Document) PreviousLineIndentLevel(indentSize int) int {
	line, ok := d.PreviousLine()
	if !ok {
		return 0
	}
	return d.IndentLevel(line, indentSize)
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

// GetWordBeforeCursor returns the identifier or punctuation run before the cursor.
// Whitespace immediately before the cursor produces an empty string.
func (d *Document) GetWordBeforeCursor() string {
	x := d.TextBeforeCursor()
	return x[d.FindStartOfPreviousWord():]
}

// GetWordAfterCursor returns the identifier or punctuation run after the cursor.
// Whitespace immediately after the cursor produces an empty string.
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

// FindStartOfPreviousWord returns a byte index in the text before the cursor
// pointing at the start of the part that ends at the cursor.
// The index is the end of that text when whitespace is immediately before the cursor.
func (d *Document) FindStartOfPreviousWord() istrings.ByteNumber {
	return startOfPartBefore(d.TextBeforeCursor(), false)
}

// Returns how many runes the cursor must move left to reach the end of the
// previous part. Parts are identifier runs and punctuation runs. Whitespace
// is skipped, so the stops in foo.bar(elo, dupa) are the ends of
// foo, ., bar, (, elo, ,, dupa, and ).
func (d *Document) FindRuneNumberUntilStartOfPreviousWord() istrings.RuneNumber {
	s := d.TextBeforeCursor()
	if s == "" {
		return 0
	}

	r, _ := utf8.DecodeLastRuneInString(s)
	var prevEnd string
	if classifyWordPart(r) == wordPartSpace {
		prevEnd = strings.TrimRightFunc(s, isWordPartSpace)
	} else {
		prevEnd = strings.TrimRightFunc(s[:startOfPartBefore(s, false)], isWordPartSpace)
	}
	return istrings.RuneCountInString(s) - istrings.RuneCountInString(prevEnd)
}

// FindStartOfPreviousWordWithSpace is almost the same as FindStartOfPreviousWord.
// Whitespace immediately before the cursor is kept, and the part before that whitespace is used.
func (d *Document) FindStartOfPreviousWordWithSpace() istrings.ByteNumber {
	return startOfPartBefore(d.TextBeforeCursor(), true)
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

// FindEndOfCurrentWord returns a byte index in the text after the cursor
// pointing at the end of the part that starts at the cursor.
// The index is 0 when whitespace is immediately after the cursor.
func (d *Document) FindEndOfCurrentWord() istrings.ByteNumber {
	return endOfPartAfter(d.TextAfterCursor(), false)
}

// FindEndOfCurrentWordWithSpace is almost the same as FindEndOfCurrentWord.
// Whitespace immediately after the cursor is kept, and the part after that whitespace is used.
func (d *Document) FindEndOfCurrentWordWithSpace() istrings.ByteNumber {
	return endOfPartAfter(d.TextAfterCursor(), true)
}

// isWordCharacter reports whether char is part of an identifier.
// A run of these characters is one part, such as foo or dupa.
func isWordCharacter(char rune) bool {
	return char == '_' || unicode.IsLetter(char) || unicode.IsDigit(char)
}

type wordPartClass uint8

const (
	wordPartWord wordPartClass = iota
	wordPartPunct
	wordPartSpace
)

// classifyWordPart groups a character into an identifier, punctuation, or whitespace.
// Consecutive characters of the same non-space class form one word.
func classifyWordPart(char rune) wordPartClass {
	switch {
	case unicode.IsSpace(char):
		return wordPartSpace
	case isWordCharacter(char):
		return wordPartWord
	default:
		return wordPartPunct
	}
}

func isWordPartSpace(char rune) bool {
	return classifyWordPart(char) == wordPartSpace
}

// startOfPartBefore returns the byte index where the part at the end of s begins.
// When skipSpace is set and s ends in whitespace, that whitespace stays in the
// returned span and the part before it is used. When s is only whitespace, the index is 0.
// When skipSpace is unset and s ends in whitespace, the index is len(s).
func startOfPartBefore(s string, skipSpace bool) istrings.ByteNumber {
	i := len(s)
	var class wordPartClass
	started := false
	for i > 0 {
		char, size := utf8.DecodeLastRuneInString(s[:i])
		charClass := classifyWordPart(char)
		if !started {
			if charClass == wordPartSpace {
				if !skipSpace {
					return istrings.ByteNumber(len(s))
				}
				i -= size
				continue
			}
			class = charClass
			started = true
			i -= size
			continue
		}
		if charClass != class {
			break
		}
		i -= size
	}
	return istrings.ByteNumber(i)
}

// endOfPartAfter returns the byte index where the part at the start of s ends.
// When skipSpace is set and s starts with whitespace, that whitespace stays in the
// returned span and the part after it is used. When s is only whitespace, the index is len(s).
// When skipSpace is unset and s starts with whitespace, the index is 0.
func endOfPartAfter(s string, skipSpace bool) istrings.ByteNumber {
	i := 0
	var class wordPartClass
	started := false
	for i < len(s) {
		char, size := utf8.DecodeRuneInString(s[i:])
		charClass := classifyWordPart(char)
		if !started {
			if charClass == wordPartSpace {
				if !skipSpace {
					return 0
				}
				i += size
				continue
			}
			class = charClass
			started = true
			i += size
			continue
		}
		if charClass != class {
			break
		}
		i += size
	}
	return istrings.ByteNumber(i)
}

// Returns how many runes the cursor must move right to reach the end of the
// current part. Parts are identifier runs and punctuation runs. Whitespace
// is skipped, so the stops in foo.bar(elo, dupa) are the ends of
// foo, ., bar, (, elo, ,, dupa, and ).
func (d *Document) FindRuneNumberUntilEndOfCurrentWord() istrings.RuneNumber {
	s := d.TextAfterCursor()
	return istrings.RuneCountInString(s[:endOfPartAfter(s, true)])
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

// Return the text of the previous line (relative to the cursor).
// If the cursor is on the first line then false is returned in the second value
// to signify that there is no previous line.
func (d *Document) PreviousLine() (s string, ok bool) {
	indices := d.lineStartIndices()
	pos := bisect.Right(indices, d.cursorPosition) - 1
	if pos == 0 {
		return "", false
	}

	prevLineStartIndex := indices[pos-1]
	lineStartIndex := indices[pos]
	return d.Text[prevLineStartIndex : lineStartIndex-1], true
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
	textLength := istrings.RuneCountInString(d.Text)
	if textLength == 0 {
		return 0
	}
	row, _ = d.findLineStartIndex(textLength - 1)
	return
}

// CursorPositionCol returns the current column. (0-based.)
func (d *Document) CursorPositionCol() (col istrings.Width) {
	_, lineStartIndex := d.findLineStartIndex(d.cursorPosition)

	text := utf8string.NewString(d.Text).Slice(int(lineStartIndex), int(d.cursorPosition))
	return istrings.GetWidth(text)
}

// Returns the amount of runes that the cursors should be moved by.
// The `count` argument tells this function by how many graphemes (visible characters)
// the cursor should be moved (to the left).
func (d *Document) GetCursorLeftPosition(count istrings.GraphemeNumber) istrings.RuneNumber {
	if count < 0 {
		return d.GetCursorRightPosition(-count)
	}
	if d.cursorPosition == 0 {
		return 0
	}
	text := d.TextBeforeCursor()
	g := uniseg.NewGraphemes(text)
	graphemeLength := istrings.GraphemeCountInString(text)
	var currentGraphemeIndex istrings.GraphemeNumber
	var currentPosition istrings.RuneNumber

	for g.Next() {
		if currentGraphemeIndex >= graphemeLength-count {
			break
		}
		currentPosition += istrings.RuneNumber(len(g.Runes()))
		currentGraphemeIndex++
	}

	result := d.cursorPosition - currentPosition
	return -result
}

// Returns the amount of runes that the cursors should be moved by.
// The `count` argument tells this function by how many runes
// the cursor should be moved (to the left).
func (d *Document) GetCursorLeftPositionRunes(count istrings.RuneNumber) istrings.RuneNumber {
	if count < 0 {
		return d.GetCursorRightPositionRunes(-count)
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

// Returns the amount of runes that the cursors should be moved by.
// The `count` argument tells this function by how many graphemes (visible characters)
// the cursor should be moved (to the right).
func (d *Document) GetCursorRightPosition(count istrings.GraphemeNumber) istrings.RuneNumber {
	if count < 0 {
		return d.GetCursorLeftPosition(-count)
	}
	text := d.TextAfterCursor()
	if len(text) == 0 {
		return 0
	}

	return istrings.RuneIndexNthGrapheme(text, count)
}

// Returns the amount of runes that the cursors should be moved by.
// The `count` argument tells this function by how many runes
// the cursor should be moved (to the right).
func (d *Document) GetCursorRightPositionRunes(count istrings.RuneNumber) istrings.RuneNumber {
	if count < 0 {
		return d.GetCursorLeftPositionRunes(-count)
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
func (d *Document) GetCursorUpPosition(count int, preferredColumn istrings.Width) istrings.RuneNumber {
	var col istrings.Width
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
func (d *Document) GetCursorDownPosition(count int, preferredColumn istrings.Width) istrings.RuneNumber {
	var col istrings.Width
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
func (d *Document) TranslateRowColToIndex(row int, column istrings.Width) (index istrings.RuneNumber) {
	indices := d.lineStartIndices()
	if row < 0 {
		row = 0
	} else if row > len(indices) {
		row = len(indices) - 1
	}
	index = indices[row]
	line := d.Lines()[row]

	index += istrings.RuneIndexNthColumn(line, column)

	runeLength := istrings.RuneCountInString(d.Text)
	// Keep in range. (len(self.text) is included, because the cursor can be
	// right after the end of the text as well.)
	// python) result = max(0, min(result, len(self.text)))
	if index > runeLength {
		index = runeLength
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
	return istrings.RuneCountInString(d.CurrentLineAfterCursor())
}

// GetStartOfLinePosition returns relative position for the start of this line.
func (d *Document) GetStartOfLinePosition() istrings.RuneNumber {
	return istrings.RuneCountInString(d.CurrentLineBeforeCursor())
}

// FindStartOfFirstWordOfLine returns how many runes the cursor must move left
// to reach the first part on this line. From inside the leading whitespace it
// returns the distance to the start of the line.
func (d *Document) FindStartOfFirstWordOfLine() istrings.RuneNumber {
	line := d.CurrentLineBeforeCursor()
	trimmed := strings.TrimLeftFunc(line, isWordPartSpace)
	if trimmed == "" {
		return istrings.RuneCountInString(line)
	}
	return istrings.RuneCountInString(trimmed)
}

func (d *Document) leadingWhitespaceInCurrentLine() (margin string) {
	trimmed := strings.TrimSpace(d.CurrentLine())
	margin = d.CurrentLine()[:len(d.CurrentLine())-len(trimmed)]
	return
}
