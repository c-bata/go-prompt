package prompt

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

// History stores the texts that are entered.
type History struct {
	histories []string
	tmp       []string
	selected  int
	size      int
}

// Add to add text in history.
func (h *History) Add(input string) {
	h.histories = append(h.histories, input)
	if len(h.histories) > h.size {
		h.histories = h.histories[1:]
	}
	h.Clear()
}

func (h *History) Get(i int) string {
	if i < 0 || i >= len(h.histories) {
		return ""
	}
	return h.histories[i]
}

func (h *History) Entries() []string {
	return h.histories
}

func (h *History) List(bShouldUsePager bool) {
	if len(h.histories) <= 0 {
		return
	}
	tableString := &bytes.Buffer{}
	table := tablewriter.NewWriter(tableString)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetAutoWrapText(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT})
	data := [][]string{}
	for i, s := range h.histories {
		data = append(data, []string{strconv.Itoa(i + 1), s})
	}
	table.AppendBulk(data)
	table.Render()
	b := tableString.Bytes()
	if bShouldUsePager {
		pager_print(b, nil)
	} else {
		fmt.Printf(string(b))
	}
}

func (h *History) DeleteAll() {
	(*h).histories = []string{}
	(*h).tmp = []string{""}
	(*h).selected = 0
}

// Clear to clear the history.
func (h *History) Clear() {
	h.tmp = make([]string, len(h.histories))
	for i := range h.histories {
		h.tmp[i] = h.histories[i]
	}
	h.tmp = append(h.tmp, "")
	h.selected = len(h.tmp) - 1
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
		size:      1000,
	}
}
