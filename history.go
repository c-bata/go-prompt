package prompt

// History stores the texts that are entered.
type History struct {
	histories []string
	selected  int
}

// Add to add text in history.
func (h *History) Add(input string) {
	h.histories = append(h.histories, input)
	h.selected = len(h.histories) - 1
}

// Clear to clear the history.
func (h *History) Clear() {
	h.histories = []string{}
	h.selected = -1
}

// Older saves a buffer of current line and get a buffer of previous line by up-arrow.
// The changes of line buffers are stored until new history is created.
func (h *History) Older(buf *Buffer) (new *Buffer, changed bool) {
	if len(h.histories) == 0 || h.selected < 0 {
		return buf, false
	}
	new = NewBuffer()
	new.InsertText(h.histories[h.selected], false, true)
	h.histories[h.selected] = buf.Text()
	h.selected--
	return new, true
}

// Newer saves a buffer of current line and get a buffer of next line by up-arrow.
// The changes of line buffers are stored until new history is created.
func (h *History) Newer(buf *Buffer) (new *Buffer, changed bool) {
	if h.selected >= len(h.histories)-1 {
		return buf, false
	}
	new = NewBuffer()
	new.InsertText(h.histories[h.selected], false, true)
	h.histories[h.selected] = buf.Text()
	h.selected++
	return new, true
}

// NewHistory returns new history object.
func NewHistory() *History {
	return &History{
		histories: []string{},
		selected:  -1,
	}
}
