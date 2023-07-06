package prompt

// Option is the type to replace default parameters.
// prompt.New accepts any number of options (this is functional option pattern).
type Option func(prompt *Prompt) error

// WithCompleter is an option that sets a custom Completer object.
func WithCompleter(c Completer) Option {
	return func(p *Prompt) error {
		p.completion.completer = c
		return nil
	}
}

// WithReader to set a custom Reader object. An argument should implement Reader interface.
func WithReader(x Reader) Option {
	return func(p *Prompt) error {
		p.reader = x
		return nil
	}
}

// WithWriter to set a custom Writer object. An argument should implement Writer interface.
func WithWriter(x Writer) Option {
	return func(p *Prompt) error {
		registerWriter(x)
		p.renderer.out = x
		return nil
	}
}

// WithTitle to set title displayed at the header bar of terminal.
func WithTitle(x string) Option {
	return func(p *Prompt) error {
		p.renderer.title = x
		return nil
	}
}

// WithPrefix to set prefix string.
func WithPrefix(x string) Option {
	return func(p *Prompt) error {
		p.renderer.prefix = x
		return nil
	}
}

// WithInitialBufferText to set the initial buffer text
func WithInitialBufferText(x string) Option {
	return func(p *Prompt) error {
		p.buf.InsertText(x, false, true)
		return nil
	}
}

// WithCompletionWordSeparator to set word separators. Enable only ' ' if empty.
func WithCompletionWordSeparator(x string) Option {
	return func(p *Prompt) error {
		p.completion.wordSeparator = x
		return nil
	}
}

// WithLivePrefix to change the prefix dynamically by callback function
func WithLivePrefix(f func() (prefix string, useLivePrefix bool)) Option {
	return func(p *Prompt) error {
		p.renderer.livePrefixCallback = f
		return nil
	}
}

// WithPrefixTextColor change a text color of prefix string
func WithPrefixTextColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.prefixTextColor = x
		return nil
	}
}

// WithPrefixBackgroundColor to change a background color of prefix string
func WithPrefixBackgroundColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.prefixBGColor = x
		return nil
	}
}

// WithInputTextColor to change a color of text which is input by user
func WithInputTextColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.inputTextColor = x
		return nil
	}
}

// WithInputBGColor to change a color of background which is input by user
func WithInputBGColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.inputBGColor = x
		return nil
	}
}

// WithPreviewSuggestionTextColor to change a text color which is completed
func WithPreviewSuggestionTextColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.previewSuggestionTextColor = x
		return nil
	}
}

// WithPreviewSuggestionBGColor to change a background color which is completed
func WithPreviewSuggestionBGColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.previewSuggestionBGColor = x
		return nil
	}
}

// WithSuggestionTextColor to change a text color in drop down suggestions.
func WithSuggestionTextColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.suggestionTextColor = x
		return nil
	}
}

// WithSuggestionBGColor change a background color in drop down suggestions.
func WithSuggestionBGColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.suggestionBGColor = x
		return nil
	}
}

// WithSelectedSuggestionTextColor to change a text color for completed text which is selected inside suggestions drop down box.
func WithSelectedSuggestionTextColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.selectedSuggestionTextColor = x
		return nil
	}
}

// WithSelectedSuggestionBGColor to change a background color for completed text which is selected inside suggestions drop down box.
func WithSelectedSuggestionBGColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.selectedSuggestionBGColor = x
		return nil
	}
}

// WithDescriptionTextColor to change a background color of description text in drop down suggestions.
func WithDescriptionTextColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.descriptionTextColor = x
		return nil
	}
}

// WithDescriptionBGColor to change a background color of description text in drop down suggestions.
func WithDescriptionBGColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.descriptionBGColor = x
		return nil
	}
}

// WithSelectedDescriptionTextColor to change a text color of description which is selected inside suggestions drop down box.
func WithSelectedDescriptionTextColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.selectedDescriptionTextColor = x
		return nil
	}
}

// WithSelectedDescriptionBGColor to change a background color of description which is selected inside suggestions drop down box.
func WithSelectedDescriptionBGColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.selectedDescriptionBGColor = x
		return nil
	}
}

