package prompt

import (
	"testing"
	"reflect"
)

func TestFilter(t *testing.T) {
	var scenarioTable = [] struct {
		scenario   string
		filter     CompletionFilter
		list       []Completion
		substr     string
		ignoreCase bool
		expected   []Completion
	} {
		{
			scenario:   "Contains don't ignore case",
			filter:     FilterContains,
			list:       []Completion{
				{Text: "abcde"},
				{Text: "fghij"},
				{Text: "ABCDE"},
			},
			substr:     "cd",
			ignoreCase: false,
			expected:  []Completion{
				{Text: "abcde"},
			},
		},
		{
			scenario:   "Contains ignore case",
			filter:     FilterContains,
			list:       []Completion{
				{Text: "abcde"},
				{Text: "fghij"},
				{Text: "ABCDE"},
			},
			substr:     "cd",
			ignoreCase: true,
			expected:   []Completion{
				{Text: "abcde"},
				{Text: "ABCDE"},
			},
		},
		{
			scenario:   "HasPrefix don't ignore case",
			filter:     FilterHasPrefix,
			list:       []Completion{
				{Text: "abcde"},
				{Text: "fghij"},
				{Text: "ABCDE"},
			},
			substr:     "abc",
			ignoreCase: false,
			expected:  []Completion{
				{Text: "abcde"},
			},
		},
		{
			scenario:   "HasPrefix ignore case",
			filter:     FilterHasPrefix,
			list:       []Completion{
				{Text: "abcde"},
				{Text: "fabcj"},
				{Text: "ABCDE"},
			},
			substr:     "abc",
			ignoreCase: true,
			expected:   []Completion{
				{Text: "abcde"},
				{Text: "ABCDE"},
			},
		},
		{
			scenario:   "HasSuffix don't ignore case",
			filter:     FilterHasSuffix,
			list:       []Completion{
				{Text: "abcde"},
				{Text: "fcdej"},
				{Text: "ABCDE"},
			},
			substr:     "cde",
			ignoreCase: false,
			expected:  []Completion{
				{Text: "abcde"},
			},
		},
		{
			scenario:   "HasSuffix ignore case",
			filter:     FilterHasSuffix,
			list:       []Completion{
				{Text: "abcde"},
				{Text: "fcdej"},
				{Text: "ABCDE"},
			},
			substr:     "cde",
			ignoreCase: true,
			expected:   []Completion{
				{Text: "abcde"},
				{Text: "ABCDE"},
			},
		},
	}

	for _, s := range scenarioTable {
		if actual := s.filter(s.list, s.substr, s.ignoreCase); !reflect.DeepEqual(actual, s.expected) {
			t.Errorf("%s: Should be %#v, but got %#v", s.scenario, s.expected, actual)
		}
	}
}
