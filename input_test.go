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

func RandomASCIIByteSequence() *rapid.Generator[*ASCIICode] {
	return rapid.Custom(func(t *rapid.T) *ASCIICode {
		return rapid.SampledFrom(ASCIISequences).Draw(t, "random ascii sequence")
	})
}

func TestSanitizeInputWithASCIISequences(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		testString := []byte("this_is_a_longer_sized_text_input_for_testing_purposes")
		inputString := make([]byte, len(testString))
		expectedString := make([]byte, 0)
		//at each index insert some random number of ascii control sequences
		for _, char := range testString {
			inputString = append(inputString, char)
			expectedString = append(expectedString, char)
			//append 1-5 ascii control sequences
			sequences := rapid.SliceOfN(RandomASCIIByteSequence(), 1, 5).Draw(t, "random number of ascii control sequences")
			for _, sequence := range sequences {
				// allow \n and \r to be inserted
				if sequence.Key == Enter || sequence.Key == ControlM {
					expectedString = append(expectedString, sequence.ASCIICode...)
				}

				inputString = append(inputString, sequence.ASCIICode...)
			}
		}
		assert.Equal(t, string(expectedString), string(RemoveASCIISequences(inputString)))
	})
}
