package prompt

import (
	istrings "github.com/elk-language/go-prompt/internal/strings"
)

// GoLineEnd Go to the End of the line
func GoLineEnd(buf *Buffer) {
	x := []rune(buf.Document().TextAfterCursor())
	buf.CursorRight(istrings.RuneCount(len(x)))
}

// GoLineBeginning Go to the beginning of the line
func GoLineBeginning(buf *Buffer) {
	x := []rune(buf.Document().TextBeforeCursor())
	buf.CursorLeft(istrings.RuneCount(len(x)))
}

// DeleteChar Delete character under the cursor
func DeleteChar(buf *Buffer) {
	buf.Delete(1)
}

// DeleteWord Delete word before the cursor
func DeleteWord(buf *Buffer) {
	buf.DeleteBeforeCursor(istrings.RuneCount(len([]rune(buf.Document().TextBeforeCursor()))) - istrings.RuneCount(buf.Document().FindStartOfPreviousWordWithSpace())) // WARN
}

// DeleteBeforeChar Go to Backspace
func DeleteBeforeChar(buf *Buffer) {
	buf.DeleteBeforeCursor(1)
}

// GoRightChar Forward one character
func GoRightChar(buf *Buffer) {
	buf.CursorRight(1)
}

// GoLeftChar Backward one character
func GoLeftChar(buf *Buffer) {
	buf.CursorLeft(1)
}

// GoRightWord Forward one word
func GoRightWord(buf *Buffer) {
	buf.CursorRight(istrings.RuneCount(buf.Document().FindEndOfCurrentWordWithSpace())) // WARN
}

// GoLeftWord Backward one word
func GoLeftWord(buf *Buffer) {
	buf.CursorLeft(istrings.RuneCount(len([]rune(buf.Document().TextBeforeCursor()))) - istrings.RuneCount(buf.Document().FindStartOfPreviousWordWithSpace())) // WARN
}
