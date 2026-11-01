package prompt

import (
	"fmt"
	"reflect"
	"testing"
	"unicode/utf8"

	istrings "github.com/elk-language/go-prompt/strings"
)

func ExampleDocument_CurrentLine() {
	d := &Document{
		Text: `Hello! my name is c-bata.
This is a example of Document component.
This component has texts displayed in terminal and cursor position.
`,
		cursorPosition: istrings.RuneCountInString(`Hello! my name is c-bata.
This is a exam`),
	}
	fmt.Println(d.CurrentLine())
	// Output:
	// This is a example of Document component.
}

func ExampleDocument_DisplayCursorPosition() {
	d := &Document{
		Text:           `Hello! my name is c-bata.`,
		cursorPosition: istrings.RuneCountInString(`Hello`),
	}
	fmt.Println("DisplayCursorPosition", d.DisplayCursorPosition(50))
	// Output:
	// DisplayCursorPosition {5 0}
}

func ExampleDocument_CursorPositionRow() {
	d := &Document{
		Text: `Hello! my name is c-bata.
This is a example of Document component.
This component has texts displayed in terminal and cursor position.
`,
		cursorPosition: istrings.RuneCountInString(`Hello! my name is c-bata.
This is a exam`),
	}
	fmt.Println("CursorPositionRow", d.CursorPositionRow())
	// Output:
	// CursorPositionRow 1
}

func ExampleDocument_CursorPositionCol() {
	d := &Document{
		Text: `Hello! my name is c-bata.
This is a example of Document component.
This component has texts displayed in terminal and cursor position.
`,
		cursorPosition: istrings.RuneCountInString(`Hello! my name is c-bata.
This is a exam`),
	}
	fmt.Println("CursorPositionCol", d.CursorPositionCol())
	// Output:
	// CursorPositionCol 14
}

func ExampleDocument_TextBeforeCursor() {
	d := &Document{
		Text: `Hello! my name is c-bata.
This is a example of Document component.
This component has texts displayed in terminal and cursor position.
`,
		cursorPosition: istrings.RuneCountInString(`Hello! my name is c-bata.
This is a exam`),
	}
	fmt.Println(d.TextBeforeCursor())
	// Output:
	// Hello! my name is c-bata.
	// This is a exam
}

func ExampleDocument_TextAfterCursor() {
	d := &Document{
		Text: `Hello! my name is c-bata.
This is a example of Document component.
This component has texts displayed in terminal and cursor position.
`,
		cursorPosition: istrings.RuneCountInString(`Hello! my name is c-bata.
This is a exam`),
	}
	fmt.Println(d.TextAfterCursor())
	// Output:
	// ple of Document component.
	// This component has texts displayed in terminal and cursor position.
}

func ExampleDocument_DisplayCursorPosition_withJapanese() {
	d := &Document{
		Text:           `こんにちは、芝田 将です。`,
		cursorPosition: 3,
	}
	fmt.Println("DisplayCursorPosition", d.DisplayCursorPosition(30))
	// Output:
	// DisplayCursorPosition {6 0}
}

func ExampleDocument_CurrentLineBeforeCursor() {
	d := &Document{
		Text: `Hello! my name is c-bata.
This is a example of Document component.
This component has texts displayed in terminal and cursor position.
`,
		cursorPosition: istrings.RuneCountInString(`Hello! my name is c-bata.
This is a exam`),
	}
	fmt.Println(d.CurrentLineBeforeCursor())
	// Output:
	// This is a exam
}

func ExampleDocument_CurrentLineAfterCursor() {
	d := &Document{
		Text: `Hello! my name is c-bata.
This is a example of Document component.
This component has texts displayed in terminal and cursor position.
`,
		cursorPosition: istrings.RuneCountInString(`Hello! my name is c-bata.
This is a exam`),
	}
	fmt.Println(d.CurrentLineAfterCursor())
	// Output:
	// ple of Document component.
}

func ExampleDocument_GetWordBeforeCursor() {
	d := &Document{
		Text: `Hello! my name is c-bata.
This is a example of Document component.
`,
		cursorPosition: istrings.RuneCountInString(`Hello! my name is c-bata.
This is a exam`),
	}
	fmt.Println(d.GetWordBeforeCursor())
	// Output:
	// exam
}

func ExampleDocument_GetWordAfterCursor() {
	d := &Document{
		Text: `Hello! my name is c-bata.
This is a example of Document component.
`,
		cursorPosition: istrings.RuneCountInString(`Hello! my name is c-bata.
This is a exam`),
	}
	fmt.Println(d.GetWordAfterCursor())
	// Output:
	// ple
}

func ExampleDocument_GetWordBeforeCursorWithSpace() {
	d := &Document{
		Text: `Hello! my name is c-bata.
This is a example of Document component.
`,
		cursorPosition: istrings.RuneCountInString(`Hello! my name is c-bata.
This is a example `),
	}
	fmt.Println(d.GetWordBeforeCursorWithSpace())
	// Output:
	// example
}

