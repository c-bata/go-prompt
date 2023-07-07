package prompt

func NoopExecutor(in string) {}

// Input get the input data from the user and return it.
func Input(opts ...Option) string {
	pt := New(NoopExecutor)
	pt.renderer.prefixTextColor = DefaultColor

	for _, opt := range opts {
		if err := opt(pt); err != nil {
			panic(err)
		}
	}
	return pt.Input()
}
