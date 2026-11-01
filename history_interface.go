package prompt

import (
	istrings "github.com/elk-language/go-prompt/strings"
)

// HistoryInterface lets users replace the builtin history.
type HistoryInterface interface {
	Add(string)
	Clear()
	Older(*Buffer, istrings.Width, int) (*Buffer, bool)
	Newer(*Buffer, istrings.Width, int) (*Buffer, bool)
	Get(i int) (string, bool)
	Entries() []string
	DeleteAll()
}