func ExampleDocument_GetWordAfterCursorWithSpace() {
	d := &Document{
		Text: `Hello! my name is c-bata.
This is a example of Document component.
`,
		cursorPosition: istrings.RuneCountInString(`Hello! my name is c-bata.
This is a`),
	}
	fmt.Println(d.GetWordAfterCursorWithSpace())
	// Output:
	//  example
}

func ExampleDocument_GetWordBeforeCursorUntilSeparator() {
	d := &Document{
		Text:           `hello,i am c-bata`,
		cursorPosition: istrings.RuneCountInString(`hello,i am c`),
	}
	fmt.Println(d.GetWordBeforeCursorUntilSeparator(","))
	// Output:
	// i am c
}

func ExampleDocument_GetWordAfterCursorUntilSeparator() {
	d := &Document{
		Text:           `hello,i am c-bata,thank you for using go-prompt`,
		cursorPosition: istrings.RuneCountInString(`hello,i a`),
	}
	fmt.Println(d.GetWordAfterCursorUntilSeparator(","))
	// Output:
	// m c-bata
}

func ExampleDocument_GetWordBeforeCursorUntilSeparatorIgnoreNextToCursor() {
	d := &Document{
		Text:           `hello,i am c-bata,thank you for using go-prompt`,
		cursorPosition: istrings.RuneCountInString(`hello,i am c-bata,`),
	}
	fmt.Println(d.GetWordBeforeCursorUntilSeparatorIgnoreNextToCursor(","))
	// Output:
	// i am c-bata,
}

func ExampleDocument_GetWordAfterCursorUntilSeparatorIgnoreNextToCursor() {
	d := &Document{
		Text:           `hello,i am c-bata,thank you for using go-prompt`,
		cursorPosition: istrings.RuneCountInString(`hello`),
	}
	fmt.Println(d.GetWordAfterCursorUntilSeparatorIgnoreNextToCursor(","))
	// Output:
	// ,i am c-bata
}

func TestDocument_DisplayCursorPosition(t *testing.T) {
	patterns := []struct {
		document *Document
		expected Position
	}{
		{
			document: &Document{
				Text:           "hello",
				cursorPosition: 2,
			},
			expected: Position{X: 2},
		},
		{
			document: &Document{
				Text:           "こんにちは",
				cursorPosition: 2,
			},
			expected: Position{X: 4},
		},
		{
			// If you're facing test failure on this test case and your terminal is iTerm2,
			// please check 'Profile -> Text' configuration. 'Use Unicode version 9 widths'
			// must be checked.
			// https://github.com/elk-language/go-prompt/pull/99
			document: &Document{
				Text:           "Добрый день",
				cursorPosition: 3,
			},
			expected: Position{X: 3},
		},
	}

	for _, p := range patterns {
		ac := p.document.DisplayCursorPosition(30)
		if ac != p.expected {
			t.Errorf("Should be %#v, got %#v", p.expected, ac)
		}
	}
}

func TestDocument_GetCharRelativeToCursor(t *testing.T) {
	patterns := []struct {
		document *Document
		expected string
	}{
		{
			document: &Document{
				Text:           "line 1\nline 2\nline 3\nline 4\n",
				cursorPosition: istrings.RuneCountInString("line 1\n" + "lin"),
			},
			expected: "e",
		},
		{
			document: &Document{
				Text:           "あいうえお\nかきくけこ\nさしすせそ\nたちつてと\n",
				cursorPosition: 8,
			},
			expected: "く",
		},
		{
			document: &Document{
				Text:           "Добрый\nдень\nДобрый день",
				cursorPosition: 9,
			},
			expected: "н",
		},
	}

	for i, p := range patterns {
		ac := p.document.GetCharRelativeToCursor(1)
		ex, _ := utf8.DecodeRuneInString(p.expected)
		if ac != ex {
			t.Errorf("[%d] Should be %s, got %s", i, string(ex), string(ac))
		}
	}
}

func TestDocument_TextBeforeCursor(t *testing.T) {
	patterns := []struct {
		document *Document
		expected string
	}{
		{
			document: &Document{
				Text:           "line 1\nline 2\nline 3\nline 4\n",
				cursorPosition: istrings.RuneCountInString("line 1\n" + "lin"),
			},
			expected: "line 1\nlin",
		},
		{
			document: &Document{
				Text:           "あいうえお\nかきくけこ\nさしすせそ\nたちつてと\n",
				cursorPosition: 8,
			},
			expected: "あいうえお\nかき",
		},
		{
			document: &Document{
				Text:           "Добрый\nдень\nДобрый день",
				cursorPosition: 9,
			},
			expected: "Добрый\nде",
		},
	}
	for i, p := range patterns {
		ac := p.document.TextBeforeCursor()
		if ac != p.expected {
			t.Errorf("[%d] Should be %s, got %s", i, p.expected, ac)
		}
	}
}

