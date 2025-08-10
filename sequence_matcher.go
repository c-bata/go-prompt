package prompt

// Trie-based sequence matcher for efficient key sequence parsing.
// This is a Go port of replkit's SequenceMatcher.

import "sync"

// Use Key type from key.go

// TrieNode represents a node in the Trie structure.
type TrieNode struct {
	key      *Key
	children map[byte]*TrieNode
}

// SequenceMatcher provides efficient key sequence parsing.
type SequenceMatcher struct {
	root *TrieNode
	mu   sync.RWMutex
}

// MatchResult represents the result of matching a byte sequence.
type MatchResult int

const (
	NoMatch MatchResult = iota
	Exact
	Prefix
)

// LongestMatchResult holds the result of finding the longest valid sequence.
type LongestMatchResult struct {
	Key           Key
	ConsumedBytes int
}

// NewSequenceMatcher creates a new matcher and builds standard sequences.
func NewSequenceMatcher() *SequenceMatcher {
	matcher := &SequenceMatcher{
		root: &TrieNode{children: make(map[byte]*TrieNode)},
	}
	matcher.buildStandardSequences()
	return matcher
}

// Insert registers a custom sequence mapping.
func (sm *SequenceMatcher) Insert(bytes []byte, key Key) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	current := sm.root
	for _, b := range bytes {
		if current.children == nil {
			current.children = make(map[byte]*TrieNode)
		}
		if _, ok := current.children[b]; !ok {
			current.children[b] = &TrieNode{children: make(map[byte]*TrieNode)}
		}
		current = current.children[b]
	}
	current.key = &key
}

// MatchSequence returns whether the given bytes represent an exact match, prefix, or no match.
func (sm *SequenceMatcher) MatchSequence(bytes []byte) (MatchResult, *Key) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	if len(bytes) == 0 {
		return NoMatch, nil
	}
	current := sm.root
	for _, b := range bytes {
		child, ok := current.children[b]
		if !ok {
			return NoMatch, nil
		}
		current = child
	}
	if current.key != nil {
		return Exact, current.key
	}
	if len(current.children) > 0 {
		return Prefix, nil
	}
	return NoMatch, nil
}

// FindLongestMatch finds the longest valid sequence from the start of bytes.
func (sm *SequenceMatcher) FindLongestMatch(bytes []byte) *LongestMatchResult {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	var longest *LongestMatchResult
	current := sm.root
	for i, b := range bytes {
		child, ok := current.children[b]
		if !ok {
			break
		}
		current = child
		if current.key != nil {
			longest = &LongestMatchResult{Key: *current.key, ConsumedBytes: i + 1}
		}
	}
	return longest
}

