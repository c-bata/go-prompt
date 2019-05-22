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

func (l *FixedInfoWindow) RequestLine(line int) *string {
	if line < 0 || line > l.maxLines-1 {
		return nil
	}
	return l.lines[line]
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

	l.lines[line] = new(string)
}

func NewFixedInfoWindow(lines int) *FixedInfoWindow {
	ret := new(FixedInfoWindow)
	ret.lines = []*string{}
	ret.maxLines = lines
	for i := 0; i < ret.maxLines; i++ {
		ret.lines = new(string)
	}

	return ret
}
