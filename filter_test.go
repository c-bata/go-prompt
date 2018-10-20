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
		{
			scenario: "Fuzzy don't ignore case",
			filter:   FilterFuzzy,
			list: []Suggest{
				{Text: "abcde"},
				{Text: "fcdej"},
				{Text: "ABCDE"},
			},
			substr:     "ae",
			ignoreCase: false,
			expected: []Suggest{
				{Text: "abcde"},
			},
		},
		{
			scenario: "Fuzzy ignore case",
			filter:   FilterFuzzy,
			list: []Suggest{
				{Text: "abcde"},
				{Text: "fcdej"},
				{Text: "ABCDE"},
			},
			substr:     "ae",
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

func TestFuzzyMatch(t *testing.T) {
	tests := []struct {
		s     string
		sub   string
		match bool
	}{
		{"dog house", "dog", true},
		{"dog house", "", true},
		{"", "", true},
		{"this is much longer", "hhg", true},
		{"this is much longer", "hhhg", false},
		{"long", "longer", false},
		{"can we do unicode 文字 with this 今日", "文字今日", true},
		{"can we do unicode 文字 with this 今日", "d文字tt今日", true},
		{"can we do unicode 文字 with this 今日", "d文字ttt今日", false},
	}

	for _, test := range tests {
		if fuzzyMatch(test.s, test.sub) != test.match {
			t.Errorf("fuzzymatch, %s in %s: expected %v, got %v", test.sub, test.s, test.match, !test.match)
		}
	}
}
