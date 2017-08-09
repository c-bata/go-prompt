package prompt

import (
	"reflect"
	"testing"
)

func TestFilter(t *testing.T) {
	var scenarioTable = []struct {
		scenario   string
		filter     Filter
		list       []Suggest
		substr     string
		ignoreCase bool
		expected   []Suggest
	}{
		{
			scenario: "Contains don't ignore case",
			filter:   FilterContains,
			list: []Suggest{
				{Text: "abcde"},
				{Text: "fghij"},
				{Text: "ABCDE"},
			},
			substr:     "cd",
			ignoreCase: false,
			expected: []Suggest{
				{Text: "abcde"},
			},
		},
		{
			scenario: "Contains ignore case",
			filter:   FilterContains,
			list: []Suggest{
				{Text: "abcde"},
				{Text: "fghij"},
				{Text: "ABCDE"},
			},
			substr:     "cd",
			ignoreCase: true,
			expected: []Suggest{
				{Text: "abcde"},
				{Text: "ABCDE"},
			},
		},
		{
			scenario: "HasPrefix don't ignore case",
			filter:   FilterHasPrefix,
			list: []Suggest{
				{Text: "abcde"},
				{Text: "fghij"},
				{Text: "ABCDE"},
			},
			substr:     "abc",
			ignoreCase: false,
			expected: []Suggest{
				{Text: "abcde"},
			},
		},
		{
			scenario: "HasPrefix ignore case",
			filter:   FilterHasPrefix,
			list: []Suggest{
				{Text: "abcde"},
				{Text: "fabcj"},
				{Text: "ABCDE"},
			},
			substr:     "abc",
			ignoreCase: true,
			expected: []Suggest{
				{Text: "abcde"},
				{Text: "ABCDE"},
			},
		},
		{
			scenario: "HasSuffix don't ignore case",
			filter:   FilterHasSuffix,
			list: []Suggest{
				{Text: "abcde"},
				{Text: "fcdej"},
				{Text: "ABCDE"},
			},
			substr:     "cde",
			ignoreCase: false,
			expected: []Suggest{
				{Text: "abcde"},
			},
		},
		{
			scenario: "HasSuffix ignore case",
			filter:   FilterHasSuffix,
			list: []Suggest{
				{Text: "abcde"},
				{Text: "fcdej"},
				{Text: "ABCDE"},
			},
			substr:     "cde",
			ignoreCase: true,
			expected: []Suggest{
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
