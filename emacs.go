package prompt

var emacsKeyBindings = []KeyBind {
	// Go to the End of the line
	{
		Key: ControlE,
		Fn: func(buf *Buffer) *Buffer {
			x := []rune(buf.Document().TextAfterCursor())
			buf.CursorRight(len(x))
			return buf
		},
	},
	// Go to the beginning of the line
	{
		Key: ControlA,
		Fn: func(buf *Buffer) *Buffer {
			x := []rune(buf.Document().TextBeforeCursor())
			buf.CursorLeft(len(x))
			return buf
		},
	},
	// Cut the Line after the cursor
	{
		Key: ControlK,
		Fn: func(buf *Buffer) *Buffer {
			x := []rune(buf.Document().TextAfterCursor())
			buf.Delete(len(x))
			return buf
		},
	},
	// Cut/delete the Line before the cursor
	{
		Key: ControlU,
		Fn: func(buf *Buffer) *Buffer {
			x := []rune(buf.Document().TextBeforeCursor())
			buf.DeleteBeforeCursor(len(x))
			return buf
		},
	},
	// Delete character under the cursor
	{
		Key: ControlD,
		Fn: func(buf *Buffer) *Buffer {
			if buf.Text() == "" {
				return buf  // This means just exit.
			}
			buf.Delete(1)
			return buf
		},
	},
	// Backspace
	{
		Key: ControlH,
		Fn: func(buf *Buffer) *Buffer {
			buf.DeleteBeforeCursor(1)
			return buf
		},
	},
	// Right allow: Forward one character
	{
		Key: ControlF,
		Fn: func(buf *Buffer) *Buffer {
			buf.CursorRight(1)
			return buf
		},
	},
	// Left allow: Backward one character
	{
		Key: ControlB,
		Fn: func(buf *Buffer) *Buffer {
			buf.CursorLeft(1)
			return buf
		},
	},
	// Cut the Word before the cursor.
	{
		Key: ControlW,
		Fn: func(buf *Buffer) *Buffer {
			buf.DeleteBeforeCursor(len([]rune(buf.Document().GetWordBeforeCursor())))
			return buf
		},
	},
}
