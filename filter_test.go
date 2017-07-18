package prompt

import (
	"testing"
	"reflect"
)

func TestFilter(t *testing.T) {
	var scenarioTable = [] struct {
		scenario   string
		filter     Filter
		list       []string
		substr     string
		ignoreCase bool
		expected   []string
	} {
		{
			scenario:   "Contains don't ignore case",
			filter:     FilterContains,
			list:       []string{"abcde", "fghij", "ABCDE"},
			substr:     "cd",
			ignoreCase: false,
			expected:  []string{"abcde"},
		},
		{
			scenario:   "Contains ignore case",
			filter:     FilterContains,
			list:       []string{"abcde", "fghij", "ABCDE"},
			substr:     "cd",
			ignoreCase: true,
			expected:   []string{"abcde", "ABCDE"},
		},
		{
			scenario:   "HasPrefix don't ignore case",
			filter:     FilterHasPrefix,
			list:       []string{"abcde", "fghij", "ABCDE"},
			substr:     "abc",
			ignoreCase: false,
			expected:  []string{"abcde"},
		},
		{
			scenario:   "HasPrefix ignore case",
			filter:     FilterHasPrefix,
			list:       []string{"abcde", "fabcj", "ABCDE"},
			substr:     "abc",
			ignoreCase: true,
			expected:   []string{"abcde", "ABCDE"},
		},
		{
			scenario:   "HasSuffix don't ignore case",
			filter:     FilterHasSuffix,
			list:       []string{"abcde", "fcdej", "ABCDE"},
			substr:     "cde",
			ignoreCase: false,
			expected:  []string{"abcde"},
		},
		{
			scenario:   "HasSuffix ignore case",
			filter:     FilterHasSuffix,
			list:       []string{"abcde", "fcdej", "ABCDE"},
			substr:     "cde",
			ignoreCase: true,
			expected:   []string{"abcde", "ABCDE"},
		},
	}

	for _, s := range scenarioTable {
		if actual := s.filter(s.list, s.substr, s.ignoreCase); !reflect.DeepEqual(actual, s.expected) {
			t.Errorf("%s: Should be %#v, but got %#v", s.scenario, s.expected, actual)
		}
	}
}
