package prompt

// KeyParser implements a state machine for key sequence parsing using SequenceMatcher.
type ParserState int

const (
	Normal ParserState = iota
	EscapeSequence
	CsiSequence
	MouseEvent
)

type KeyParser struct {
	state   ParserState
	buffer  []byte
	matcher *SequenceMatcher
}

// NewKeyParser creates a new KeyParser with a SequenceMatcher.
func NewKeyParser() *KeyParser {
	return &KeyParser{
		state:   Normal,
		buffer:  make([]byte, 0, 32),
		matcher: NewSequenceMatcher(),
	}
}

// Feed processes input bytes and returns parsed key events.
func (kp *KeyParser) Feed(input []byte) []KeyEvent {
	var events []KeyEvent
	kp.buffer = append(kp.buffer, input...)
	for len(kp.buffer) > 0 {
		match := kp.matcher.FindLongestMatch(kp.buffer)
		if match != nil {
			events = append(events, KeyEvent{Key: match.Key, RawBytes: kp.buffer[:match.ConsumedBytes]})
			kp.buffer = kp.buffer[match.ConsumedBytes:]
		} else {
			// Not a known sequence, treat first byte as text
			events = append(events, KeyEvent{Key: NotDefined, RawBytes: kp.buffer[:1], Text: string(kp.buffer[:1])})
			kp.buffer = kp.buffer[1:]
		}
	}
	return events
}

type KeyEvent struct {
	Key      Key
	RawBytes []byte
	Text     string // For unicode/unknown input
}
