package prompt

import (
	"testing"
)

func TestKeyParser_Feed(t *testing.T) {
	parser := NewKeyParser()

	// Test ControlC
	input := []byte{0x03}
	events := parser.Feed(input)
	if len(events) != 1 || events[0].Key != ControlC {
		t.Errorf("Expected ControlC, got %+v", events)
	}

	// Test Up Arrow
	input = []byte{0x1b, 0x5b, 0x41}
	events = parser.Feed(input)
	if len(events) != 1 || events[0].Key != Up {
		t.Errorf("Expected Up, got %+v", events)
	}

	// Test unknown sequence (should fallback to text)
	input = []byte{0xff}
	events = parser.Feed(input)
	if len(events) != 1 || events[0].Key != NotDefined || events[0].Text != string([]byte{0xff}) {
		t.Errorf("Expected NotDefined with text, got %+v", events)
	}

	// Test multiple keys in one input
	input = []byte{0x03, 0x1b, 0x5b, 0x41}
	events = parser.Feed(input)
	if len(events) != 2 || events[0].Key != ControlC || events[1].Key != Up {
		t.Errorf("Expected ControlC and Up, got %+v", events)
	}

	// Test custom sequence
	parser.matcher.Insert([]byte("gg"), F24)
	input = []byte("gg")
	events = parser.Feed(input)
	if len(events) != 1 || events[0].Key != F24 {
		t.Errorf("Expected F24 for 'gg', got %+v", events)
	}
}
