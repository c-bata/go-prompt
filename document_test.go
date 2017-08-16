package prompt

import (
	"reflect"
	"testing"
	"unicode/utf8"
)

func TestDocument_GetCharRelativeToCursor(t *testing.T) {
	d := &Document{
		Text:           "line 1\nline 2\nline 3\nline 4\n",
		CursorPosition: len([]rune("line 1\n" + "lin")),
	}
	ac := d.GetCharRelativeToCursor(1)
	ex, _ := utf8.DecodeRuneInString("e")
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}
}

func TestDocument_TextBeforeCursor(t *testing.T) {
	d := &Document{
		Text:           "line 1\nline 2\nline 3\nline 4\n",
		CursorPosition: len("line 1\n" + "lin"),
	}
	ac := d.TextBeforeCursor()
	ex := "line 1\nlin"
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}
}

func TestDocument_TextAfterCursor(t *testing.T) {
	pattern := []struct {
		document *Document
		expected string
	}{
		{
			document: &Document{
				Text:           "line 1\nline 2\nline 3\nline 4\n",
				CursorPosition: len("line 1\n" + "lin"),
			},
			expected: "e 2\nline 3\nline 4\n",
		},
		{
			document: &Document{
				Text:           "",
				CursorPosition: 0,
			},
			expected: "",
		},
	}

	for _, p := range pattern {
		ac := p.document.TextAfterCursor()
		if ac != p.expected {
			t.Errorf("Should be %#v, got %#v", p.expected, ac)
		}
	}
}

func TestDocument_GetWordBeforeCursor(t *testing.T) {
	pattern := []struct {
		document *Document
		expected string
	}{
		{
			document: &Document{
				Text:           "apple bana",
				CursorPosition: len("apple bana"),
			},
			expected: "bana",
		},
		{
			document: &Document{
				Text:           "apple ",
				CursorPosition: len("apple "),
			},
			expected: "",
		},
	}

	for _, p := range pattern {
		ac := p.document.GetWordBeforeCursor()
		if ac != p.expected {
			t.Errorf("Should be %#v, got %#v", p.expected, ac)
		}
	}
}

func TestDocument_GetWordBeforeCursorWithSpace(t *testing.T) {
	pattern := []struct {
		document *Document
		expected string
	}{
		{
			document: &Document{
				Text:           "apple bana ",
				CursorPosition: len("apple bana "),
			},
			expected: "bana ",
		},
		{
			document: &Document{
				Text:           "apple ",
				CursorPosition: len("apple "),
			},
			expected: "apple ",
		},
	}

	for _, p := range pattern {
		ac := p.document.GetWordBeforeCursorWithSpace()
		if ac != p.expected {
			t.Errorf("Should be %#v, got %#v", p.expected, ac)
		}
	}
}

func TestDocument_FindStartOfPreviousWord(t *testing.T) {
	pattern := []struct {
		document *Document
		expected int
	}{
		{
			document: &Document{
				Text:           "apple bana",
				CursorPosition: len("apple bana"),
			},
			expected: len("apple "),
		},
		{
			document: &Document{
				Text:           "apple ",
				CursorPosition: len("apple "),
			},
			expected: len("apple "),
		},
	}

	for _, p := range pattern {
		ac := p.document.FindStartOfPreviousWord()
		if ac != p.expected {
			t.Errorf("Should be %#v, got %#v", p.expected, ac)
		}
	}
}

func TestDocument_FindStartOfPreviousWordWithSpace(t *testing.T) {
	pattern := []struct {
		document *Document
		expected int
	}{
		{
			document: &Document{
				Text:           "apple bana ",
				CursorPosition: len("apple bana "),
			},
			expected: len("apple "),
		},
		{
			document: &Document{
				Text:           "apple ",
				CursorPosition: len("apple "),
			},
			expected: len(""),
		},
	}

	for _, p := range pattern {
		ac := p.document.FindStartOfPreviousWordWithSpace()
		if ac != p.expected {
			t.Errorf("Should be %#v, got %#v", p.expected, ac)
		}
	}
}

func TestDocument_CurrentLineBeforeCursor(t *testing.T) {
	d := &Document{
		Text:           "line 1\nline 2\nline 3\nline 4\n",
		CursorPosition: len("line 1\n" + "lin"),
	}
	ac := d.CurrentLineBeforeCursor()
	ex := "lin"
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}
}

func TestDocument_CurrentLineAfterCursor(t *testing.T) {
	d := &Document{
		Text:           "line 1\nline 2\nline 3\nline 4\n",
		CursorPosition: len("line 1\n" + "lin"),
	}
	ac := d.CurrentLineAfterCursor()
	ex := "e 2"
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}
}

func TestDocument_CurrentLine(t *testing.T) {
	d := &Document{
		Text:           "line 1\nline 2\nline 3\nline 4\n",
		CursorPosition: len("line 1\n" + "lin"),
	}
	ac := d.CurrentLine()
	ex := "line 2"
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}
}

// Table Driven Tests for CursorPositionRow and CursorPositionCol
type cursorPositionTest struct {
	document    *Document
	expectedRow int
	expectedCol int
}