func TestDocument_LastLineIndentLevel(t *testing.T) {
	tests := []struct {
		document   *Document
		indentSize int
		want       int
	}{
		{
			document: &Document{
				Text: "line 1\nline 2\nline 3\n      line 4",
			},
			indentSize: 2,
			want:       3,
		},
		{
			document: &Document{
				Text: "line 1\nline 2\nline 3\n      line 4",
			},
			indentSize: 1,
			want:       6,
		},
		{
			document: &Document{
				Text: "line 1\nline 2\nline 3\n      line 4",
			},
			indentSize: 4,
			want:       1,
		},
		{
			document: &Document{
				Text: "line 1\nline 2\nline 3\n      line 4",
			},
			indentSize: 0,
			want:       0,
		},
	}

	for i, tc := range tests {
		got := tc.document.LastLineIndentLevel(tc.indentSize)
		if got != tc.want {
			t.Errorf("[%d] Should be %#v, got %#v", i, tc.want, got)
		}
	}
}

func TestDocument_PreviousLine(t *testing.T) {
	tests := []struct {
		document *Document
		want     string
		ok       bool
	}{
		{
			document: &Document{
				Text:           "line 1\nline 2\nline 3\n      line 4",
				cursorPosition: 32,
			},
			want: "line 3",
			ok:   true,
		},
		{
			document: &Document{
				Text:           "line 1\nline 2\nline 3\n      line 4",
				cursorPosition: 19,
			},
			want: "line 2",
			ok:   true,
		},
		{
			document: &Document{
				Text:           "line 1\nline 2\nline 3\n      line 4",
				cursorPosition: 1,
			},
			want: "",
			ok:   false,
		},
	}

	for i, tc := range tests {
		got, ok := tc.document.PreviousLine()
		if got != tc.want || ok != tc.ok {
			t.Errorf("[%d] Should be (%#v, %#v), got (%#v, %#v)", i, tc.want, tc.ok, got, ok)
		}
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
				cursorPosition: istrings.RuneCountInString("line 1\n" + "lin"),
			},
			expected: "e 2\nline 3\nline 4\n",
		},
		{
			document: &Document{
				Text:           "",
				cursorPosition: 0,
			},
			expected: "",
		},
		{
			document: &Document{
				Text:           "あいうえお\nかきくけこ\nさしすせそ\nたちつてと\n",
				cursorPosition: 8,
			},
			expected: "くけこ\nさしすせそ\nたちつてと\n",
		},
		{
			document: &Document{
				Text:           "Добрый\nдень\nДобрый день",
				cursorPosition: 9,
			},
			expected: "нь\nДобрый день",
		},
	}

	for i, p := range pattern {
		ac := p.document.TextAfterCursor()
		if ac != p.expected {
			t.Errorf("[%d] Should be %#v, got %#v", i, p.expected, ac)
		}
	}
}

func TestDocument_GetWordBeforeCursor(t *testing.T) {
	pattern := []struct {
		document *Document
		expected string
		sep      string
	}{
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: istrings.RuneCountInString("apple bana"),
			},
			expected: "bana",
		},
		{
			document: &Document{
				Text:           "apply -f ./file/foo.json",
				cursorPosition: istrings.RuneCountInString("apply -f ./file/foo.json"),
			},
			expected: "foo.json",
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "apple banana orange",
				cursorPosition: istrings.RuneCountInString("apple ba"),
			},
			expected: "ba",
		},
		{
			document: &Document{
				Text:           "apply -f ./file/foo.json",
				cursorPosition: istrings.RuneCountInString("apply -f ./fi"),
			},
			expected: "fi",
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "apple ",
				cursorPosition: istrings.RuneCountInString("apple "),
			},
			expected: "",
		},
		{
			document: &Document{
				Text:           "あいうえお かきくけこ さしすせそ",
				cursorPosition: 8,
			},
			expected: "かき",
		},
		{
			document: &Document{
				Text:           "Добрый день Добрый день",
				cursorPosition: 9,
			},
			expected: "де",
		},
	}

	for i, p := range pattern {
		if p.sep == "" {
			ac := p.document.GetWordBeforeCursor()
			if ac != p.expected {
				t.Errorf("[%d] Should be %#v, got %#v", i, p.expected, ac)
			}
			ac = p.document.GetWordBeforeCursorUntilSeparator("")
			if ac != p.expected {
				t.Errorf("[%d] Should be %#v, got %#v", i, p.expected, ac)
			}
		} else {
			ac := p.document.GetWordBeforeCursorUntilSeparator(p.sep)
			if ac != p.expected {
				t.Errorf("[%d] Should be %#v, got %#v", i, p.expected, ac)
			}
		}
	}
}

