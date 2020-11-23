// +build !windows

package prompt

import (
	"reflect"
	"syscall"
	"testing"

	fcolor "github.com/fatih/color"
)

func TestFormatCompletion(t *testing.T) {
	scenarioTable := []struct {
		scenario      string
		completions   []Suggest
		prefix        string
		suffix        string
		expected      []Suggest
		maxWidth      int
		expectedWidth int
	}{
		{
			scenario: "",
			completions: []Suggest{
				{Text: "select"},
				{Text: "from"},
				{Text: "insert"},
				{Text: "where"},
			},
			prefix: " ",
			suffix: " ",
			expected: []Suggest{
				{Text: " select "},
				{Text: " from   "},
				{Text: " insert "},
				{Text: " where  "},
			},
			maxWidth:      20,
			expectedWidth: 8,
		},
		{
			scenario: "",
			completions: []Suggest{
				{Text: "select", Description: "select description"},
				{Text: "from", Description: "from description"},
				{Text: "insert", Description: "insert description"},
				{Text: "where", Description: "where description"},
			},
			prefix: " ",
			suffix: " ",
			expected: []Suggest{
				{Text: " select ", Description: " select description "},
				{Text: " from   ", Description: " from description   "},
				{Text: " insert ", Description: " insert description "},
				{Text: " where  ", Description: " where description  "},
			},
			maxWidth:      40,
			expectedWidth: 28,
		},
	}

	for _, s := range scenarioTable {
		ac, width := formatSuggestions(s.completions, s.maxWidth)
		if !reflect.DeepEqual(ac, s.expected) {
			t.Errorf("Should be %#v, but got %#v", s.expected, ac)
		}
		if width != s.expectedWidth {
			t.Errorf("Should be %#v, but got %#v", s.expectedWidth, width)
		}
	}
}

func TestBreakLineCallback(t *testing.T) {
	var i int
	r := &Render{
		prefix: "> ",
		out: &PosixWriter{
			fd: syscall.Stdin, // "write" to stdin just so we don't mess with the output of the tests
		},
		livePrefixCallback:       func() (string, bool) { return "", false },
		prefixColor:              fcolor.New(fcolor.FgBlue),
		inputColor:               nil,
		previewSuggestionColor:   fcolor.New(fcolor.FgGreen),
		suggestionColor:          fcolor.New(fcolor.FgWhite, fcolor.BgCyan),
		selectedSuggestionColor:  fcolor.New(fcolor.FgBlack, fcolor.BgHiCyan),
		descriptionColor:         fcolor.New(fcolor.FgBlack, fcolor.BgHiCyan),
		selectedDescriptionColor: fcolor.New(fcolor.FgWhite, fcolor.BgCyan),
		scrollbarColor:           fcolor.New(fcolor.FgCyan),
		scrollbarThumbColor:      fcolor.New(fcolor.FgHiBlack),
		col:                      1,
	}
	b := NewBuffer()
	r.BreakLine(b)

	if i != 0 {
		t.Errorf("i should initially be 0, before applying a break line callback")
	}

	r.breakLineCallback = func(doc *Document) {
		i++
	}
	r.BreakLine(b)
	r.BreakLine(b)
	r.BreakLine(b)

	if i != 3 {
		t.Errorf("BreakLine callback not called, i should be 3")
	}
}
