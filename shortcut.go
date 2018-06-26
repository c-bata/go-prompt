package prompt

// Input get the input data from the user and return it.
func Input(prefix string, completer Completer, opts ...Option) string {
	pt := New(nil, completer, opts...)
	pt.rendererOptions = append(pt.rendererOptions, func(r *Renderer) {
		r.prefix = prefix
	})
	return pt.Input()
}

// Choose to the shortcut of input function to select from string array.
// Deprecated: Maybe anyone want to use this.
func Choose(prefix string, choices []string, opts ...Option) string {
	completer := newChoiceCompleter(choices, FilterHasPrefix)
	pt := New(nil, completer, opts...)
	pt.rendererOptions = append(pt.rendererOptions, func(r *Renderer) {
		r.prefix = prefix
	})
	return pt.Input()
}

func newChoiceCompleter(choices []string, filter Filter) Completer {
	s := make([]Suggest, len(choices))
	for i := range choices {
		s[i] = Suggest{Text: choices[i]}
	}
	return func(x Document) []Suggest {
		return filter(s, x.GetWordBeforeCursor(), true)
	}
}
