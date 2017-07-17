package prompt

import (
	"reflect"
	"testing"
)

func TestFormatCompletion(t *testing.T) {
	scenarioTable := []struct {
		scenario      string
		completions   []string
		prefix        string
		suffix        string
		expected      []string
		maxWidth      int
		expectedWidth int
	}{
		{
			scenario: "",
			completions: []string{
				"select",
				"from",
				"insert",
				"where",
			},
			prefix: " ",
			suffix: " ",
			expected: []string{
				" select ",
				" from   ",
				" insert ",
				" where  ",
			},
			maxWidth:      20,
			expectedWidth: 8,
		},
	}

	for _, s := range scenarioTable {
		ac, width := formatCompletions(s.completions, s.maxWidth, s.prefix, s.suffix)
		if !reflect.DeepEqual(ac, s.expected) {
			t.Errorf("Should be %#v, but got %#v", s.expected, ac)
		}
		if width != s.expectedWidth {
			t.Errorf("Should be %#v, but got %#v", s.expectedWidth, width)
		}
	}
}
