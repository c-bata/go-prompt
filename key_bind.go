package prompt

type KeyBindFunc func(*Buffer) *Buffer

type KeyBind struct {
	Key Key
	Fn  KeyBindFunc
}

var commonKeyBindings = []KeyBind {
	// Go to the End of the line
	{
		Key: End,
		Fn: func(buf *Buffer) *Buffer {
			x := []rune(buf.Document().TextAfterCursor())
			buf.CursorRight(len(x))
			return buf
		},
	},
	// Go to the beginning of the line
	{
		Key: Home,
		Fn: func(buf *Buffer) *Buffer {
			x := []rune(buf.Document().TextBeforeCursor())
			buf.CursorLeft(len(x))
			return buf
		},
	},
	// Delete character under the cursor
	{
		Key: Delete,
		Fn: func(buf *Buffer) *Buffer {
			buf.Delete(1)
			return buf
		},
	},
	// Backspace
	{
		Key: Backspace,
		Fn: func(buf *Buffer) *Buffer {
			buf.DeleteBeforeCursor(1)
			return buf
		},
	},
	// Right allow: Forward one character
	{
		Key: Right,
		Fn: func(buf *Buffer) *Buffer {
			buf.CursorRight(1)
			return buf
		},
	},
	// Left allow: Backward one character
	{
		Key: Left,
		Fn: func(buf *Buffer) *Buffer {
			buf.CursorLeft(1)
			return buf
		},
	},

}