func TestDocument_GetWordBeforeCursorWithSpace(t *testing.T) {
	pattern := []struct {
		document *Document
		expected string
		sep      string
	}{
		{
			document: &Document{
				Text:           "apple bana ",
				cursorPosition: istrings.RuneCountInString("apple bana "),
			},
			expected: "bana ",
		},
		{
			document: &Document{
				Text:           "apply -f /path/to/file/",
				cursorPosition: istrings.RuneCountInString("apply -f /path/to/file/"),
			},
			expected: "file/",
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "apple ",
				cursorPosition: istrings.RuneCountInString("apple "),
			},
			expected: "apple ",
		},
		{
			document: &Document{
				Text:           "path/",
				cursorPosition: istrings.RuneCountInString("path/"),
			},
			expected: "path/",
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "あいうえお かきくけこ ",
				cursorPosition: 12,
			},
			expected: "かきくけこ ",
		},
		{
			document: &Document{
				Text:           "Добрый день ",
				cursorPosition: 12,
			},
			expected: "день ",
		},
	}

	for _, p := range pattern {
		if p.sep == "" {
			ac := p.document.GetWordBeforeCursorWithSpace()
			if ac != p.expected {
				t.Errorf("Should be %#v, got %#v", p.expected, ac)
			}
			ac = p.document.GetWordBeforeCursorUntilSeparatorIgnoreNextToCursor("")
			if ac != p.expected {
				t.Errorf("Should be %#v, got %#v", p.expected, ac)
			}
		} else {
			ac := p.document.GetWordBeforeCursorUntilSeparatorIgnoreNextToCursor(p.sep)
			if ac != p.expected {
				t.Errorf("Should be %#v, got %#v", p.expected, ac)
			}
		}
	}
}

func TestDocument_FindStartOfPreviousWord(t *testing.T) {
	pattern := []struct {
		document *Document
		expected istrings.ByteNumber
		sep      string
	}{
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: istrings.RuneCountInString("apple bana"),
			},
			expected: istrings.Len("apple "),
		},
		{
			document: &Document{
				Text:           "apply -f ./file/foo.json",
				cursorPosition: istrings.RuneCountInString("apply -f ./file/foo.json"),
			},
			expected: istrings.Len("apply -f ./file/"),
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "apple ",
				cursorPosition: istrings.RuneCountInString("apple "),
			},
			expected: istrings.Len("apple "),
		},
		{
			document: &Document{
				Text:           "apply -f ./file/foo.json",
				cursorPosition: istrings.RuneCountInString("apply -f ./"),
			},
			expected: istrings.Len("apply -f ./"),
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "あいうえお かきくけこ さしすせそ",
				cursorPosition: 8, // between 'き' and 'く'
			},
			expected: istrings.Len("あいうえお "), // this function returns index byte in string
		},
		{
			document: &Document{
				Text:           "Добрый день Добрый день",
				cursorPosition: 9,
			},
			expected: istrings.Len("Добрый "), // this function returns index byte in string
		},
	}

	for _, p := range pattern {
		if p.sep == "" {
			ac := p.document.FindStartOfPreviousWord()
			if ac != p.expected {
				t.Errorf("Should be %#v, got %#v", p.expected, ac)
			}
			ac = p.document.FindStartOfPreviousWordUntilSeparator("")
			if ac != p.expected {
				t.Errorf("Should be %#v, got %#v", p.expected, ac)
			}
		} else {
			ac := p.document.FindStartOfPreviousWordUntilSeparator(p.sep)
			if ac != p.expected {
				t.Errorf("Should be %#v, got %#v", p.expected, ac)
			}
		}
	}
}

func TestDocument_FindStartOfPreviousWordWithSpace(t *testing.T) {
	pattern := []struct {
		document *Document
		expected istrings.ByteNumber
		sep      string
	}{
		{
			document: &Document{
				Text:           "apple bana ",
				cursorPosition: istrings.RuneCountInString("apple bana "),
			},
			expected: istrings.Len("apple "),
		},
		{
			document: &Document{
				Text:           "apply -f /file/foo/",
				cursorPosition: istrings.RuneCountInString("apply -f /file/foo/"),
			},
			expected: istrings.Len("apply -f /file/"),
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "apple ",
				cursorPosition: istrings.RuneCountInString("apple "),
			},
			expected: istrings.Len(""),
		},
		{
			document: &Document{
				Text:           "file/",
				cursorPosition: istrings.RuneCountInString("file/"),
			},
			expected: istrings.Len(""),
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "あいうえお かきくけこ ",
				cursorPosition: 12, // cursor points to last
			},
			expected: istrings.Len("あいうえお "), // this function returns index byte in string
		},
		{
			document: &Document{
				Text:           "Добрый день ",
				cursorPosition: 12,
			},
			expected: istrings.Len("Добрый "), // this function returns index byte in string
		},
	}

	for _, p := range pattern {
		if p.sep == "" {
			ac := p.document.FindStartOfPreviousWordWithSpace()
			if ac != p.expected {
				t.Errorf("Should be %#v, got %#v", p.expected, ac)
			}
			ac = p.document.FindStartOfPreviousWordUntilSeparatorIgnoreNextToCursor("")
			if ac != p.expected {
				t.Errorf("Should be %#v, got %#v", p.expected, ac)
			}
		} else {
			ac := p.document.FindStartOfPreviousWordUntilSeparatorIgnoreNextToCursor(p.sep)
			if ac != p.expected {
				t.Errorf("Should be %#v, got %#v", p.expected, ac)
			}
		}
	}
}

