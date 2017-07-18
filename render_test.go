package prompt

import (
	"reflect"
	"testing"
)

func TestFormatCompletion(t *testing.T) {
	scenarioTable := []struct {
		scenario      string
		completions   []Completion
		prefix        string
		suffix        string
		expected      []Completion
		maxWidth      int
		expectedWidth int
	}{
		{
			scenario: "",
			completions: []Completion{
				{Text: "select"},
				{Text: "from"},
				{Text: "insert"},
				{Text: "where"},
			},
			prefix: " ",
			suffix: " ",
			expected: []Completion{
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
			completions: []Completion{
				{Text: "select", Description: "select description"},
				{Text: "from", Description: "from description"},
				{Text: "insert", Description: "insert description"},
				{Text: "where", Description: "where description"},
			},
			prefix: " ",
			suffix: " ",
			expected: []Completion{
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
		ac, width := formatCompletions(s.completions, s.maxWidth)
		if !reflect.DeepEqual(ac, s.expected) {
			t.Errorf("Should be %#v, but got %#v", s.expected, ac)
		}
		if width != s.expectedWidth {
			t.Errorf("Should be %#v, but got %#v", s.expectedWidth, width)
		}
	}
}
