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
