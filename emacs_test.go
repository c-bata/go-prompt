package prompt

import (
	"testing"

	istrings "github.com/elk-language/go-prompt/internal/strings"
)

func TestEmacsKeyBindings(t *testing.T) {
	buf := NewBuffer()
	buf.InsertText("abcde", false, true)
	if buf.cursorPosition != istrings.RuneIndex(len("abcde")) {
		t.Errorf("Want %d, but got %d", len("abcde"), buf.cursorPosition)
	}

	// Go to the beginning of the line
	applyEmacsKeyBind(buf, ControlA)
	if buf.cursorPosition != 0 {
		t.Errorf("Want %d, but got %d", 0, buf.cursorPosition)
	}

	// Go to the end of the line
	applyEmacsKeyBind(buf, ControlE)
	if buf.cursorPosition != istrings.RuneIndex(len("abcde")) {
		t.Errorf("Want %d, but got %d", len("abcde"), buf.cursorPosition)
	}
}

func applyEmacsKeyBind(buf *Buffer, key Key) {
	for i := range emacsKeyBindings {
		kb := emacsKeyBindings[i]
		if kb.Key == key {
			kb.Fn(buf)
		}
	}
}
