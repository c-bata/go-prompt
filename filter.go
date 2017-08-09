package prompt

import "strings"

type Filter func([]Suggest, string, bool) []Suggest

func FilterHasPrefix(completions []Suggest, sub string, ignoreCase bool) []Suggest {
	if sub == "" {
		return completions
	}
	if ignoreCase {
		sub = strings.ToUpper(sub)
	}

	ret := make([]Suggest, 0, len(completions))
	for i := range completions {
		c := completions[i].Text
		if ignoreCase {
			c = strings.ToUpper(c)
		}
		if strings.HasPrefix(c, sub) {
			ret = append(ret, completions[i])
		}
	}
	return ret
}

func FilterHasSuffix(completions []Suggest, sub string, ignoreCase bool) []Suggest {
	if sub == "" {
		return completions
	}
	if ignoreCase {
		sub = strings.ToUpper(sub)
	}

	ret := make([]Suggest, 0, len(completions))
	for i := range completions {
		c := completions[i].Text
		if ignoreCase {
			c = strings.ToUpper(c)
		}
		if strings.HasSuffix(c, sub) {
			ret = append(ret, completions[i])
		}
	}
	return ret
}

func FilterContains(completions []Suggest, sub string, ignoreCase bool) []Suggest {
	if sub == "" {
		return completions
	}
	if ignoreCase {
		sub = strings.ToUpper(sub)
	}

	ret := make([]Suggest, 0, len(completions))
	for i := range completions {
		c := completions[i].Text
		if ignoreCase {
			c = strings.ToUpper(c)
		}
		if strings.Contains(c, sub) {
			ret = append(ret, completions[i])
		}
	}
	return ret
}
