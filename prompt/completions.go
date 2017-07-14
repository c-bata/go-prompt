package prompt

type Completion struct {
	// The new string that will be inserted into document.
	text          string
	// Position relative to the cursor position where the new text will start.
	startPosition int
}

func (c *Completion) NewCompletionFromPosition(position int) *Completion {
	if position < c.startPosition {
		panic("position argument must be smaller than start position.")
	}

	return &Completion{
		text: c.text[position - c.startPosition:],
	}
}

func NewCompletion(text string) *Completion {
	return &Completion{
		text: text,
	}
}
