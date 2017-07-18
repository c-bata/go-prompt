package prompt

import "strings"

type CompletionFilter func([]Completion, string, bool) []Completion

func FilterHasPrefix(completions []Completion, sub string, ignoreCase bool) []Completion {
	if sub == "" {
		return completions
	}
	if ignoreCase {
		sub = strings.ToUpper(sub)
	}

	ret := make([]Completion, 0, len(completions))
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

func FilterHasSuffix(completions []Completion, sub string, ignoreCase bool) []Completion {
	if sub == "" {
		return completions
	}
	if ignoreCase {
		sub = strings.ToUpper(sub)
	}

	ret := make([]Completion, 0, len(completions))
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

func FilterContains(completions []Completion, sub string, ignoreCase bool) []Completion {
	if sub == "" {
		return completions
	}
	if ignoreCase {
		sub = strings.ToUpper(sub)
	}

	ret := make([]Completion, 0, len(completions))
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
