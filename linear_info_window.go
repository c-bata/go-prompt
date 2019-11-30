package prompt

// LinearInfoWindow is the info window that
// is scrollable and can be used to have a
// linear output information
type LinearInfoWindow struct {
	lines      []string
	maxLines   int
	actLine    int
	autoscroll bool
	InfoWindow
}

// GetLines returns an array of the Lines from the info window
// count tells how many lines should be returned
func (l *LinearInfoWindow) GetLines(count int) []string {
	if l.actLine < 0 || l.actLine >= len(l.lines) {
		return []string{}
	}

	if l.actLine+count >= len(l.lines) {
		if len(l.lines)-(count+1) >= 0 {
			return l.lines[len(l.lines)-(count):]
		}
		return l.lines
	}

	return l.lines[l.actLine : l.actLine+count]
}

// AddLine will append a new line to the info window
func (l *LinearInfoWindow) AddLine(line string) {
	if l.maxLines != 0 {
		if len(l.lines) >= l.maxLines {
			l.lines = l.lines[1 : len(l.lines)-1]
		}
	}
	l.lines = append(l.lines, line)
	if l.autoscroll {
		l.LineDown()
	}
}

// LineUp will scroll up a line
func (l *LinearInfoWindow) LineUp() {
	if l.actLine > 0 {
		l.actLine--
	}
}

// LineDown will scroll down a line
func (l *LinearInfoWindow) LineDown() {
	if l.actLine < len(l.lines)-1 {
		l.actLine++
	}
}

// Clear cleans the whole info window
func (l *LinearInfoWindow) Clear() {
	l.lines = []string{}
	l.actLine = 0
}

// ClearLine will clear one specific line of the info window
func (l *LinearInfoWindow) ClearLine(line int) {
	if line < 0 || line >= len(l.lines) {
		return
	}

	l.lines = append(l.lines[:line], l.lines[line+1:]...)
}

// NewLinearInfoWindow will create a new info window that is scrollable
// with the maximum number of lines give by the parameter maxLines
// If autoscroll is set to true it will scroll up when adding lines
func NewLinearInfoWindow(maxLines int, autoscroll bool) *LinearInfoWindow {
	ret := new(LinearInfoWindow)
	ret.actLine = 0
	ret.maxLines = maxLines
	ret.autoscroll = autoscroll
	ret.lines = []string{}

	return ret
}
