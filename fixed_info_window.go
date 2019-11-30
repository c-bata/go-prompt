package prompt

import ()

// FixedInfoWindow is the info window that
// uses fixed number of lines and is not scrollable
type FixedInfoWindow struct {
	lines    []*string
	maxLines int
	InfoWindow
}

// GetLines returns an array of the Lines from the info window
// count tells how many lines should be returned
func (l *FixedInfoWindow) GetLines(count int) []string {
	ret := []string{}

	if count >= len(l.lines) {
		for i := 0; i < len(l.lines); i++ {
			ret = append(ret, *l.lines[i])
		}
	} else {
		for i := 0; i < count; i++ {
			ret = append(ret, *l.lines[i])
		}
	}
	return ret
}

// RequestLine will return a pointer to one specific line
// which can be updated to new content
func (l *FixedInfoWindow) RequestLine(line int) *string {
	if line < 0 || line > l.maxLines-1 {
		return nil
	}
	return l.lines[line]
}

// Clear cleans the whole info window
func (l *FixedInfoWindow) Clear() {
	for i, _ := range l.lines {
		*l.lines[i] = ""
	}
}

// ClearLine will clear one specific line of the info window
func (l *FixedInfoWindow) ClearLine(line int) {
	if line < 0 || line >= len(l.lines) {
		return
	}

	l.lines[line] = new(string)
}

// Len returns the number of lines that the info window
// can print
func (l *FixedInfoWindow) Len() int {
	return l.maxLines
}

// NewFixedInfoWindow will create a new fixed size info window
// with the number of lines given by the lines parameter
func NewFixedInfoWindow(lines int) *FixedInfoWindow {
	ret := new(FixedInfoWindow)
	ret.lines = []*string{}
	ret.maxLines = lines
	for i := 0; i < ret.maxLines; i++ {
		ret.lines = append(ret.lines, new(string))
	}

	return ret
}
