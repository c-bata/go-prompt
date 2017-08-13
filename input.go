package prompt

func dummyExecutor(in string) { return }

func Input(prefix string, completer Completer, opts ...option) string {
	pt := New(dummyExecutor, completer)
	pt.renderer.prefixTextColor = DefaultColor
	pt.renderer.prefix = prefix

	for _, opt := range opts {
		if err := opt(pt); err != nil {
			panic(err)
		}
	}
	return pt.Input()
}

func Choose(prefix string, choices []string, opts ...option) string {
	completer := newChoiceCompleter(choices, FilterHasPrefix)
	pt := New(dummyExecutor, completer)
	pt.renderer.prefixTextColor = DefaultColor
	pt.renderer.prefix = prefix

	for _, opt := range opts {
		if err := opt(pt); err != nil {
			panic(err)
		}
	}
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