func TestDocument_GetWordAfterCursor(t *testing.T) {
	pattern := []struct {
		document *Document
		expected string
		sep      string
	}{
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: istrings.RuneCountInString("apple bana"),
			},
			expected: "",
		},
		{
			document: &Document{
				Text:           "apply -f ./file/foo.json",
				cursorPosition: istrings.RuneCountInString("apply -f ./fi"),
			},
			expected: "le",
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: istrings.RuneCountInString("apple "),
			},
			expected: "bana",
		},
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: istrings.RuneCountInString("apple"),
			},
			expected: "",
		},
		{
			document: &Document{
				Text:           "apply -f ./file/foo.json",
				cursorPosition: istrings.RuneCountInString("apply -f ."),
			},
			expected: "",
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: istrings.RuneCountInString("ap"),
			},
			expected: "ple",
		},
		{
			document: &Document{
				Text:           "あいうえお かきくけこ さしすせそ",
				cursorPosition: 8,
			},
			expected: "くけこ",
		},
		{
			document: &Document{
				Text:           "Добрый день Добрый день",
				cursorPosition: 9,
			},
			expected: "нь",
		},
	}

	for k, p := range pattern {
		if p.sep == "" {
			ac := p.document.GetWordAfterCursor()
			if ac != p.expected {
				t.Errorf("[%d] Should be %#v, got %#v", k, p.expected, ac)
			}
			ac = p.document.GetWordAfterCursorUntilSeparator("")
			if ac != p.expected {
				t.Errorf("[%d] Should be %#v, got %#v", k, p.expected, ac)
			}
		} else {
			ac := p.document.GetWordAfterCursorUntilSeparator(p.sep)
			if ac != p.expected {
				t.Errorf("[%d] Should be %#v, got %#v", k, p.expected, ac)
			}
		}
	}
}

func TestDocument_GetWordAfterCursorWithSpace(t *testing.T) {
	pattern := []struct {
		document *Document
		expected string
		sep      string
	}{
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: istrings.RuneCountInString("apple bana"),
			},
			expected: "",
		},
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: istrings.RuneCountInString("apple "),
			},
			expected: "bana",
		},
		{
			document: &Document{
				Text:           "/path/to",
				cursorPosition: istrings.RuneCountInString("/path/"),
			},
			expected: "to",
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "/path/to/file",
				cursorPosition: istrings.RuneCountInString("/path/"),
			},
			expected: "to",
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: istrings.RuneCountInString("apple"),
			},
			expected: " bana",
		},
		{
			document: &Document{
				Text:           "path/to",
				cursorPosition: istrings.RuneCountInString("path"),
			},
			expected: "/to",
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: istrings.RuneCountInString("ap"),
			},
			expected: "ple",
		},
		{
			document: &Document{
				Text:           "あいうえお かきくけこ さしすせそ",
				cursorPosition: 5,
			},
			expected: " かきくけこ",
		},
		{
			document: &Document{
				Text:           "Добрый день Добрый день",
				cursorPosition: 6,
			},
			expected: " день",
		},
	}

	for k, p := range pattern {
		if p.sep == "" {
			ac := p.document.GetWordAfterCursorWithSpace()
			if ac != p.expected {
				t.Errorf("[%d] Should be %#v, got %#v", k, p.expected, ac)
			}
			ac = p.document.GetWordAfterCursorUntilSeparatorIgnoreNextToCursor("")
			if ac != p.expected {
				t.Errorf("[%d] Should be %#v, got %#v", k, p.expected, ac)
			}
		} else {
			ac := p.document.GetWordAfterCursorUntilSeparatorIgnoreNextToCursor(p.sep)
			if ac != p.expected {
				t.Errorf("[%d] Should be %#v, got %#v", k, p.expected, ac)
			}
		}
	}
}

func TestDocument_FindEndOfCurrentWord(t *testing.T) {
	pattern := []struct {
		document *Document
		expected istrings.ByteNumber
		sep      string
	}{
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: istrings.RuneCountInString("apple bana"),
			},
			expected: istrings.Len(""),
		},
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: istrings.RuneCountInString("apple "),
			},
			expected: istrings.Len("bana"),
		},
		{
			document: &Document{
				Text:           "apply -f ./file/foo.json",
				cursorPosition: istrings.RuneCountInString("apply -f ./"),
			},
			expected: istrings.Len("file"),
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: istrings.RuneCountInString("apple"),
			},
			expected: istrings.Len(""),
		},
		{
			document: &Document{
				Text:           "apply -f ./file/foo.json",
				cursorPosition: istrings.RuneCountInString("apply -f ."),
			},
			expected: istrings.Len(""),
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: istrings.RuneCountInString("ap"),
			},
			expected: istrings.Len("ple"),
		},
		{
			// りん(cursor)ご ばなな
			document: &Document{
				Text:           "りんご ばなな",
				cursorPosition: 2,
			},
			expected: istrings.Len("ご"),
		},
		{
			document: &Document{
				Text:           "りんご ばなな",
				cursorPosition: 3,
			},
			expected: 0,
		},
		{
			// Доб(cursor)рый день
			document: &Document{
				Text:           "Добрый день",
				cursorPosition: 3,
			},
			expected: istrings.Len("рый"),
		},
	}

	for k, p := range pattern {
		if p.sep == "" {
			ac := p.document.FindEndOfCurrentWord()
			if ac != p.expected {
				t.Errorf("[%d] Should be %#v, got %#v", k, p.expected, ac)
			}
			ac = p.document.FindEndOfCurrentWordUntilSeparator("")
			if ac != p.expected {
				t.Errorf("[%d] Should be %#v, got %#v", k, p.expected, ac)
			}
		} else {
			ac := p.document.FindEndOfCurrentWordUntilSeparator(p.sep)
			if ac != p.expected {
				t.Errorf("[%d] Should be %#v, got %#v", k, p.expected, ac)
			}
		}
	}
}

