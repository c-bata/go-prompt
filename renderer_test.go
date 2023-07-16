//go:build !windows
// +build !windows

package prompt

import (
	"reflect"
	"syscall"
	"testing"

	istrings "github.com/elk-language/go-prompt/strings"
)

func TestFormatCompletion(t *testing.T) {
	scenarioTable := []struct {
		scenario      string
		completions   []Suggest
		prefix        string
		suffix        string
		expected      []Suggest
		maxWidth      istrings.Width
		expectedWidth istrings.Width
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
	r := NewRenderer()
	r.out = &PosixWriter{
		fd: syscall.Stdin, // "write" to stdin just so we don't mess with the output of the tests
	}
	r.col = 1
	b := NewBuffer()
	r.BreakLine(b, nil)

	if i != 0 {
		t.Errorf("i should initially be 0, before applying a break line callback")
	}

	r.breakLineCallback = func(doc *Document) {
		i++
	}
	r.BreakLine(b, nil)
	r.BreakLine(b, nil)
	r.BreakLine(b, nil)

	if i != 3 {
		t.Errorf("BreakLine callback not called, i should be 3")
	}
}

func TestGetMultilinePrefix(t *testing.T) {
	tests := map[string]struct {
		prefix string
		want   string
	}{
		"single width chars": {
			prefix: ">>",
			want:   "..",
		},
		"double width chars": {
			prefix: "本日",
			want:   "....",
		},
		"trailing spaces and single width chars": {
			prefix: ">!>   ",
			want:   "...   ",
		},
		"trailing spaces and double width chars": {
			prefix: "本日:   ",
			want:   ".....   ",
		},
		"leading spaces and single width chars": {
			prefix: "  ah:   ",
			want:   ".....   ",
		},
		"leading spaces and double width chars": {
			prefix: "  本日:   ",
			want:   ".......   ",
		},
	}

	r := NewRenderer()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := r.getMultilinePrefix(tc.prefix)
			if tc.want != got {
				t.Errorf("Expected %#v, but got %#v", tc.want, got)
			}
		})
	}
}
