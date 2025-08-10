package prompt

import (
	"testing"
)

func TestKeyParser_BracketedPaste(t *testing.T) {
	parser := NewKeyParser()
	input := []byte{0x1b, 0x5b, '2', '0', '0', '~', 'h', 'e', 'l', 'l', 'o', 0x1b, 0x5b, '2', '0', '1', '~'}
	events := parser.Feed(input)
	if len(events) != 1 || events[0].Key != BracketedPaste {
		t.Errorf("Expected BracketedPaste, got %+v", events)
	}
}

func TestKeyParser_MouseEvent(t *testing.T) {
	parser := NewKeyParser()
	input := []byte{0x1b, 0x5b, 'M', 0x20, 0x21, 0x22}
	events := parser.Feed(input)
	if len(events) != 1 || events[0].Key != Vt100MouseEvent {
		t.Errorf("Expected Vt100MouseEvent, got %+v", events)
	}
}

func TestKeyParser_CPRResponse(t *testing.T) {
	parser := NewKeyParser()
	input := []byte{0x1b, 0x5b, '1', ';', '2', '4', 'R'}
	events := parser.Feed(input)
	if len(events) != 1 || events[0].Key != CPRResponse {
		t.Errorf("Expected CPRResponse, got %+v", events)
	}
}
