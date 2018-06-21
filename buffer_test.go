package prompt

import (
	"reflect"
	"testing"
)

func TestNewBuffer(t *testing.T) {
	b := NewBuffer()
	if b.workingIndex != 0 {
		t.Errorf("workingIndex should be %#v, got %#v", 0, b.workingIndex)
	}
	if !reflect.DeepEqual(b.workingLines, []string{""}) {
		t.Errorf("workingLines should be %#v, got %#v", []string{""}, b.workingLines)
	}
}

func TestBuffer_InsertText(t *testing.T) {
	b := NewBuffer()
	b.InsertText("some_text", false, true)

	if b.Text() != "some_text" {
		t.Errorf("Text should be %#v, got %#v", "some_text", b.Text())
	}

	if b.cursorPosition != len("some_text") {
		t.Errorf("cursorPosition should be %#v, got %#v", len("some_text"), b.cursorPosition)
	}
}

func TestBuffer_CursorMovement(t *testing.T) {
	b := NewBuffer()
	b.InsertText("some_text", false, true)

	b.CursorLeft(1)
	b.CursorLeft(2)
	b.CursorRight(1)
	b.InsertText("A", false, true)
	if b.Text() != "some_teAxt" {
		t.Errorf("Text should be %#v, got %#v", "some_teAxt", b.Text())
	}
	if b.cursorPosition != len("some_teA") {
		t.Errorf("Text should be %#v, got %#v", len("some_teA"), b.cursorPosition)
	}

	// Moving over left character counts.
	b.CursorLeft(100)
	b.InsertText("A", false, true)
	if b.Text() != "Asome_teAxt" {
		t.Errorf("Text should be %#v, got %#v", "some_teAxt", b.Text())
	}
	if b.cursorPosition != len("A") {
		t.Errorf("Text should be %#v, got %#v", len("some_teA"), b.cursorPosition)
	}

	// TODO: Going right already at right end.
}

func TestBuffer_CursorMovement_WithMultiByte(t *testing.T) {
	b := NewBuffer()
	b.InsertText("あいうえお", false, true)
	b.CursorLeft(1)
	if l := b.Document().TextAfterCursor(); l != "お" {
		t.Errorf("Should be 'お', but got %s", l)
	}
}

func TestBuffer_CursorUp(t *testing.T) {
	b := NewBuffer()
	b.InsertText("long line1\nline2", false, true)
	b.CursorUp(1)
	if b.Document().cursorPosition != 5 {
		t.Errorf("Should be %#v, got %#v", 5, b.Document().cursorPosition)
	}

	// Going up when already at the top.
	b.CursorUp(1)
	if b.Document().cursorPosition != 5 {
		t.Errorf("Should be %#v, got %#v", 5, b.Document().cursorPosition)
	}

	// Going up to a line that's shorter.
	b.setDocument(&Document{})
	b.InsertText("line1\nlong line2", false, true)
	b.CursorUp(1)
	if b.Document().cursorPosition != 5 {
		t.Errorf("Should be %#v, got %#v", 5, b.Document().cursorPosition)
	}
}

func TestBuffer_CursorDown(t *testing.T) {
	b := NewBuffer()
	b.InsertText("line1\nline2", false, true)
	b.cursorPosition = 3

	// Normally going down
	b.CursorDown(1)
	if b.Document().cursorPosition != len("line1\nlin") {
		t.Errorf("Should be %#v, got %#v", len("line1\nlin"), b.Document().cursorPosition)
	}

	// Going down to a line that's storter.
	b = NewBuffer()
	b.InsertText("long line1\na\nb", false, true)
	b.cursorPosition = 3
	b.CursorDown(1)
	if b.Document().cursorPosition != len("long line1\na") {
		t.Errorf("Should be %#v, got %#v", len("long line1\na"), b.Document().cursorPosition)
	}
}

func TestBuffer_DeleteBeforeCursor(t *testing.T) {
	b := NewBuffer()
	b.InsertText("some_text", false, true)
	b.CursorLeft(2)
	deleted := b.DeleteBeforeCursor(1)

	if b.Text() != "some_txt" {
		t.Errorf("Should be %#v, got %#v", "some_txt", b.Text())
	}
	if deleted != "e" {
		t.Errorf("Should be %#v, got %#v", deleted, "e")
	}
	if b.cursorPosition != len("some_t") {
		t.Errorf("Should be %#v, got %#v", len("some_t"), b.cursorPosition)
	}

	// Delete over the characters length before cursor.
	deleted = b.DeleteBeforeCursor(100)
	if deleted != "some_t" {
		t.Errorf("Should be %#v, got %#v", "some_t", deleted)
	}
	if b.Text() != "xt" {
		t.Errorf("Should be %#v, got %#v", "xt", b.Text())
	}

	// If cursor position is a beginning of line, it has no effect.
	deleted = b.DeleteBeforeCursor(1)
	if deleted != "" {
		t.Errorf("Should be empty, got %#v", deleted)
	}
}

func TestBuffer_NewLine(t *testing.T) {
	b := NewBuffer()
	b.InsertText("  hello", false, true)
	b.NewLine(false)
	ac := b.Text()
	ex := "  hello\n"
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}

	b = NewBuffer()
	b.InsertText("  hello", false, true)
	b.NewLine(true)
	ac = b.Text()
	ex = "  hello\n  "
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}
}

func TestBuffer_JoinNextLine(t *testing.T) {
	b := NewBuffer()
	b.InsertText("line1\nline2\nline3", false, true)
	b.CursorUp(1)
	b.JoinNextLine(" ")

	ac := b.Text()
	ex := "line1\nline2 line3"
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}

	// Test when there is no '\n' in the text
	b = NewBuffer()
	b.InsertText("line1", false, true)
	b.cursorPosition = 0
	b.JoinNextLine(" ")
	ac = b.Text()
	ex = "line1"
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}
}

func TestBuffer_SwapCharactersBeforeCursor(t *testing.T) {
	b := NewBuffer()
	b.InsertText("hello world", false, true)
	b.CursorLeft(2)
	b.SwapCharactersBeforeCursor()
	ac := b.Text()
	ex := "hello wrold"
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}
}
