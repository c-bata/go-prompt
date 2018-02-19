package prompt

import (
	"reflect"
	"testing"
)

func TestFormatShortSuggestion(t *testing.T) {
	var scenarioTable = []struct {
		in       []Suggest
		expected []Suggest
		max      int
		exWidth  int
	}{
		{
			in: []Suggest{
				{Text: "foo"},
				{Text: "bar"},
				{Text: "fuga"},
			},
			expected: []Suggest{
				{Text: " foo  "},
				{Text: " bar  "},
				{Text: " fuga "},
			},
			max:     100,
			exWidth: 6,
		},
		{
			in: []Suggest{
				{Text: "apple", Description: "This is apple."},
				{Text: "banana", Description: "This is banana."},
				{Text: "coconut", Description: "This is coconut."},
			},
			expected: []Suggest{
				{Text: " apple   ", Description: " This is apple.   "},
				{Text: " banana  ", Description: " This is banana.  "},
				{Text: " coconut ", Description: " This is coconut. "},
			},
			max:     100,
			exWidth: len(" apple   " + " This is apple.   "),
		},
		{
			in: []Suggest{
				{Text: "This is apple."},
				{Text: "This is banana."},
				{Text: "This is coconut."},
			},
			expected: []Suggest{
				{Text: " Thi... "},
				{Text: " Thi... "},
				{Text: " Thi... "},
			},
			max:     8,
			exWidth: 8,
		},
		{
			in: []Suggest{
				{Text: "This is apple."},
				{Text: "This is banana."},
				{Text: "This is coconut."},
			},
			expected: []Suggest{},
			max:      3,
			exWidth:  0,
		},
		{
			in: []Suggest{
				{Text: "--all-namespaces", Description: "-------------------------------------------------------------------------------------------------------------------------------------------"},
				{Text: "--allow-missing-template-keys", Description: "-----------------------------------------------------------------------------------------------------------------------------------------------"},
				{Text: "--export", Description: "----------------------------------------------------------------------------------------------------------"},
				{Text: "-f", Description: "-----------------------------------------------------------------------------------"},
				{Text: "--filename", Description: "-----------------------------------------------------------------------------------"},
				{Text: "--include-extended-apis", Description: "------------------------------------------------------------------------------------"},
			},
			expected: []Suggest{
				{Text: " --all-namespaces              ", Description: " --------------... "},
				{Text: " --allow-missing-template-keys ", Description: " --------------... "},
				{Text: " --export                      ", Description: " --------------... "},
				{Text: " -f                            ", Description: " --------------... "},
				{Text: " --filename                    ", Description: " --------------... "},
				{Text: " --include-extended-apis       ", Description: " --------------... "},
			},
			max:     50,
			exWidth: len(" --include-extended-apis       " + " ---------------..."),
		},
		{
			in: []Suggest{
				{Text: "--all-namespaces", Description: "If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace."},
				{Text: "--allow-missing-template-keys", Description: "If true, ignore any errors in templates when a field or map key is missing in the template. Only applies to golang and jsonpath output formats."},
				{Text: "--export", Description: "If true, use 'export' for the resources.  Exported resources are stripped of cluster-specific information."},
				{Text: "-f", Description: "Filename, directory, or URL to files identifying the resource to get from a server."},
				{Text: "--filename", Description: "Filename, directory, or URL to files identifying the resource to get from a server."},
				{Text: "--include-extended-apis", Description: "If true, include definitions of new APIs via calls to the API server. [default true]"},
			},
			expected: []Suggest{
				{Text: " --all-namespaces              ", Description: " If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.     "},
				{Text: " --allow-missing-template-keys ", Description: " If true, ignore any errors in templates when a field or map key is missing in the template. Only applies to golang and jsonpath output formats. "},
				{Text: " --export                      ", Description: " If true, use 'export' for the resources.  Exported resources are stripped of cluster-specific information.                                      "},
				{Text: " -f                            ", Description: " Filename, directory, or URL to files identifying the resource to get from a server.                                                             "},
				{Text: " --filename                    ", Description: " Filename, directory, or URL to files identifying the resource to get from a server.                                                             "},
				{Text: " --include-extended-apis       ", Description: " If true, include definitions of new APIs via calls to the API server. [default true]                                                            "},
			},
			max:     500,
			exWidth: len(" --include-extended-apis       " + " If true, include definitions of new APIs via calls to the API server. [default true]                                                            "),
		},
	}

	for i, s := range scenarioTable {
		actual, width := formatSuggestions(s.in, s.max)
		if width != s.exWidth {
			t.Errorf("[scenario %d] Want %d but got %d\n", i, s.exWidth, width)
		}
		if !reflect.DeepEqual(actual, s.expected) {
			t.Errorf("[scenario %d] Want %#v, but got %#v\n", i, s.expected, actual)
		}
	}
}

func TestFormatText(t *testing.T) {
	var scenarioTable = []struct {
		in       []string
		expected []string
		max      int
		exWidth  int
	}{
		{
			in: []string{
				"",
				"",
			},
			expected: []string{
				"",
				"",
			},
			max:     10,
			exWidth: 0,
		},
		{
			in: []string{
				"apple",
				"banana",
				"coconut",
			},
			expected: []string{
				"",
				"",
				"",
			},
			max:     2,
			exWidth: 0,
		},
		{
			in: []string{
				"apple",
				"banana",
				"coconut",
			},
			expected: []string{
				"",
				"",
				"",
			},
			max:     len(" " + " " + shortenSuffix),
			exWidth: 0,
		},
		{
			in: []string{
				"apple",
				"banana",
				"coconut",
			},
			expected: []string{
				" apple   ",
				" banana  ",
				" coconut ",
			},
			max:     100,
			exWidth: len(" coconut "),
		},
		{
			in: []string{
				"apple",
				"banana",
				"coconut",
			},
			expected: []string{
				" a... ",
				" b... ",
				" c... ",
			},
			max:     6,
			exWidth: 6,
		},
	}

	for i, s := range scenarioTable {
		actual, width := formatTexts(s.in, s.max, " ", " ")
		if width != s.exWidth {
			t.Errorf("[scenario %d] Want %d but got %d\n", i, s.exWidth, width)
		}
		if !reflect.DeepEqual(actual, s.expected) {
			t.Errorf("[scenario %d] Want %#v, but got %#v\n", i, s.expected, actual)
		}
	}
}