var cursorPositionTests = []cursorPositionTest{
	{
		document:    &Document{Text: "line 1\nline 2\nline 3\n", CursorPosition: len("line 1\n" + "lin")},
		expectedRow: 1,
		expectedCol: 3,
	},
	{
		document:    &Document{Text: "", CursorPosition: 0},
		expectedRow: 0,
		expectedCol: 0,
	},
}

func TestDocument_CursorPositionRowAndCol(t *testing.T) {
	for _, test := range cursorPositionTests {
		ac := test.document.CursorPositionRow()
		if ac != test.expectedRow {
			t.Errorf("Should be %#v, got %#v", test.expectedRow, ac)
		}
		ac = test.document.CursorPositionCol()
		if ac != test.expectedCol {
			t.Errorf("Should be %#v, got %#v", test.expectedCol, ac)
		}
	}
}

func TestDocument_GetCursorLeftPosition(t *testing.T) {
	d := &Document{
		Text:           "line 1\nline 2\nline 3\nline 4\n",
		CursorPosition: len("line 1\n" + "line 2\n" + "lin"),
	}
	ac := d.GetCursorLeftPosition(2)
	ex := -2
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}
	ac = d.GetCursorLeftPosition(10)
	ex = -3
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}
}

func TestDocument_GetCursorUpPosition(t *testing.T) {
	d := &Document{
		Text:           "line 1\nline 2\nline 3\nline 4\n",
		CursorPosition: len("line 1\n" + "line 2\n" + "lin"),
	}
	ac := d.GetCursorUpPosition(2, -1)
	ex := len("lin") - len("line 1\n"+"line 2\n"+"lin")
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}

	ac = d.GetCursorUpPosition(100, -1)
	ex = len("lin") - len("line 1\n"+"line 2\n"+"lin")
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}
}

func TestDocument_GetCursorDownPosition(t *testing.T) {
	d := &Document{
		Text:           "line 1\nline 2\nline 3\nline 4\n",
		CursorPosition: len("lin"),
	}
	ac := d.GetCursorDownPosition(2, -1)
	ex := len("line 1\n"+"line 2\n"+"lin") - len("lin")
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}

	ac = d.GetCursorDownPosition(100, -1)
	ex = len("line 1\n"+"line 2\n"+"line 3\n"+"line 4\n") - len("lin")
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}
}

func TestDocument_GetCursorRightPosition(t *testing.T) {
	d := &Document{
		Text:           "line 1\nline 2\nline 3\nline 4\n",
		CursorPosition: len("line 1\n" + "line 2\n" + "lin"),
	}
	ac := d.GetCursorRightPosition(2)
	ex := 2
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}
	ac = d.GetCursorRightPosition(10)
	ex = 3
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}
}

func TestDocument_Lines(t *testing.T) {
	d := &Document{
		Text:           "line 1\nline 2\nline 3\nline 4\n",
		CursorPosition: len("line 1\n" + "lin"),
	}
	ac := d.Lines()
	ex := []string{"line 1", "line 2", "line 3", "line 4", ""}
	if !reflect.DeepEqual(ac, ex) {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}
}

func TestDocument_LineCount(t *testing.T) {
	d := &Document{
		Text:           "line 1\nline 2\nline 3\nline 4\n",
		CursorPosition: len("line 1\n" + "lin"),
	}
	ac := d.LineCount()
	ex := 5
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}
}

func TestDocument_TranslateIndexToPosition(t *testing.T) {
	d := &Document{
		Text:           "line 1\nline 2\nline 3\nline 4\n",
		CursorPosition: len("line 1\n" + "lin"),
	}
	row, col := d.TranslateIndexToPosition(len("line 1\nline 2\nlin"))
	if row != 2 {
		t.Errorf("Should be %#v, got %#v", 2, row)
	}
	if col != 3 {
		t.Errorf("Should be %#v, got %#v", 3, col)
	}
	row, col = d.TranslateIndexToPosition(0)
	if row != 0 {
		t.Errorf("Should be %#v, got %#v", 0, row)
	}
	if col != 0 {
		t.Errorf("Should be %#v, got %#v", 0, col)
	}
}

func TestDocument_TranslateRowColToIndex(t *testing.T) {
	d := &Document{
		Text:           "line 1\nline 2\nline 3\nline 4\n",
		CursorPosition: len("line 1\n" + "lin"),
	}
	ac := d.TranslateRowColToIndex(2, 3)
	ex := len("line 1\nline 2\nlin")
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}
	ac = d.TranslateRowColToIndex(0, 0)
	ex = 0
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}
}

func TestDocument_OnLastLine(t *testing.T) {
	d := &Document{
		Text:           "line 1\nline 2\nline 3",
		CursorPosition: len("line 1\nline"),
	}
	ac := d.OnLastLine()
	if ac {
		t.Errorf("Should be %#v, got %#v", false, ac)
	}
	d.CursorPosition = len("line 1\nline 2\nline")
	ac = d.OnLastLine()
	if !ac {
		t.Errorf("Should be %#v, got %#v", true, ac)
	}
}

func TestDocument_GetEndOfLinePosition(t *testing.T) {
	d := &Document{
		Text:           "line 1\nline 2\nline 3",
		CursorPosition: len("line 1\nli"),
	}
	ac := d.GetEndOfLinePosition()
	ex := len("ne 2")
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}
}
