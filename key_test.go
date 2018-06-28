package prompt

import (
	"testing"
)

func TestGetKey(t *testing.T) {
	scenarioTable := []struct {
		input    []byte
		expected Key
	}{
		{
			input:    []byte{0x1b},
			expected: Escape,
		},
		{
			input:    []byte{'a'},
			expected: NotDefined,
		},
	}

	for _, s := range scenarioTable {
		key := GetKey(s.input)
		if key != s.expected {
			t.Errorf("Should be %s, but got %s", key, s.expected)
		}
	}
}
