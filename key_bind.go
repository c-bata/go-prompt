package prompt

// KeyBindFunc receives buffer and processed it.
type KeyBindFunc func(*Buffer)

// KeyBind represents which key should do what operation.
type KeyBind struct {
	Key Key
	Fn  KeyBindFunc
}

// KeyBindMode to switch a key binding flexibly.
type KeyBindMode string

const (
	// CommonKeyBind is a mode without any keyboard shortcut
	CommonKeyBind KeyBindMode = "common"
	// EmacsKeyBind is a mode to use emacs-like keyboard shortcut
	EmacsKeyBind KeyBindMode = "emacs"
)

var commonKeyBindings = []KeyBind{
	// Go to the End of the line
	{
		Key: End,
		Fn: func(buf *Buffer) {
			x := []rune(buf.Document().TextAfterCursor())
			buf.CursorRight(len(x))
		},
	},
	// Go to the beginning of the line
	{
		Key: Home,
		Fn: func(buf *Buffer) {
			x := []rune(buf.Document().TextBeforeCursor())
			buf.CursorLeft(len(x))
		},
	},
	// Delete character under the cursor
	{
		Key: Delete,
		Fn: func(buf *Buffer) {
			buf.Delete(1)
		},
	},
	// Backspace
	{
		Key: Backspace,
		Fn: func(buf *Buffer) {
			buf.DeleteBeforeCursor(1)
		},
	},
	// Right allow: Forward one character
	{
		Key: Right,
		Fn: func(buf *Buffer) {
			buf.CursorRight(1)
		},
	},
	// Left allow: Backward one character
	{
		Key: Left,
		Fn: func(buf *Buffer) {
			buf.CursorLeft(1)
		},
	},
}
