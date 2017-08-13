package prompt

/*

========
PROGRESS
========

Moving the cursor
-----------------

* [x] Ctrl + a   Go to the beginning of the line (Home)
* [x] Ctrl + e   Go to the End of the line (End)
* [x] Ctrl + p   Previous command (Up arrow)
* [x] Ctrl + n   Next command (Down arrow)
* [ ] Alt + b   Back (left) one word
* [ ] Alt + f   Forward (right) one word
* [x] Ctrl + f   Forward one character
* [x] Ctrl + b   Backward one character
* [x] Ctrl + xx  Toggle between the start of line and current cursor position

Editing
-------

* [ ] Ctrl + L   Clear the Screen, similar to the clear command
* [ ] Alt + Del Delete the Word before the cursor.
* [ ] Alt + d   Delete the Word after the cursor.
* [x] Ctrl + d   Delete character under the cursor
* [x] Ctrl + h   Delete character before the cursor (Backspace)

* [x] Ctrl + w   Cut the Word before the cursor to the clipboard.
* [x] Ctrl + k   Cut the Line after the cursor to the clipboard.
* [x] Ctrl + u   Cut/delete the Line before the cursor to the clipboard.

* [ ] Alt + t   Swap current word with previous
* [ ] Ctrl + t   Swap the last two characters before the cursor (typo).
* [ ] Esc  + t   Swap the last two words before the cursor.

* [ ] ctrl + y   Paste the last thing to be cut (yank)
* [ ] Alt + u   UPPER capitalize every character from the cursor to the end of the current word.
* [ ] Alt + l   Lower the case of every character from the cursor to the end of the current word.
* [ ] Alt + c   Capitalize the character under the cursor and move to the end of the word.
* [ ] Alt + r   Cancel the changes and put back the line as it was in the history (revert).
* [ ] ctrl + _   Undo

*/

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
