package prompt

import (
	"strings"
)

// History stores the texts that are entered.
type History struct {
	histories []string
	tmp       []string
	selected  int
	searchAt  int
}

// Add to add text in history.
func (h *History) Add(input string) {
	h.histories = append(h.histories, input)
	h.Clear()
}

// Clear to clear the history.
func (h *History) Clear() {
	h.tmp = make([]string, len(h.histories))
	for i := range h.histories {
		h.tmp[i] = h.histories[i]
	}
	h.tmp = append(h.tmp, "")
	h.selected = len(h.tmp) - 1
	h.searchAt = -1
}

func (h *History) SearchReset(begin bool) {
	hlen := len(h.histories)
	if begin {
		h.searchAt = 0
	} else {
		h.searchAt = hlen - 1
	}
}

func (h *History) Search(pattern string, fwd bool, skipCur bool) string {
	hlen := len(h.histories)
	if skipCur {
		if h.searchAt < 0 && fwd {
			h.searchAt = 0
		}
		if h.searchAt >= hlen && !fwd {
			h.searchAt = hlen -1
		}
	}
	if h.searchAt < 0 || h.searchAt >= hlen {
		return ""
	}
	if fwd {
		for idx := h.searchAt; idx < hlen; idx++ {
			hstr := h.histories[idx]
			if strings.Contains(hstr, pattern) {
				if skipCur {
					h.searchAt = idx + 1
					skipCur = false
				} else {
					return hstr
				}
			}
		}
		if skipCur {
			h.searchAt = hlen
		}
	} else {
		for idx := h.searchAt; idx >= 0; idx-- {
			hstr := h.histories[idx]
			if strings.Contains(hstr, pattern) {
				if skipCur {
					h.searchAt = idx - 1
					skipCur = false
				} else {
					return hstr
				}
			}
		}
		if skipCur {
			h.searchAt = -1
		}
	}
	return ""
}

// Older saves a buffer of current line and get a buffer of previous line by up-arrow.
// The changes of line buffers are stored until new history is created.
func (h *History) Older(buf *Buffer) (new *Buffer, changed bool) {
	if len(h.tmp) == 1 || h.selected == 0 {
		return buf, false
	}
	h.tmp[h.selected] = buf.Text()

	h.selected--
	new = NewBuffer()
	new.InsertText(h.tmp[h.selected], false, true)
	return new, true
}

// Newer saves a buffer of current line and get a buffer of next line by up-arrow.
// The changes of line buffers are stored until new history is created.
func (h *History) Newer(buf *Buffer) (new *Buffer, changed bool) {
	if h.selected >= len(h.tmp)-1 {
		return buf, false
	}
	h.tmp[h.selected] = buf.Text()

	h.selected++
	new = NewBuffer()
	new.InsertText(h.tmp[h.selected], false, true)
	return new, true
}

// NewHistory returns new history object.
func NewHistory() *History {
	return &History{
		histories: []string{},
		tmp:       []string{""},
		selected:  0,
		searchAt:  -1,
	}
}
