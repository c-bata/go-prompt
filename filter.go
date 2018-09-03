package prompt

import "strings"

// Filter is the type to filter the prompt.Suggestion array.
type Filter func([]Suggest, string, bool) []Suggest

// FilterHasPrefix checks whether the string completions.Text begins with sub.
func FilterHasPrefix(completions []Suggest, sub string, ignoreCase bool) []Suggest {
	return filterCommon(completions, sub, ignoreCase, strings.HasPrefix)
}

// FilterHasSuffix checks whether the completion.Text ends with sub.
func FilterHasSuffix(completions []Suggest, sub string, ignoreCase bool) []Suggest {
	return filterCommon(completions, sub, ignoreCase, strings.HasSuffix)
}

// FilterContains checks whether the completion.Text contains sub.
func FilterContains(completions []Suggest, sub string, ignoreCase bool) []Suggest {
	return filterCommon(completions, sub, ignoreCase, strings.Contains)
}

// FilterFuzzy checks whether the completion.Text fuzzy matches sub.
// Fuzzy searching for "dog" is equivalent to "*d*o*g*". This search term
// would match, for example, "Good food is gone"
//                               ^  ^      ^
func FilterFuzzy(completions []Suggest, sub string, ignoreCase bool) []Suggest {
	return filterCommon(completions, sub, ignoreCase, fuzzyMatch)
}

func fuzzyMatch(s, sub string) bool {
	sChars := []rune(s)
	subChars := []rune(sub)
	sIdx := 0

	for _, c := range subChars {
		found := false
		for ; sIdx < len(sChars); sIdx++ {
			if sChars[sIdx] == c {
				found = true
				sIdx++
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func filterCommon(completions []Suggest, sub string, ignoreCase bool, test func(string, string) bool) []Suggest {
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
		if test(c, sub) {
			ret = append(ret, completions[i])
		}
	}
	return ret
}