// WithScrollbarThumbColor to change a thumb color on scrollbar.
func WithScrollbarThumbColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.scrollbarThumbColor = x
		return nil
	}
}

// WithScrollbarBGColor to change a background color of scrollbar.
func WithScrollbarBGColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.scrollbarBGColor = x
		return nil
	}
}

// WithMaxSuggestion specify the max number of displayed suggestions.
func WithMaxSuggestion(x uint16) Option {
	return func(p *Prompt) error {
		p.completion.max = x
		return nil
	}
}

// WithHistory to set history expressed by string array.
func WithHistory(x []string) Option {
	return func(p *Prompt) error {
		p.history.histories = x
		p.history.Clear()
		return nil
	}
}

// WithSwitchKeyBindMode set a key bind mode.
func WithSwitchKeyBindMode(m KeyBindMode) Option {
	return func(p *Prompt) error {
		p.keyBindMode = m
		return nil
	}
}

// WithCompletionOnDown allows for Down arrow key to trigger completion.
func WithCompletionOnDown() Option {
	return func(p *Prompt) error {
		p.completionOnDown = true
		return nil
	}
}

// SwitchKeyBindMode to set a key bind mode.
// Deprecated: Please use WithSwitchKeyBindMode.
var SwitchKeyBindMode = WithSwitchKeyBindMode

// WithAddKeyBind to set a custom key bind.
func WithAddKeyBind(b ...KeyBind) Option {
	return func(p *Prompt) error {
		p.keyBindings = append(p.keyBindings, b...)
		return nil
	}
}

// WithAddASCIICodeBind to set a custom key bind.
func WithAddASCIICodeBind(b ...ASCIICodeBind) Option {
	return func(p *Prompt) error {
		p.ASCIICodeBindings = append(p.ASCIICodeBindings, b...)
		return nil
	}
}

// WithShowCompletionAtStart to set completion window is open at start.
func WithShowCompletionAtStart() Option {
	return func(p *Prompt) error {
		p.completion.showAtStart = true
		return nil
	}
}

// WithBreakLineCallback to run a callback at every break line
func WithBreakLineCallback(fn func(*Document)) Option {
	return func(p *Prompt) error {
		p.renderer.breakLineCallback = fn
		return nil
	}
}

// WithExitChecker set an exit function which checks if go-prompt exits its Run loop
func WithExitChecker(fn ExitChecker) Option {
	return func(p *Prompt) error {
		p.exitChecker = fn
		return nil
	}
}

// WithLexer set lexer function and enable it.
func WithLexer(lex Lexer) Option {
	return func(p *Prompt) error {
		p.lexer = lex
		return nil
	}
}

// New returns a Prompt with powerful auto-completion.
func New(executor Executor, opts ...Option) *Prompt {
	defaultWriter := NewStdoutWriter()
	registerWriter(defaultWriter)

	pt := &Prompt{
		reader: NewStdinReader(),
		renderer: &Render{
			prefix:                       "> ",
			out:                          defaultWriter,
			livePrefixCallback:           func() (string, bool) { return "", false },
			prefixTextColor:              Blue,
			prefixBGColor:                DefaultColor,
			inputTextColor:               DefaultColor,
			inputBGColor:                 DefaultColor,
			previewSuggestionTextColor:   Green,
			previewSuggestionBGColor:     DefaultColor,
			suggestionTextColor:          White,
			suggestionBGColor:            Cyan,
			selectedSuggestionTextColor:  Black,
			selectedSuggestionBGColor:    Turquoise,
			descriptionTextColor:         Black,
			descriptionBGColor:           Turquoise,
			selectedDescriptionTextColor: White,
			selectedDescriptionBGColor:   Cyan,
			scrollbarThumbColor:          DarkGray,
			scrollbarBGColor:             Cyan,
		},
		buf:         NewBuffer(),
		executor:    executor,
		history:     NewHistory(),
		completion:  NewCompletionManager(6),
		keyBindMode: EmacsKeyBind, // All the above assume that bash is running in the default Emacs setting
	}

	for _, opt := range opts {
		if err := opt(pt); err != nil {
			panic(err)
		}
	}
	return pt
}