func TestDocument_FindEndOfCurrentWordWithSpace(t *testing.T) {
	pattern := []struct {
		document *Document
		expected istrings.ByteNumber
		sep      string
	}{
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: istrings.RuneCountInString("apple bana"),
			},
			expected: istrings.Len(""),
		},
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: istrings.RuneCountInString("apple "),
			},
			expected: istrings.Len("bana"),
		},
		{
			document: &Document{
				Text:           "apply -f /file/foo.json",
				cursorPosition: istrings.RuneCountInString("apply -f /"),
			},
			expected: istrings.Len("file"),
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: istrings.RuneCountInString("apple"),
			},
			expected: istrings.Len(" bana"),
		},
		{
			document: &Document{
				Text:           "apply -f /path/to",
				cursorPosition: istrings.RuneCountInString("apply -f /path"),
			},
			expected: istrings.Len("/to"),
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: istrings.RuneCountInString("ap"),
			},
			expected: istrings.Len("ple"),
		},
		{
			document: &Document{
				Text:           "あいうえお かきくけこ",
				cursorPosition: 6,
			},
			expected: istrings.Len("かきくけこ"),
		},
		{
			document: &Document{
				Text:           "あいうえお かきくけこ",
				cursorPosition: 5,
			},
			expected: istrings.Len(" かきくけこ"),
		},
		{
			document: &Document{
				Text:           "Добрый день",
				cursorPosition: 6,
			},
			expected: istrings.Len(" день"),
		},
	}

	for k, p := range pattern {
		if p.sep == "" {
			ac := p.document.FindEndOfCurrentWordWithSpace()
			if ac != p.expected {
				t.Errorf("[%d] Should be %#v, got %#v", k, p.expected, ac)
			}
			ac = p.document.FindEndOfCurrentWordUntilSeparatorIgnoreNextToCursor("")
			if ac != p.expected {
				t.Errorf("[%d] Should be %#v, got %#v", k, p.expected, ac)
			}
		} else {
			ac := p.document.FindEndOfCurrentWordUntilSeparatorIgnoreNextToCursor(p.sep)
			if ac != p.expected {
				t.Errorf("[%d] Should be %#v, got %#v", k, p.expected, ac)
			}
		}
	}
}

func TestDocument_CurrentLineBeforeCursor(t *testing.T) {
	d := &Document{
		Text:           "line 1\nline 2\nline 3\nline 4\n",
		cursorPosition: istrings.RuneCountInString("line 1\n" + "lin"),
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
		cursorPosition: istrings.RuneCountInString("line 1\n" + "lin"),
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
		cursorPosition: istrings.RuneCountInString("line 1\n" + "lin"),
	}
	ac := d.CurrentLine()
	ex := "line 2"
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}
}

func TestDocument_CursorPositionRowAndCol(t *testing.T) {
	var cursorPositionTests = []struct {
		document    *Document
		expectedRow istrings.RuneNumber
		expectedCol istrings.Width
	}{
		{
			document:    &Document{Text: "line 1\nline 2\nline 3\n", cursorPosition: istrings.RuneCountInString("line 1\n" + "lin")},
			expectedRow: 1,
			expectedCol: 3,
		},
		{
			document:    &Document{Text: "", cursorPosition: 0},
			expectedRow: 0,
			expectedCol: 0,
		},
	}
	for _, test := range cursorPositionTests {
		r := test.document.CursorPositionRow()
		if r != test.expectedRow {
			t.Errorf("Should be %#v, got %#v", test.expectedRow, r)
		}
		c := test.document.CursorPositionCol()
		if c != test.expectedCol {
			t.Errorf("Should be %#v, got %#v", test.expectedCol, c)
		}
	}
}

func TestDocument_GetCursorLeftPosition(t *testing.T) {
	d := &Document{
		Text:           "line 1\nline 2\nline 3\nline 4\n",
		cursorPosition: istrings.RuneCountInString("line 1\n" + "line 2\n" + "lin"),
	}
	ac := d.GetCursorLeftPosition(2)
	var ex istrings.RuneNumber = -2
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}
	ac = d.GetCursorLeftPosition(10)
	ex = -10
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}

	ac = d.GetCursorLeftPosition(30)
	ex = -17
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}
}

