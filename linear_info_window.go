package prompt

type LinearInfoWindow struct {
	lines      []string
	maxLines   int
	actLine    int
	autoscroll bool
	InfoWindow
}

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

func (l *LinearInfoWindow) LineUp() {
	if l.actLine > 0 {
		l.actLine--
	}
}

func (l *LinearInfoWindow) LineDown() {
	if l.actLine < len(l.lines)-1 {
		l.actLine++
	}
}

func (l *LinearInfoWindow) Clear() {
	l.lines = []string{}
	l.actLine = 0
}

func (l *LinearInfoWindow) ClearLine(line int) {
	if line < 0 || line >= len(l.lines) {
		return
	}

	l.lines = append(l.lines[:line], l.lines[line+1:]...)
}

func NewLinearInfoWindow(maxLines int, autoscroll bool) *LinearInfoWindow {
	ret := new(LinearInfoWindow)
	ret.actLine = 0
	ret.maxLines = maxLines
	ret.autoscroll = autoscroll
	ret.lines = []string{}

	return ret
}
