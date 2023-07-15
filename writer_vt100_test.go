package prompt

import (
	"bytes"
	"testing"
)

func TestVT100WriterWrite(t *testing.T) {
	scenarioTable := []struct {
		input    []byte
		expected []byte
	}{
		{
			input:    []byte{0x1b},
			expected: []byte{'?'},
		},
		{
			input:    []byte{'a'},
			expected: []byte{'a'},
		},
	}

	for _, s := range scenarioTable {
		pw := &VT100Writer{}
		if _, err := pw.Write(s.input); err != nil {
			panic(err)
		}

		if !bytes.Equal(pw.buffer, s.expected) {
			t.Errorf("Should be %+#v, but got %+#v", pw.buffer, s.expected)
		}
	}
}

func TestVT100WriterWriteString(t *testing.T) {
	scenarioTable := []struct {
		input    string
		expected []byte
	}{
		{
			input:    "\x1b",
			expected: []byte{'?'},
		},
		{
			input:    "a",
			expected: []byte{'a'},
		},
	}

	for _, s := range scenarioTable {
		pw := &VT100Writer{}
		if _, err := pw.WriteString(s.input); err != nil {
			panic(err)
		}

		if !bytes.Equal(pw.buffer, s.expected) {
			t.Errorf("Should be %+#v, but got %+#v", pw.buffer, s.expected)
		}
	}
}

func TestVT100WriterWriteRawString(t *testing.T) {
	scenarioTable := []struct {
		input    string
		expected []byte
	}{
		{
			input:    "\x1b",
			expected: []byte{0x1b},
		},
		{
			input:    "a",
			expected: []byte{'a'},
		},
	}

	for _, s := range scenarioTable {
		pw := &VT100Writer{}
		pw.WriteRawString(s.input)

		if !bytes.Equal(pw.buffer, s.expected) {
			t.Errorf("Should be %+#v, but got %+#v", pw.buffer, s.expected)
		}
	}
}