func TestDocument_GetCursorUpPosition(t *testing.T) {
	d := &Document{
		Text:           "line 1\nline 2\nline 3\nline 4\n",
		cursorPosition: istrings.RuneCountInString("line 1\n" + "line 2\n" + "lin"),
	}
	ac := d.GetCursorUpPosition(2, -1)
	ex := istrings.RuneCountInString("lin") - istrings.RuneCountInString("line 1\n"+"line 2\n"+"lin")
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}

	ac = d.GetCursorUpPosition(100, -1)
	ex = istrings.RuneCountInString("lin") - istrings.RuneCountInString("line 1\n"+"line 2\n"+"lin")
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}
}

func TestDocument_GetCursorDownPosition(t *testing.T) {
	d := &Document{
		Text:           "line 1\nline 2\nline 3\nline 4\n",
		cursorPosition: istrings.RuneCountInString("lin"),
	}
	ac := d.GetCursorDownPosition(2, -1)
	ex := istrings.RuneCountInString("line 1\n"+"line 2\n"+"lin") - istrings.RuneCountInString("lin")
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}

	ac = d.GetCursorDownPosition(100, -1)
	ex = istrings.RuneCountInString("line 1\n"+"line 2\n"+"line 3\n"+"line 4\n") - istrings.RuneCountInString("lin")
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}
}

func TestDocument_GetCursorLeftPositionBytes(t *testing.T) {
	const text = "a日bうc"
	// byte offsets: a=1, 日=3, b=1, う=3, c=1
	d := &Document{
		Text:           text,
		cursorPosition: istrings.RuneCountInString(text),
	}

	if ac := d.GetCursorLeftPositionBytes(0); ac != 0 {
		t.Errorf("zero move: got %#v, want 0", ac)
	}
	if ac := d.GetCursorLeftPositionBytes(1); ac != -1 {
		t.Errorf("one ASCII byte: got %#v, want -1", ac)
	}
	// 1 byte of 'c' plus 1 byte into 'う' consumes the whole rune.
	if ac := d.GetCursorLeftPositionBytes(2); ac != -2 {
		t.Errorf("byte inside multibyte rune: got %#v, want -2", ac)
	}
	if ac := d.GetCursorLeftPositionBytes(1+3); ac != -2 {
		t.Errorf("exact rune boundary: got %#v, want -2", ac)
	}
	if ac := d.GetCursorLeftPositionBytes(100); ac != -istrings.RuneCountInString(text) {
		t.Errorf("clamp to start: got %#v, want %d", ac, -istrings.RuneCountInString(text))
	}

	d.cursorPosition = 0
	if ac := d.GetCursorLeftPositionBytes(1); ac != 0 {
		t.Errorf("already at start: got %#v, want 0", ac)
	}

	d.cursorPosition = istrings.RuneCountInString("a日")
	if ac := d.GetCursorLeftPositionBytes(-1); ac != 1 {
		t.Errorf("negative count moves right: got %#v, want 1", ac)
	}
}

func TestDocument_GetCursorRightPositionBytes(t *testing.T) {
	const text = "a日bうc"
	d := &Document{
		Text:           text,
		cursorPosition: 0,
	}

	if ac := d.GetCursorRightPositionBytes(0); ac != 0 {
		t.Errorf("zero move: got %#v, want 0", ac)
	}
	if ac := d.GetCursorRightPositionBytes(1); ac != 1 {
		t.Errorf("one ASCII byte: got %#v, want 1", ac)
	}
	// 1 byte of 'a' plus 1 byte into '日' consumes the whole rune.
	if ac := d.GetCursorRightPositionBytes(2); ac != 2 {
		t.Errorf("byte inside multibyte rune: got %#v, want 2", ac)
	}
	if ac := d.GetCursorRightPositionBytes(1+3); ac != 2 {
		t.Errorf("exact rune boundary: got %#v, want 2", ac)
	}
	if ac := d.GetCursorRightPositionBytes(100); ac != istrings.RuneCountInString(text) {
		t.Errorf("clamp to end: got %#v, want %d", ac, istrings.RuneCountInString(text))
	}

	d.cursorPosition = istrings.RuneCountInString(text)
	if ac := d.GetCursorRightPositionBytes(1); ac != 0 {
		t.Errorf("already at end: got %#v, want 0", ac)
	}

	d.cursorPosition = istrings.RuneCountInString("a日")
	if ac := d.GetCursorRightPositionBytes(-1); ac != -1 {
		t.Errorf("negative count moves left: got %#v, want -1", ac)
	}
}

func TestDocument_CursorByteOffset(t *testing.T) {
	d := &Document{
		Text:           "a日b",
		cursorPosition: istrings.RuneCountInString("a日"),
	}
	if ac := d.CursorByteOffset(); ac != istrings.ByteNumber(len("a日")) {
		t.Errorf("got %#v, want %d", ac, len("a日"))
	}

	d.cursorPosition = 0
	if ac := d.CursorByteOffset(); ac != 0 {
		t.Errorf("cursor at start: got %#v, want 0", ac)
	}
}

