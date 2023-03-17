package prompt

import (
	"github.com/stretchr/testify/assert"
	"pgregory.net/rapid"
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
			assert.Equal(t, s.expected, key)
		})
	}
}

func RandomASCIIByteSequence() *rapid.Generator[[]byte] {
	return rapid.Custom(func(t *rapid.T) []byte {
		return rapid.SampledFrom(ASCIISequences).Draw(t, "random ascii sequence").ASCIICode
	})
}

func TestSanitizeInputWithASCIISequences(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		expectedString := []byte("this_is_a_longer_sized_text_input_for_testing_purposes")
		inputString := make([]byte, len(expectedString))
		//at each index insert some random number of ascii control sequences
		for _, char := range expectedString {
			inputString = append(inputString, char)
			//append 1-5 ascii control sequences
			sequences := rapid.SliceOfN(RandomASCIIByteSequence(), 1, 5).Draw(t, "random number of ascii control sequences")
			for _, sequence := range sequences {
				inputString = append(inputString, sequence...)
			}
		}
		assert.Equal(t, expectedString, RemoveASCIISequences(inputString))
	})
}
