package prompt

import "strings"

// Filter is the type to filter the prompt.Suggestion array.
type Filter func([]Suggest, string, bool) []Suggest

// FilterHasPrefix checks whether the string completions.Text begins with sub.
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

// FilterHasSuffix checks whether the completion.Text ends with sub.
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

// FilterContains checks whether the completion.Text contains sub.
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
