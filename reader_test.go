package prompt

import (
	"testing"
)

func TestPosixParserGetKey(t *testing.T) {
	scenarioTable := []struct {
		name     string
		input    []byte
		expected Key
	}{
		{
			name:     "escape",
			input:    []byte{0x1b},
			expected: Escape,
		},
		{
			name:     "undefined",
			input:    []byte{'a'},
			expected: NotDefined,
		},
	}

	for _, s := range scenarioTable {
		t.Run(s.name, func(t *testing.T) {
			key := GetKey(s.input)
			if key != s.expected {
				t.Errorf("Should be %s, but got %s", key, s.expected)
			}
		})
	}
}
