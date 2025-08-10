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
		// Bracketed Paste: ESC[200~ ... ESC[201~
		if len(kp.buffer) >= 6 && string(kp.buffer[:6]) == "\x1b[200~" {
			// Start of bracketed paste
			endIdx := -1
			for i := 6; i+5 <= len(kp.buffer); i++ {
				if string(kp.buffer[i:i+6]) == "\x1b[201~" {
					endIdx = i + 6
					break
				}
			}
			if endIdx != -1 {
				events = append(events, KeyEvent{Key: BracketedPaste, RawBytes: kp.buffer[:endIdx]})
				kp.buffer = kp.buffer[endIdx:]
				continue
			} else {
				// Wait for more bytes
				break
			}
		}

		// Mouse event: ESC[M or ESC[<
		if len(kp.buffer) >= 3 && kp.buffer[0] == 0x1b && kp.buffer[1] == 0x5b && (kp.buffer[2] == 'M' || kp.buffer[2] == '<') {
			// For simplicity, treat as mouse event and consume 6 bytes (typical length)
			n := 6
			if len(kp.buffer) >= n {
				events = append(events, KeyEvent{Key: Vt100MouseEvent, RawBytes: kp.buffer[:n]})
				kp.buffer = kp.buffer[n:]
				continue
			} else {
				break
			}
		}

		// CPR: ESC[{row};{col}R
		if len(kp.buffer) >= 2 && kp.buffer[0] == 0x1b && kp.buffer[1] == 0x5b {
			found := false
			for i := 2; i < len(kp.buffer); i++ {
				if kp.buffer[i] == 'R' {
					events = append(events, KeyEvent{Key: CPRResponse, RawBytes: kp.buffer[:i+1]})
					kp.buffer = kp.buffer[i+1:]
					found = true
					break
				}
			}
			if found {
				continue
			}
		}

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
