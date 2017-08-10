package prompt

import (
	"testing"
	"reflect"
)

func TestFormatShortSuggestion(t *testing.T) {
	var scenarioTable = []struct{
		in       []Suggest
		expected []Suggest
		max      int
		exWidth  int
	} {
		{
			in: []Suggest{
				{Text:"foo"},
				{Text:"bar"},
				{Text:"fuga"},
			},
			expected: []Suggest{
				{Text:" foo  "},
				{Text:" bar  "},
				{Text:" fuga "},
			},
			max:     100,
			exWidth: 6,
		},
		{
			in: []Suggest{
				{Text:"apple", Description: "This is apple."},
				{Text:"banana", Description: "This is banana."},
				{Text:"coconut", Description: "This is coconut."},
			},
			expected: []Suggest{
				{Text:" apple   ", Description: " This is apple.   "},
				{Text:" banana  ", Description: " This is banana.  "},
				{Text:" coconut ", Description: " This is coconut. "},
			},
			max:     100,
			exWidth: len(" apple   " + " This is apple.   "),
		},
		{
			in: []Suggest{
				{Text:"--all-namespaces", Description:"If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace."},
				{Text:"--allow-missing-template-keys", Description:"If true, ignore any errors in templates when a field or map key is missing in the template. Only applies to golang and jsonpath output formats."},
				{Text:"--export", Description:"If true, use 'export' for the resources.  Exported resources are stripped of cluster-specific information."},
				{Text:"-f", Description:"Filename, directory, or URL to files identifying the resource to get from a server."},
				{Text:"--filename", Description:"Filename, directory, or URL to files identifying the resource to get from a server."},
				{Text:"--include-extended-apis", Description:"If true, include definitions of new APIs via calls to the API server. [default true]"},
			},
			expected: []Suggest{
				{Text:" --all-namespaces              ", Description:" If present, li... "},
				{Text:" --allow-missing-template-keys ", Description:" If true, ignor... "},
				{Text:" --export                      ", Description:" If true, use '... "},
				{Text:" -f                            ", Description:" Filename, dire... "},
				{Text:" --filename                    ", Description:" Filename, dire... "},
				{Text:" --include-extended-apis       ", Description:" If true, inclu... "},
			},
			max: 50,
			exWidth: len(" --include-extended-apis       " + " If true, includ..."),
		},
		{
			in: []Suggest{
				{Text:"--all-namespaces", Description:"If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace."},
				{Text:"--allow-missing-template-keys", Description:"If true, ignore any errors in templates when a field or map key is missing in the template. Only applies to golang and jsonpath output formats."},
				{Text:"--export", Description:"If true, use 'export' for the resources.  Exported resources are stripped of cluster-specific information."},
				{Text:"-f", Description:"Filename, directory, or URL to files identifying the resource to get from a server."},
				{Text:"--filename", Description:"Filename, directory, or URL to files identifying the resource to get from a server."},
				{Text:"--include-extended-apis", Description:"If true, include definitions of new APIs via calls to the API server. [default true]"},
			},
			expected: []Suggest{
				{Text:" --all-namespaces              ", Description:" If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.     "},
				{Text:" --allow-missing-template-keys ", Description:" If true, ignore any errors in templates when a field or map key is missing in the template. Only applies to golang and jsonpath output formats. "},
				{Text:" --export                      ", Description:" If true, use 'export' for the resources.  Exported resources are stripped of cluster-specific information.                                      "},
				{Text:" -f                            ", Description:" Filename, directory, or URL to files identifying the resource to get from a server.                                                             "},
				{Text:" --filename                    ", Description:" Filename, directory, or URL to files identifying the resource to get from a server.                                                             "},
				{Text:" --include-extended-apis       ", Description:" If true, include definitions of new APIs via calls to the API server. [default true]                                                            "},
			},
			max: 500,
			exWidth: len(" --include-extended-apis       " + " If true, include definitions of new APIs via calls to the API server. [default true]                                                            "),
		},
	}

	for i, s := range scenarioTable {
		actual, width := formatCompletions(s.in, s.max)
		if width != s.exWidth {
			t.Errorf("[scenario %d] Want %d but got %d\n", i, s.exWidth, width)
		}
		if !reflect.DeepEqual(actual, s.expected) {
			t.Errorf("[scenario %d] Want %#v, but got %#v\n", i, s.expected, actual)
		}
	}
}
