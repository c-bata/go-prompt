package prompt

func NoopExecutor(in string) {}

// Input get the input data from the user and return it.
func Input(prefix string, opts ...Option) string {
	pt := New(NoopExecutor)
	pt.renderer.prefixTextColor = DefaultColor
	pt.renderer.prefix = prefix

	for _, opt := range opts {
		if err := opt(pt); err != nil {
			panic(err)
		}
	}
	return pt.Input()
}