// buildStandardSequences populates the Trie with standard key sequences.
func (sm *SequenceMatcher) buildStandardSequences() {
	// Control characters (single byte)
	sm.Insert([]byte{0x1b}, Escape)
	sm.Insert([]byte{0x00}, ControlSpace)
	sm.Insert([]byte{0x01}, ControlA)
	sm.Insert([]byte{0x02}, ControlB)
	sm.Insert([]byte{0x03}, ControlC)
	sm.Insert([]byte{0x04}, ControlD)
	sm.Insert([]byte{0x05}, ControlE)
	sm.Insert([]byte{0x06}, ControlF)
	sm.Insert([]byte{0x07}, ControlG)
	sm.Insert([]byte{0x08}, ControlH)
	sm.Insert([]byte{0x09}, Tab)
	sm.Insert([]byte{0x0a}, ControlJ)
	sm.Insert([]byte{0x0b}, ControlK)
	sm.Insert([]byte{0x0c}, ControlL)
	sm.Insert([]byte{0x0d}, Enter)
	sm.Insert([]byte{0x0e}, ControlN)
	sm.Insert([]byte{0x0f}, ControlO)
	sm.Insert([]byte{0x10}, ControlP)
	sm.Insert([]byte{0x11}, ControlQ)
	sm.Insert([]byte{0x12}, ControlR)
	sm.Insert([]byte{0x13}, ControlS)
	sm.Insert([]byte{0x14}, ControlT)
	sm.Insert([]byte{0x15}, ControlU)
	sm.Insert([]byte{0x16}, ControlV)
	sm.Insert([]byte{0x17}, ControlW)
	sm.Insert([]byte{0x18}, ControlX)
	sm.Insert([]byte{0x19}, ControlY)
	sm.Insert([]byte{0x1a}, ControlZ)
	sm.Insert([]byte{0x1c}, ControlBackslash)
	sm.Insert([]byte{0x1d}, ControlSquareClose)
	sm.Insert([]byte{0x1e}, ControlCircumflex)
	sm.Insert([]byte{0x1f}, ControlUnderscore)
	sm.Insert([]byte{0x7f}, Backspace)

	// Arrow keys (standard VT100)
	sm.Insert([]byte{0x1b, 0x5b, 0x41}, Up)
	sm.Insert([]byte{0x1b, 0x5b, 0x42}, Down)
	sm.Insert([]byte{0x1b, 0x5b, 0x43}, Right)
	sm.Insert([]byte{0x1b, 0x5b, 0x44}, Left)

	// Arrow keys (alternative sequences for some terminals)
	sm.Insert([]byte{0x1b, 0x4f, 0x41}, Up)
	sm.Insert([]byte{0x1b, 0x4f, 0x42}, Down)
	sm.Insert([]byte{0x1b, 0x4f, 0x43}, Right)
	sm.Insert([]byte{0x1b, 0x4f, 0x44}, Left)

	// Home and End keys (multiple variants)
	sm.Insert([]byte{0x1b, 0x5b, 0x48}, Home)
	sm.Insert([]byte{0x1b, 0x30, 0x48}, Home)
	sm.Insert([]byte{0x1b, 0x5b, 0x46}, End)
	sm.Insert([]byte{0x1b, 0x30, 0x46}, End)
	sm.Insert([]byte{0x1b, 0x5b, 0x31, 0x7e}, Home)
	sm.Insert([]byte{0x1b, 0x5b, 0x34, 0x7e}, End)
	sm.Insert([]byte{0x1b, 0x5b, 0x37, 0x7e}, Home)
	sm.Insert([]byte{0x1b, 0x5b, 0x38, 0x7e}, End)

	// Delete keys
	sm.Insert([]byte{0x1b, 0x5b, 0x33, 0x7e}, Delete)
	sm.Insert([]byte{0x1b, 0x5b, 0x33, 0x3b, 0x32, 0x7e}, ShiftDelete)
	sm.Insert([]byte{0x1b, 0x5b, 0x33, 0x3b, 0x35, 0x7e}, ControlDelete)

	// Page Up/Down
	sm.Insert([]byte{0x1b, 0x5b, 0x35, 0x7e}, PageUp)
	sm.Insert([]byte{0x1b, 0x5b, 0x36, 0x7e}, PageDown)

	// Insert and BackTab
	sm.Insert([]byte{0x1b, 0x5b, 0x32, 0x7e}, Insert)
	sm.Insert([]byte{0x1b, 0x5b, 0x5a}, BackTab)

	// Function keys F1-F4 (standard VT100)
	sm.Insert([]byte{0x1b, 0x4f, 0x50}, F1)
	sm.Insert([]byte{0x1b, 0x4f, 0x51}, F2)
	sm.Insert([]byte{0x1b, 0x4f, 0x52}, F3)
	sm.Insert([]byte{0x1b, 0x4f, 0x53}, F4)

	// Function keys F1-F5 (Linux console variants)
	sm.Insert([]byte{0x1b, 0x4f, 0x50, 0x41}, F1)
	sm.Insert([]byte{0x1b, 0x5b, 0x5b, 0x42}, F2)
	sm.Insert([]byte{0x1b, 0x5b, 0x5b, 0x43}, F3)
	sm.Insert([]byte{0x1b, 0x5b, 0x5b, 0x44}, F4)
	sm.Insert([]byte{0x1b, 0x5b, 0x5b, 0x45}, F5)

	// Function keys F1-F4 (rxvt-unicode variants)
	sm.Insert([]byte{0x1b, 0x5b, 0x31, 0x31, 0x7e}, F1)
	sm.Insert([]byte{0x1b, 0x5b, 0x31, 0x32, 0x7e}, F2)
	sm.Insert([]byte{0x1b, 0x5b, 0x31, 0x33, 0x7e}, F3)
	sm.Insert([]byte{0x1b, 0x5b, 0x31, 0x34, 0x7e}, F4)

	// Function keys F5-F12
	sm.Insert([]byte{0x1b, 0x5b, 0x31, 0x35, 0x7e}, F5)
	sm.Insert([]byte{0x1b, 0x5b, 0x31, 0x37, 0x7e}, F6)
	sm.Insert([]byte{0x1b, 0x5b, 0x31, 0x38, 0x7e}, F7)
	sm.Insert([]byte{0x1b, 0x5b, 0x31, 0x39, 0x7e}, F8)
	sm.Insert([]byte{0x1b, 0x5b, 0x32, 0x30, 0x7e}, F9)
	sm.Insert([]byte{0x1b, 0x5b, 0x32, 0x31, 0x7e}, F10)
	sm.Insert([]byte{0x1b, 0x5b, 0x32, 0x33, 0x7e}, F11)
	sm.Insert([]byte{0x1b, 0x5b, 0x32, 0x34, 0x7e}, F12)

	// Function keys F13-F20 (basic sequences)
	sm.Insert([]byte{0x1b, 0x5b, 0x32, 0x35, 0x7e}, F13)
	sm.Insert([]byte{0x1b, 0x5b, 0x32, 0x36, 0x7e}, F14)
	sm.Insert([]byte{0x1b, 0x5b, 0x32, 0x38, 0x7e}, F15)
	sm.Insert([]byte{0x1b, 0x5b, 0x32, 0x39, 0x7e}, F16)
	sm.Insert([]byte{0x1b, 0x5b, 0x33, 0x31, 0x7e}, F17)
	sm.Insert([]byte{0x1b, 0x5b, 0x33, 0x32, 0x7e}, F18)
	sm.Insert([]byte{0x1b, 0x5b, 0x33, 0x33, 0x7e}, F19)
	sm.Insert([]byte{0x1b, 0x5b, 0x33, 0x34, 0x7e}, F20)

	// Function keys F13-F24 (Xterm variants)
	sm.Insert([]byte{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x50}, F13)
	sm.Insert([]byte{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x51}, F14)
	sm.Insert([]byte{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x52}, F16)
	sm.Insert([]byte{0x1b, 0x5b, 0x31, 0x35, 0x3b, 0x32, 0x7e}, F17)
	sm.Insert([]byte{0x1b, 0x5b, 0x31, 0x37, 0x3b, 0x32, 0x7e}, F18)
	sm.Insert([]byte{0x1b, 0x5b, 0x31, 0x38, 0x3b, 0x32, 0x7e}, F19)
	sm.Insert([]byte{0x1b, 0x5b, 0x31, 0x39, 0x3b, 0x32, 0x7e}, F20)
	sm.Insert([]byte{0x1b, 0x5b, 0x32, 0x30, 0x3b, 0x32, 0x7e}, F21)
	sm.Insert([]byte{0x1b, 0x5b, 0x32, 0x31, 0x3b, 0x32, 0x7e}, F22)
	sm.Insert([]byte{0x1b, 0x5b, 0x32, 0x33, 0x3b, 0x32, 0x7e}, F23)
	sm.Insert([]byte{0x1b, 0x5b, 0x32, 0x34, 0x3b, 0x32, 0x7e}, F24)

	// Control + Arrow keys
	sm.Insert([]byte{0x1b, 0x5b, 0x31, 0x3b, 0x35, 0x41}, ControlUp)
	sm.Insert([]byte{0x1b, 0x5b, 0x31, 0x3b, 0x35, 0x42}, ControlDown)
	sm.Insert([]byte{0x1b, 0x5b, 0x31, 0x3b, 0x35, 0x43}, ControlRight)
	sm.Insert([]byte{0x1b, 0x5b, 0x31, 0x3b, 0x35, 0x44}, ControlLeft)

	// Alternative Control + Arrow keys
	sm.Insert([]byte{0x1b, 0x5b, 0x35, 0x41}, ControlUp)
	sm.Insert([]byte{0x1b, 0x5b, 0x35, 0x42}, ControlDown)
	sm.Insert([]byte{0x1b, 0x5b, 0x35, 0x43}, ControlRight)
	sm.Insert([]byte{0x1b, 0x5b, 0x35, 0x44}, ControlLeft)

	// rxvt Control + Arrow keys
	sm.Insert([]byte{0x1b, 0x5b, 0x4f, 0x63}, ControlRight)
	sm.Insert([]byte{0x1b, 0x5b, 0x4f, 0x64}, ControlLeft)

	// Shift + Arrow keys
	sm.Insert([]byte{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x41}, ShiftUp)
	sm.Insert([]byte{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x42}, ShiftDown)
	sm.Insert([]byte{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x43}, ShiftRight)
	sm.Insert([]byte{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x44}, ShiftLeft)

	// Ignore sequences (terminal-specific sequences that should be ignored)
	sm.Insert([]byte{0x1b, 0x5b, 0x45}, Ignore) // Xterm
	sm.Insert([]byte{0x1b, 0x5b, 0x46}, Ignore) // Linux console
}
