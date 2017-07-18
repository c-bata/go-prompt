package prompt

import "strings"

func FilterHasPrefix(completions []string, sub string, ignoreCase bool) []string {
	if sub == "" {
		return completions
	}
	if ignoreCase {
		sub = strings.ToUpper(sub)
	}

	ret := make([]string, 0, len(completions))
	for i, n := range completions {
		if ignoreCase {
			n = strings.ToUpper(n)
		}
		if strings.HasPrefix(n, sub) {
			ret = append(ret, completions[i])
		}
	}
	return ret
}

func FilterHasSuffix(completions []string, sub string, ignoreCase bool) []string {
	if sub == "" {
		return completions
	}
	if ignoreCase {
		sub = strings.ToUpper(sub)
	}

	ret := make([]string, 0, len(completions))
	for i, n := range completions {
		if ignoreCase {
			n = strings.ToUpper(n)
		}
		if strings.HasSuffix(n, sub) {
			ret = append(ret, completions[i])
		}
	}
	return ret
}

func FilterContains(completions []string, sub string, ignoreCase bool) []string {
	if sub == "" {
		return completions
	}
	if ignoreCase {
		sub = strings.ToUpper(sub)
	}

	ret := make([]string, 0, len(completions))
	for i, n := range completions {
		if ignoreCase {
			n = strings.ToUpper(n)
		}
		if strings.Contains(n, sub) {
			ret = append(ret, completions[i])
		}
	}
	return ret
}
