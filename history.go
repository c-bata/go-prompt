package prompt

import "log"

type History struct {
	histories []string
	tmp       []string
	selected  int
}

func (h *History) Add(input string) {
	h.histories = append(h.histories, input)
	h.Clear()
}

func (h *History) Clear() {
	copy(h.tmp, h.histories)
	h.tmp = append(h.tmp, "")
	h.selected = len(h.tmp) - 1
}

func (h *History) Older(buf *Buffer) (new *Buffer, changed bool) {
	log.Printf("[DEBUG] Before %#v\n", h)
	if len(h.tmp) == 1 || h.selected == 0 {
		return buf, false
	}
	h.tmp[h.selected] = buf.Text()

	h.selected--
	new = NewBuffer()
	new.InsertText(h.tmp[h.selected], false, true)
	log.Printf("[DEBUG] After %#v\n", h)
	return new, true
}

func (h *History) Newer(buf *Buffer) (new *Buffer, changed bool) {
	log.Printf("[DEBUG] Before %#v\n", h)
	if h.selected >= len(h.tmp)-1 {
		return buf, false
	}
	h.tmp[h.selected] = buf.Text()

	h.selected++
	new = NewBuffer()
	new.InsertText(h.tmp[h.selected], false, true)
	log.Printf("[DEBUG] After %#v\n", h)
	return new, true
}

func NewHistory() *History {
	return &History{
		histories: []string{},
		tmp:       []string{""},
		selected:  0,
	}
}
