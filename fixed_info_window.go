package prompt

import (
	"fmt"
)

type FixedInfoWindow struct {
	lines    []*string
	maxLines int
	InfoWindow
}

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

func (l *FixedInfoWindow) AddLine(line *string) error {
	if len(l.lines) <= l.maxLines {
		l.lines = append(l.lines, line)
		return nil
	}

	for i, il := range l.lines {
		if il != nil {
			l.lines[i] = line
			return nil
		}
	}
	return fmt.Errorf("No free lines")
}

func (l *FixedInfoWindow) Clear() {
	for i, _ := range l.lines {
		l.lines[i] = nil
	}
}

func (l *FixedInfoWindow) ClearLine(line int) {
	if line < 0 || line >= len(l.lines) {
		return
	}

	l.lines[line] = nil
}

func NewFixedInfoWindow(lines int) *FixedInfoWindow {
	ret := new(FixedInfoWindow)
	ret.lines = []*string{}
	ret.maxLines = lines

	return ret
}
