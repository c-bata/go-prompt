package filter

import (
	"strings"

	"github.com/c-bata/go-prompt"
)

// HasPrefix checks whether the string completions.Text begins with sub.
func HasPrefix(suggests []prompt.Suggest, sub string, ignoreCase bool) []prompt.Suggest {
	if sub == "" {
		return suggests
	}
	if ignoreCase {
		sub = strings.ToUpper(sub)
	}

	ret := make([]prompt.Suggest, 0, len(suggests))
	for i := range suggests {
		c := suggests[i].Text
		if ignoreCase {
			c = strings.ToUpper(c)
		}
		if strings.HasPrefix(c, sub) {
			ret = append(ret, suggests[i])
		}
	}
	return ret
}

// HasSuffix checks whether the completion.Text ends with sub.
func HasSuffix(suggests []prompt.Suggest, sub string, ignoreCase bool) []prompt.Suggest {
	if sub == "" {
		return suggests
	}
	if ignoreCase {
		sub = strings.ToUpper(sub)
	}

	ret := make([]prompt.Suggest, 0, len(suggests))
	for i := range suggests {
		c := suggests[i].Text
		if ignoreCase {
			c = strings.ToUpper(c)
		}
		if strings.HasSuffix(c, sub) {
			ret = append(ret, suggests[i])
		}
	}
	return ret
}

// Contains checks whether the completion.Text contains sub.
func Contains(suggests []prompt.Suggest, sub string, ignoreCase bool) []prompt.Suggest {
	if sub == "" {
		return suggests
	}
	if ignoreCase {
		sub = strings.ToUpper(sub)
	}

	ret := make([]prompt.Suggest, 0, len(suggests))
	for i := range suggests {
		c := suggests[i].Text
		if ignoreCase {
			c = strings.ToUpper(c)
		}
		if strings.Contains(c, sub) {
			ret = append(ret, suggests[i])
		}
	}
	return ret
}
