package prompt

import (
	"fmt"
	"reflect"
	"testing"
	"unicode/utf8"
)

func ExampleDocument_CurrentLine() {
	d := &Document{
		Text: `Hello! my name is c-bata.
This is a example of Document component.
This component has texts displayed in terminal and cursor position.
`,
		cursorPosition: len(`Hello! my name is c-bata.
This is a exam`),
	}
	fmt.Println(d.CurrentLine())
	// Output:
	// This is a example of Document component.
}

func ExampleDocument_DisplayCursorPosition() {
	d := &Document{
		Text:           `Hello! my name is c-bata.`,
		cursorPosition: len(`Hello`),
	}
	fmt.Println("DisplayCursorPosition", d.DisplayCursorPosition())
	// Output:
	// DisplayCursorPosition 5
}

func ExampleDocument_CursorPositionRow() {
	d := &Document{
		Text: `Hello! my name is c-bata.
This is a example of Document component.
This component has texts displayed in terminal and cursor position.
`,
		cursorPosition: len(`Hello! my name is c-bata.
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
		cursorPosition: len(`Hello! my name is c-bata.
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
		cursorPosition: len(`Hello! my name is c-bata.
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
		cursorPosition: len(`Hello! my name is c-bata.
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
	fmt.Println("DisplayCursorPosition", d.DisplayCursorPosition())
	// Output:
	// DisplayCursorPosition 6
}

func ExampleDocument_CurrentLineBeforeCursor() {
	d := &Document{
		Text: `Hello! my name is c-bata.
This is a example of Document component.
This component has texts displayed in terminal and cursor position.
`,
		cursorPosition: len(`Hello! my name is c-bata.
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
		cursorPosition: len(`Hello! my name is c-bata.
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
		cursorPosition: len(`Hello! my name is c-bata.
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
		cursorPosition: len(`Hello! my name is c-bata.
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
		cursorPosition: len(`Hello! my name is c-bata.
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
		cursorPosition: len(`Hello! my name is c-bata.
This is a`),
	}
	fmt.Println(d.GetWordAfterCursorWithSpace())
	// Output:
	//  example
}

func ExampleDocument_GetWordBeforeCursorUntilSeparator() {
	d := &Document{
		Text:           `hello,i am c-bata`,
		cursorPosition: len(`hello,i am c`),
	}
	fmt.Println(d.GetWordBeforeCursorUntilSeparator(","))
	// Output:
	// i am c
}

func ExampleDocument_GetWordAfterCursorUntilSeparator() {
	d := &Document{
		Text:           `hello,i am c-bata,thank you for using go-prompt`,
		cursorPosition: len(`hello,i a`),
	}
	fmt.Println(d.GetWordAfterCursorUntilSeparator(","))
	// Output:
	// m c-bata
}

func ExampleDocument_GetWordBeforeCursorUntilSeparatorIgnoreNextToCursor() {
	d := &Document{
		Text:           `hello,i am c-bata,thank you for using go-prompt`,
		cursorPosition: len(`hello,i am c-bata,`),
	}
	fmt.Println(d.GetWordBeforeCursorUntilSeparatorIgnoreNextToCursor(","))
	// Output:
	// i am c-bata,
}

func ExampleDocument_GetWordAfterCursorUntilSeparatorIgnoreNextToCursor() {
	d := &Document{
		Text:           `hello,i am c-bata,thank you for using go-prompt`,
		cursorPosition: len(`hello`),
	}
	fmt.Println(d.GetWordAfterCursorUntilSeparatorIgnoreNextToCursor(","))
	// Output:
	// ,i am c-bata
}

func TestDocument_DisplayCursorPosition(t *testing.T) {
	patterns := []struct {
		document *Document
		expected int
	}{
		{
			document: &Document{
				Text:           "hello",
				cursorPosition: 2,
			},
			expected: 2,
		},
		{
			document: &Document{
				Text:           "こんにちは",
				cursorPosition: 2,
			},
			expected: 4,
		},
		{
			// If you're facing test failure on this test case and your terminal is iTerm2,
			// please check 'Profile -> Text' configuration. 'Use Unicode version 9 widths'
			// must be checked.
			// https://github.com/c-bata/go-prompt/pull/99
			document: &Document{
				Text:           "Добрый день",
				cursorPosition: 3,
			},
			expected: 3,
		},
	}

	for _, p := range patterns {
		ac := p.document.DisplayCursorPosition()
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
				cursorPosition: len([]rune("line 1\n" + "lin")),
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
				cursorPosition: len("line 1\n" + "lin"),
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

func TestDocument_TextAfterCursor(t *testing.T) {
	pattern := []struct {
		document *Document
		expected string
	}{
		{
			document: &Document{
				Text:           "line 1\nline 2\nline 3\nline 4\n",
				cursorPosition: len("line 1\n" + "lin"),
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
				cursorPosition: len("apple bana"),
			},
			expected: "bana",
		},
		{
			document: &Document{
				Text:           "apply -f ./file/foo.json",
				cursorPosition: len("apply -f ./file/foo.json"),
			},
			expected: "foo.json",
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "apple banana orange",
				cursorPosition: len("apple ba"),
			},
			expected: "ba",
		},
		{
			document: &Document{
				Text:           "apply -f ./file/foo.json",
				cursorPosition: len("apply -f ./fi"),
			},
			expected: "fi",
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "apple ",
				cursorPosition: len("apple "),
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
				cursorPosition: len("apple bana "),
			},
			expected: "bana ",
		},
		{
			document: &Document{
				Text:           "apply -f /path/to/file/",
				cursorPosition: len("apply -f /path/to/file/"),
			},
			expected: "file/",
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "apple ",
				cursorPosition: len("apple "),
			},
			expected: "apple ",
		},
		{
			document: &Document{
				Text:           "path/",
				cursorPosition: len("path/"),
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
		expected int
		sep      string
	}{
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: len("apple bana"),
			},
			expected: len("apple "),
		},
		{
			document: &Document{
				Text:           "apply -f ./file/foo.json",
				cursorPosition: len("apply -f ./file/foo.json"),
			},
			expected: len("apply -f ./file/"),
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "apple ",
				cursorPosition: len("apple "),
			},
			expected: len("apple "),
		},
		{
			document: &Document{
				Text:           "apply -f ./file/foo.json",
				cursorPosition: len("apply -f ./"),
			},
			expected: len("apply -f ./"),
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "あいうえお かきくけこ さしすせそ",
				cursorPosition: 8, // between 'き' and 'く'
			},
			expected: len("あいうえお "), // this function returns index byte in string
		},
		{
			document: &Document{
				Text:           "Добрый день Добрый день",
				cursorPosition: 9,
			},
			expected: len("Добрый "), // this function returns index byte in string
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
		expected int
		sep      string
	}{
		{
			document: &Document{
				Text:           "apple bana ",
				cursorPosition: len("apple bana "),
			},
			expected: len("apple "),
		},
		{
			document: &Document{
				Text:           "apply -f /file/foo/",
				cursorPosition: len("apply -f /file/foo/"),
			},
			expected: len("apply -f /file/"),
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "apple ",
				cursorPosition: len("apple "),
			},
			expected: len(""),
		},
		{
			document: &Document{
				Text:           "file/",
				cursorPosition: len("file/"),
			},
			expected: len(""),
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "あいうえお かきくけこ ",
				cursorPosition: 12, // cursor points to last
			},
			expected: len("あいうえお "), // this function returns index byte in string
		},
		{
			document: &Document{
				Text:           "Добрый день ",
				cursorPosition: 12,
			},
			expected: len("Добрый "), // this function returns index byte in string
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
				cursorPosition: len("apple bana"),
			},
			expected: "",
		},
		{
			document: &Document{
				Text:           "apply -f ./file/foo.json",
				cursorPosition: len("apply -f ./fi"),
			},
			expected: "le",
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: len("apple "),
			},
			expected: "bana",
		},
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: len("apple"),
			},
			expected: "",
		},
		{
			document: &Document{
				Text:           "apply -f ./file/foo.json",
				cursorPosition: len("apply -f ."),
			},
			expected: "",
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: len("ap"),
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
				cursorPosition: len("apple bana"),
			},
			expected: "",
		},
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: len("apple "),
			},
			expected: "bana",
		},
		{
			document: &Document{
				Text:           "/path/to",
				cursorPosition: len("/path/"),
			},
			expected: "to",
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "/path/to/file",
				cursorPosition: len("/path/"),
			},
			expected: "to",
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: len("apple"),
			},
			expected: " bana",
		},
		{
			document: &Document{
				Text:           "path/to",
				cursorPosition: len("path"),
			},
			expected: "/to",
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: len("ap"),
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
		expected int
		sep      string
	}{
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: len("apple bana"),
			},
			expected: len(""),
		},
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: len("apple "),
			},
			expected: len("bana"),
		},
		{
			document: &Document{
				Text:           "apply -f ./file/foo.json",
				cursorPosition: len("apply -f ./"),
			},
			expected: len("file"),
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: len("apple"),
			},
			expected: len(""),
		},
		{
			document: &Document{
				Text:           "apply -f ./file/foo.json",
				cursorPosition: len("apply -f ."),
			},
			expected: len(""),
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: len("ap"),
			},
			expected: len("ple"),
		},
		{
			// りん(cursor)ご ばなな
			document: &Document{
				Text:           "りんご ばなな",
				cursorPosition: 2,
			},
			expected: len("ご"),
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
			expected: len("рый"),
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
		expected int
		sep      string
	}{
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: len("apple bana"),
			},
			expected: len(""),
		},
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: len("apple "),
			},
			expected: len("bana"),
		},
		{
			document: &Document{
				Text:           "apply -f /file/foo.json",
				cursorPosition: len("apply -f /"),
			},
			expected: len("file"),
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: len("apple"),
			},
			expected: len(" bana"),
		},
		{
			document: &Document{
				Text:           "apply -f /path/to",
				cursorPosition: len("apply -f /path"),
			},
			expected: len("/to"),
			sep:      " /",
		},
		{
			document: &Document{
				Text:           "apple bana",
				cursorPosition: len("ap"),
			},
			expected: len("ple"),
		},
		{
			document: &Document{
				Text:           "あいうえお かきくけこ",
				cursorPosition: 6,
			},
			expected: len("かきくけこ"),
		},
		{
			document: &Document{
				Text:           "あいうえお かきくけこ",
				cursorPosition: 5,
			},
			expected: len(" かきくけこ"),
		},
		{
			document: &Document{
				Text:           "Добрый день",
				cursorPosition: 6,
			},
			expected: len(" день"),
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
		cursorPosition: len("line 1\n" + "lin"),
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
		cursorPosition: len("line 1\n" + "lin"),
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
		cursorPosition: len("line 1\n" + "lin"),
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
		expectedRow int
		expectedCol int
	}{
		{
			document:    &Document{Text: "line 1\nline 2\nline 3\n", cursorPosition: len("line 1\n" + "lin")},
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
		cursorPosition: len("line 1\n" + "line 2\n" + "lin"),
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
		cursorPosition: len("line 1\n" + "line 2\n" + "lin"),
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
		cursorPosition: len("lin"),
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
		cursorPosition: len("line 1\n" + "line 2\n" + "lin"),
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
		cursorPosition: len("line 1\n" + "lin"),
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
		cursorPosition: len("line 1\n" + "lin"),
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
		cursorPosition: len("line 1\n" + "lin"),
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
		cursorPosition: len("line 1\n" + "lin"),
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
		cursorPosition: len("line 1\nline"),
	}
	ac := d.OnLastLine()
	if ac {
		t.Errorf("Should be %#v, got %#v", false, ac)
	}
	d.cursorPosition = len("line 1\nline 2\nline")
	ac = d.OnLastLine()
	if !ac {
		t.Errorf("Should be %#v, got %#v", true, ac)
	}
}

func TestDocument_GetEndOfLinePosition(t *testing.T) {
	d := &Document{
		Text:           "line 1\nline 2\nline 3",
		cursorPosition: len("line 1\nli"),
	}
	ac := d.GetEndOfLinePosition()
	ex := len("ne 2")
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}
}
