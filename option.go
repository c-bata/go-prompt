package prompt

// Option is the type to replace default parameters.
// prompt.New accepts any number of options (this is functional option pattern).
type Option func(prompt *Prompt) error

// OptionParser to set a custom ConsoleParser object. An argument should implement ConsoleParser interface.
func OptionParser(x ConsoleParser) Option {
	return func(p *Prompt) error {
		p.in = x
		return nil
	}
}

// OptionWriter to set a custom ConsoleWriter object. An argument should implement ConsoleWriter interace.
func OptionWriter(x ConsoleWriter) Option {
	return func(p *Prompt) error {
		p.renderer.out = x
		return nil
	}
}

// OptionTitle to set title displayed at the header bar of terminal.
func OptionTitle(x string) Option {
	return func(p *Prompt) error {
		p.renderer.title = x
		return nil
	}
}

// OptionPrefix to set prefix string.
func OptionPrefix(x string) Option {
	return func(p *Prompt) error {
		p.renderer.prefix = x
		return nil
	}
}

// OptionLivePrefix to change the prefix dynamically by callback function
func OptionLivePrefix(f func() (prefix string, useLivePrefix bool)) Option {
	return func(p *Prompt) error {
		p.renderer.livePrefixCallback = f
		return nil
	}
}

func OptionPrefixTextColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.prefixTextColor = x
		return nil
	}
}

func OptionPrefixBackgroundColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.prefixBGColor = x
		return nil
	}
}

func OptionInputTextColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.inputTextColor = x
		return nil
	}
}

func OptionInputBGColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.inputBGColor = x
		return nil
	}
}

func OptionPreviewSuggestionTextColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.previewSuggestionTextColor = x
		return nil
	}
}

func OptionPreviewSuggestionBGColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.previewSuggestionBGColor = x
		return nil
	}
}

func OptionSuggestionTextColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.suggestionTextColor = x
		return nil
	}
}

func OptionSuggestionBGColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.suggestionBGColor = x
		return nil
	}
}

func OptionSelectedSuggestionTextColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.selectedSuggestionTextColor = x
		return nil
	}
}

func OptionSelectedSuggestionBGColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.selectedSuggestionBGColor = x
		return nil
	}
}

func OptionDescriptionTextColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.descriptionTextColor = x
		return nil
	}
}

func OptionDescriptionBGColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.descriptionBGColor = x
		return nil
	}
}

func OptionSelectedDescriptionTextColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.selectedDescriptionTextColor = x
		return nil
	}
}

func OptionSelectedDescriptionBGColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.selectedDescriptionBGColor = x
		return nil
	}
}

func OptionScrollbarThumbColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.scrollbarThumbColor = x
		return nil
	}
}

func OptionScrollbarBGColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.scrollbarBGColor = x
		return nil
	}
}

// OptionMaxSuggestion specify the max number of displayed suggestions.
func OptionMaxSuggestion(x uint16) Option {
	return func(p *Prompt) error {
		p.completion.max = x
		return nil
	}
}

// OptionHistory to set history expressed by string array.
func OptionHistory(x []string) Option {
	return func(p *Prompt) error {
		p.history.histories = x
		p.history.Clear()
		return nil
	}
}

// OptionSwitchKeyBindMode set a key bind mode.
func OptionSwitchKeyBindMode(m KeyBindMode) Option {
	return func(p *Prompt) error {
		p.keyBindMode = m
		return nil
	}
}

// SwitchKeyBindMode to set a key bind mode.
// Deprecated: Please use OptionSwitchKeyBindMode.
var SwitchKeyBindMode = OptionSwitchKeyBindMode

// OptionAddKeyBind to set a custom key bind.
func OptionAddKeyBind(b ...KeyBind) Option {
	return func(p *Prompt) error {
		p.keyBindings = append(p.keyBindings, b...)
		return nil
	}
}

// New returns a Prompt with powerful auto-completion.
func New(executor Executor, completer Completer, opts ...Option) *Prompt {
	pt := &Prompt{
		in: NewStandardInputParser(),
		renderer: &Render{
			prefix:                       "> ",
			out:                          NewStandardOutputWriter(),
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
		completion:  NewCompletionManager(completer, 6),
		keyBindMode: EmacsKeyBind, // All the above assume that bash is running in the default Emacs setting
	}

	for _, opt := range opts {
		if err := opt(pt); err != nil {
			panic(err)
		}
	}
	return pt
}
