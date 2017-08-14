package prompt

import "testing"

func TestEmacsKeyBindings(t *testing.T) {
	buf := NewBuffer()
	buf.InsertText("abcde", false, true)
	if buf.CursorPosition != len("abcde") {
		t.Errorf("Want %d, but got %d", len("abcde"), buf.CursorPosition)
	}

	// Go to the beginning of the line
	applyEmacsKeyBind(buf, ControlA)
	if buf.CursorPosition != 0 {
		t.Errorf("Want %d, but got %d", 0, buf.CursorPosition)
	}

	// Go to the end of the line
	applyEmacsKeyBind(buf, ControlE)
	if buf.CursorPosition != len("abcde") {
		t.Errorf("Want %d, but got %d", len("abcde"), buf.CursorPosition)
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
