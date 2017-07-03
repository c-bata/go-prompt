package prompt

// SelectionType expresses how characters selected.
type SelectionType int

const (
	// CHARACTERS selected freely.
	CHARACTERS SelectionType = iota
	// LINES selected current line.
	LINES
	// BLOCK selected the word block.
	BLOCK
)

// SelectionState holds cursor position and selected characters.
type SelectionState struct {
	OriginalCursorPosition int
	Type                   SelectionType
}
