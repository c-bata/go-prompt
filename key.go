// Code generated "This is a fake comment to avoid golint errors"; DO NOT EDIT.
// FIXME: This is a little bit stupid, but there are many public constants which is no value for writing godoc comment.

package prompt

import "bytes"

// Key is the type express the key inserted from user.
type Key int

// ASCIICode is the type contains Key and it's ascii byte array.
type ASCIICode struct {
	Key       Key
	ASCIICode []byte
}

const (
	Escape Key = iota

	ControlA
	ControlB
	ControlC
	ControlD
	ControlE
	ControlF
	ControlG
	ControlH
	ControlI
	ControlJ
	ControlK
	ControlL
	ControlM
	ControlN
	ControlO
	ControlP
	ControlQ
	ControlR
	ControlS
	ControlT
	ControlU
	ControlV
	ControlW
	ControlX
	ControlY
	ControlZ

	ControlSpace
	ControlBackslash
	ControlSquareClose
	ControlCircumflex
	ControlUnderscore
	ControlLeft
	ControlRight
	ControlUp
	ControlDown

	Up
	Down
	Right
	Left

	ShiftLeft
	ShiftUp
	ShiftDown
	ShiftRight

	Home
	End
	Delete
	ShiftDelete
	ControlDelete
	PageUp
	PageDown
	BackTab
	Insert
	Backspace

	// Aliases.
	Tab
	Enter
	// Actually Enter equals ControlM, not ControlJ,
	// However, in prompt_toolkit, we made the mistake of translating
	// \r into \n during the input, so everyone is now handling the
	// enter key by binding ControlJ.

	// From now on, it's better to bind `ASCII_SEQUENCES.Enter` everywhere,
	// because that's future compatible, and will still work when we
	// stop replacing \r by \n.

	F1
	F2
	F3
	F4
	F5
	F6
	F7
	F8
	F9
	F10
	F11
	F12
	F13
	F14
	F15
	F16
	F17
	F18
	F19
	F20
	F21
	F22
	F23
	F24

	// Matches any key.
	Any

	// Special
	CPRResponse
	Vt100MouseEvent
	WindowsMouseEvent
	BracketedPaste

	// Key which is ignored. (The key binding for this key should not do anything.)
	Ignore

	// Key is not defined
	NotDefined
)

// GetKey returns Key correspond to input byte codes.
func GetKey(b []byte) Key {
	for _, k := range asciiSequences {
		if bytes.Equal(k.ASCIICode, b) {
			return k.Key
		}
	}
	return NotDefined
}