func TestDocument_GetCursorRightPosition(t *testing.T) {
	d := &Document{
		Text:           "line 1\nline 2\nline 3\nline 4\n",
		cursorPosition: istrings.RuneCountInString("line 1\n" + "line 2\n" + "lin"),
	}
	ac := d.GetCursorRightPosition(2)
	var ex istrings.RuneNumber = 2
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}
	ac = d.GetCursorRightPosition(10)
	ex = 10
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}

	ac = d.GetCursorRightPosition(30)
	ex = 11
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}
}

func TestDocument_Lines(t *testing.T) {
	d := &Document{
		Text:           "line 1\nline 2\nline 3\nline 4\n",
		cursorPosition: istrings.RuneCountInString("line 1\n" + "lin"),
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
		cursorPosition: istrings.RuneCountInString("line 1\n" + "lin"),
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
		cursorPosition: istrings.RuneCountInString("line 1\n" + "lin"),
	}
	row, col := d.TranslateIndexToPosition(istrings.RuneCountInString("line 1\nline 2\nlin"))
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
		cursorPosition: istrings.RuneCountInString("line 1\n" + "lin"),
	}
	ac := d.TranslateRowColToIndex(2, 3)
	ex := istrings.RuneCountInString("line 1\nline 2\nlin")
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
		cursorPosition: istrings.RuneCountInString("line 1\nline"),
	}
	ac := d.OnLastLine()
	if ac {
		t.Errorf("Should be %#v, got %#v", false, ac)
	}
	d.cursorPosition = istrings.RuneCountInString("line 1\nline 2\nline")
	ac = d.OnLastLine()
	if !ac {
		t.Errorf("Should be %#v, got %#v", true, ac)
	}
}

func TestDocument_GetEndOfLinePosition(t *testing.T) {
	d := &Document{
		Text:           "line 1\nline 2\nline 3",
		cursorPosition: istrings.RuneCountInString("line 1\nli"),
	}
	ac := d.GetEndOfLinePosition()
	ex := istrings.RuneCountInString("ne 2")
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}
}

func TestDocument_WordPartCursorStops(t *testing.T) {
	const text = "foo.bar(elo, dupa)"
	// Ends of foo, ., bar, (, elo, ,, dupa, and ). The space is not a stop.
	stops := []istrings.RuneNumber{3, 4, 7, 8, 11, 12, 17, 18}

	d := &Document{Text: text}
	for _, stop := range stops {
		n := d.FindRuneNumberUntilEndOfCurrentWord()
		d.cursorPosition += n
		if d.cursorPosition != stop {
			t.Fatalf("forward: cursor %d, want %d", d.cursorPosition, stop)
		}
	}
	if n := d.FindRuneNumberUntilEndOfCurrentWord(); n != 0 {
		t.Fatalf("forward at end: moved %d", n)
	}

	for i := len(stops) - 2; i >= 0; i-- {
		n := d.FindRuneNumberUntilStartOfPreviousWord()
		d.cursorPosition -= n
		if d.cursorPosition != stops[i] {
			t.Fatalf("backward: cursor %d, want %d", d.cursorPosition, stops[i])
		}
	}
	n := d.FindRuneNumberUntilStartOfPreviousWord()
	d.cursorPosition -= n
	if d.cursorPosition != 0 {
		t.Fatalf("backward to start: cursor %d", d.cursorPosition)
	}
}

func TestDocument_WordPartBoundaries(t *testing.T) {
	const text = "foo.bar(elo, dupa)"
	parts := []struct {
		start istrings.RuneNumber
		end   istrings.RuneNumber
		text  string
	}{
		{0, 3, "foo"},
		{3, 4, "."},
		{4, 7, "bar"},
		{7, 8, "("},
		{8, 11, "elo"},
		{11, 12, ","},
		{13, 17, "dupa"},
		{17, 18, ")"},
	}

	for _, part := range parts {
		d := &Document{Text: text, cursorPosition: part.end}
		if got := d.GetWordBeforeCursor(); got != part.text {
			t.Errorf("before %q: got %q", part.text, got)
		}
		d.cursorPosition = part.start
		if got := d.GetWordAfterCursor(); got != part.text {
			t.Errorf("after %q: got %q", part.text, got)
		}
	}

	d := &Document{Text: text, cursorPosition: 13}
	if got := d.GetWordBeforeCursor(); got != "" {
		t.Errorf("before space: got %q", got)
	}
	if got := d.GetWordBeforeCursorWithSpace(); got != ", " {
		t.Errorf("before space with space: got %q", got)
	}
	d.cursorPosition = 12
	if got := d.GetWordAfterCursor(); got != "" {
		t.Errorf("after comma: got %q", got)
	}
	if got := d.GetWordAfterCursorWithSpace(); got != " dupa" {
		t.Errorf("after comma with space: got %q", got)
	}
}
