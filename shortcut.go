package prompt

// Input get the input data from the user and return it.
func Input(prefix string, completer Completer, opts ...Option) string {
	pt := New(nil, completer, opts...)
	pt.rendererOptions = append(pt.rendererOptions, func(r *Renderer) {
		r.prefix = prefix
	})
	return pt.Input()
}