var asciiSequences = []*ASCIICode{
	{Key: Escape, ASCIICode: []byte{0x1b}},

	{Key: ControlSpace, ASCIICode: []byte{0x00}},
	{Key: ControlA, ASCIICode: []byte{0x1}},
	{Key: ControlB, ASCIICode: []byte{0x2}},
	{Key: ControlC, ASCIICode: []byte{0x3}},
	{Key: ControlD, ASCIICode: []byte{0x4}},
	{Key: ControlE, ASCIICode: []byte{0x5}},
	{Key: ControlF, ASCIICode: []byte{0x6}},
	{Key: ControlG, ASCIICode: []byte{0x7}},
	{Key: ControlH, ASCIICode: []byte{0x8}},
	//{Key: ControlI, ASCIICode: []byte{0x9}},
	//{Key: ControlJ, ASCIICode: []byte{0xa}},
	{Key: ControlK, ASCIICode: []byte{0xb}},
	{Key: ControlL, ASCIICode: []byte{0xc}},
	{Key: ControlM, ASCIICode: []byte{0xd}},
	{Key: ControlN, ASCIICode: []byte{0xe}},
	{Key: ControlO, ASCIICode: []byte{0xf}},
	{Key: ControlP, ASCIICode: []byte{0x10}},
	{Key: ControlQ, ASCIICode: []byte{0x11}},
	{Key: ControlR, ASCIICode: []byte{0x12}},
	{Key: ControlS, ASCIICode: []byte{0x13}},
	{Key: ControlT, ASCIICode: []byte{0x14}},
	{Key: ControlU, ASCIICode: []byte{0x15}},
	{Key: ControlV, ASCIICode: []byte{0x16}},
	{Key: ControlW, ASCIICode: []byte{0x17}},
	{Key: ControlX, ASCIICode: []byte{0x18}},
	{Key: ControlY, ASCIICode: []byte{0x19}},
	{Key: ControlZ, ASCIICode: []byte{0x1a}},

	{Key: ControlBackslash, ASCIICode: []byte{0x1c}},
	{Key: ControlSquareClose, ASCIICode: []byte{0x1d}},
	{Key: ControlCircumflex, ASCIICode: []byte{0x1e}},
	{Key: ControlUnderscore, ASCIICode: []byte{0x1f}},
	{Key: Backspace, ASCIICode: []byte{0x7f}},

	{Key: Up, ASCIICode: []byte{0x1b, 0x5b, 0x41}},
	{Key: Down, ASCIICode: []byte{0x1b, 0x5b, 0x42}},
	{Key: Right, ASCIICode: []byte{0x1b, 0x5b, 0x43}},
	{Key: Left, ASCIICode: []byte{0x1b, 0x5b, 0x44}},
	{Key: Home, ASCIICode: []byte{0x1b, 0x5b, 0x48}},
	{Key: Home, ASCIICode: []byte{0x1b, 0x30, 0x48}},
	{Key: End, ASCIICode: []byte{0x1b, 0x5b, 0x46}},
	{Key: End, ASCIICode: []byte{0x1b, 0x30, 0x46}},

	{Key: Enter, ASCIICode: []byte{0xa}},
	{Key: Delete, ASCIICode: []byte{0x1b, 0x5b, 0x33, 0x7e}},
	{Key: ShiftDelete, ASCIICode: []byte{0x1b, 0x5b, 0x33, 0x3b, 0x32, 0x7e}},
	{Key: ControlDelete, ASCIICode: []byte{0x1b, 0x5b, 0x33, 0x3b, 0x35, 0x7e}},
	{Key: Home, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x7e}},
	{Key: End, ASCIICode: []byte{0x1b, 0x5b, 0x34, 0x7e}},
	{Key: PageUp, ASCIICode: []byte{0x1b, 0x5b, 0x35, 0x7e}},
	{Key: PageDown, ASCIICode: []byte{0x1b, 0x5b, 0x36, 0x7e}},
	{Key: Home, ASCIICode: []byte{0x1b, 0x5b, 0x37, 0x7e}},
	{Key: End, ASCIICode: []byte{0x1b, 0x5b, 0x38, 0x7e}},
	{Key: Tab, ASCIICode: []byte{0x9}},
	{Key: BackTab, ASCIICode: []byte{0x1b, 0x5b, 0x5a}},
	{Key: Insert, ASCIICode: []byte{0x1b, 0x5b, 0x32, 0x7e}},

	{Key: F1, ASCIICode: []byte{0x1b, 0x4f, 0x50}},
	{Key: F2, ASCIICode: []byte{0x1b, 0x4f, 0x51}},
	{Key: F3, ASCIICode: []byte{0x1b, 0x4f, 0x52}},
	{Key: F4, ASCIICode: []byte{0x1b, 0x4f, 0x53}},

	{Key: F1, ASCIICode: []byte{0x1b, 0x4f, 0x50, 0x41}}, // Linux console
	{Key: F2, ASCIICode: []byte{0x1b, 0x5b, 0x5b, 0x42}}, // Linux console
	{Key: F3, ASCIICode: []byte{0x1b, 0x5b, 0x5b, 0x43}}, // Linux console
	{Key: F4, ASCIICode: []byte{0x1b, 0x5b, 0x5b, 0x44}}, // Linux console
	{Key: F5, ASCIICode: []byte{0x1b, 0x5b, 0x5b, 0x45}}, // Linux console

	{Key: F1, ASCIICode: []byte{0x1b, 0x5b, 0x11, 0x7e}}, // rxvt-unicode
	{Key: F2, ASCIICode: []byte{0x1b, 0x5b, 0x12, 0x7e}}, // rxvt-unicode
	{Key: F3, ASCIICode: []byte{0x1b, 0x5b, 0x13, 0x7e}}, // rxvt-unicode
	{Key: F4, ASCIICode: []byte{0x1b, 0x5b, 0x14, 0x7e}}, // rxvt-unicode

	{Key: F5, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x35, 0x7e}},
	{Key: F6, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x37, 0x7e}},
	{Key: F7, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x38, 0x7e}},
	{Key: F8, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x39, 0x7e}},
	{Key: F9, ASCIICode: []byte{0x1b, 0x5b, 0x32, 0x30, 0x7e}},
	{Key: F10, ASCIICode: []byte{0x1b, 0x5b, 0x32, 0x31, 0x7e}},
	{Key: F11, ASCIICode: []byte{0x1b, 0x5b, 0x32, 0x32, 0x7e}},
	{Key: F12, ASCIICode: []byte{0x1b, 0x5b, 0x32, 0x34, 0x7e, 0x8}},
	{Key: F13, ASCIICode: []byte{0x1b, 0x5b, 0x25, 0x7e}},
	{Key: F14, ASCIICode: []byte{0x1b, 0x5b, 0x26, 0x7e}},
	{Key: F15, ASCIICode: []byte{0x1b, 0x5b, 0x28, 0x7e}},
	{Key: F16, ASCIICode: []byte{0x1b, 0x5b, 0x29, 0x7e}},
	{Key: F17, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x7e}},
	{Key: F18, ASCIICode: []byte{0x1b, 0x5b, 0x32, 0x7e}},
	{Key: F19, ASCIICode: []byte{0x1b, 0x5b, 0x33, 0x7e}},
	{Key: F20, ASCIICode: []byte{0x1b, 0x5b, 0x34, 0x7e}},

	// Xterm
	{Key: F13, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x50}},
	{Key: F14, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x51}},
	// &ASCIICode{Key: F15, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x52}},  // Conflicts with CPR response
	{Key: F16, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x52}},
	{Key: F17, ASCIICode: []byte{0x1b, 0x5b, 0x15, 0x3b, 0x32, 0x7e}},
	{Key: F18, ASCIICode: []byte{0x1b, 0x5b, 0x17, 0x3b, 0x32, 0x7e}},
	{Key: F19, ASCIICode: []byte{0x1b, 0x5b, 0x18, 0x3b, 0x32, 0x7e}},
	{Key: F20, ASCIICode: []byte{0x1b, 0x5b, 0x19, 0x3b, 0x32, 0x7e}},
	{Key: F21, ASCIICode: []byte{0x1b, 0x5b, 0x20, 0x3b, 0x32, 0x7e}},
	{Key: F22, ASCIICode: []byte{0x1b, 0x5b, 0x21, 0x3b, 0x32, 0x7e}},
	{Key: F23, ASCIICode: []byte{0x1b, 0x5b, 0x23, 0x3b, 0x32, 0x7e}},
	{Key: F24, ASCIICode: []byte{0x1b, 0x5b, 0x24, 0x3b, 0x32, 0x7e}},

	{Key: ControlUp, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x35, 0x41}},
	{Key: ControlDown, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x35, 0x42}},
	{Key: ControlRight, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x35, 0x43}},
	{Key: ControlLeft, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x35, 0x44}},

	{Key: ShiftUp, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x41}},
	{Key: ShiftDown, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x42}},
	{Key: ShiftRight, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x43}},
	{Key: ShiftLeft, ASCIICode: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x44}},

	// Tmux sends following keystrokes when control+arrow is pressed, but for
	// Emacs ansi-term sends the same sequences for normal arrow keys. Consider
	// it a normal arrow press, because that's more important.
	{Key: Up, ASCIICode: []byte{0x1b, 0x4f, 0x41}},
	{Key: Down, ASCIICode: []byte{0x1b, 0x4f, 0x42}},
	{Key: Right, ASCIICode: []byte{0x1b, 0x4f, 0x43}},
	{Key: Left, ASCIICode: []byte{0x1b, 0x4f, 0x44}},

	{Key: ControlUp, ASCIICode: []byte{0x1b, 0x5b, 0x35, 0x41}},
	{Key: ControlDown, ASCIICode: []byte{0x1b, 0x5b, 0x35, 0x42}},
	{Key: ControlRight, ASCIICode: []byte{0x1b, 0x5b, 0x35, 0x43}},
	{Key: ControlLeft, ASCIICode: []byte{0x1b, 0x5b, 0x35, 0x44}},

	{Key: ControlRight, ASCIICode: []byte{0x1b, 0x5b, 0x4f, 0x63}}, // rxvt
	{Key: ControlLeft, ASCIICode: []byte{0x1b, 0x5b, 0x4f, 0x64}},  // rxvt

	{Key: Ignore, ASCIICode: []byte{0x1b, 0x5b, 0x45}}, // Xterm
	{Key: Ignore, ASCIICode: []byte{0x1b, 0x5b, 0x46}}, // Linux console
}
