package prompt

import "syscall"

type option func(prompt *Prompt) error

func OptionParser(x ConsoleParser) option {
	return func(p *Prompt) error {
		p.in = x
		return nil
	}
}

func OptionWriter(x ConsoleWriter) option {
	return func(p *Prompt) error {
		p.renderer.out = x
		return nil
	}
}

func OptionTitle(x string) option {
	return func(p *Prompt) error {
		p.renderer.title = x
		return nil
	}
}

func OptionPrefix(x string) option {
	return func(p *Prompt) error {
		p.renderer.prefix = x
		return nil
	}
}

func OptionPrefixTextColor(x Color) option {
	return func(p *Prompt) error {
		p.renderer.prefixTextColor = x
		return nil
	}
}

func OptionPrefixBackgroundColor(x Color) option {
	return func(p *Prompt) error {
		p.renderer.prefixBGColor = x
		return nil
	}
}

func OptionInputTextColor(x Color) option {
	return func(p *Prompt) error {
		p.renderer.inputTextColor = x
		return nil
	}
}

func OptionInputBGColor(x Color) option {
	return func(p *Prompt) error {
		p.renderer.inputBGColor = x
		return nil
	}
}

func OptionPreviewSuggestionTextColor(x Color) option {
	return func(p *Prompt) error {
		p.renderer.previewSuggestionTextColor = x
		return nil
	}
}

func OptionPreviewSuggestionBGColor(x Color) option {
	return func(p *Prompt) error {
		p.renderer.previewSuggestionBGColor = x
		return nil
	}
}

func OptionSuggestionTextColor(x Color) option {
	return func(p *Prompt) error {
		p.renderer.suggestionTextColor = x
		return nil
	}
}

func OptionSuggestionBGColor(x Color) option {
	return func(p *Prompt) error {
		p.renderer.suggestionBGColor = x
		return nil
	}
}

func OptionSelectedSuggestionTextColor(x Color) option {
	return func(p *Prompt) error {
		p.renderer.selectedSuggestionTextColor = x
		return nil
	}
}

func OptionSelectedSuggestionBGColor(x Color) option {
	return func(p *Prompt) error {
		p.renderer.selectedSuggestionBGColor = x
		return nil
	}
}

func OptionDescriptionTextColor(x Color) option {
	return func(p *Prompt) error {
		p.renderer.descriptionTextColor = x
		return nil
	}
}

func OptionDescriptionBGColor(x Color) option {
	return func(p *Prompt) error {
		p.renderer.descriptionBGColor = x
		return nil
	}
}

func OptionSelectedDescriptionTextColor(x Color) option {
	return func(p *Prompt) error {
		p.renderer.selectedDescriptionTextColor = x
		return nil
	}
}

func OptionSelectedDescriptionBGColor(x Color) option {
	return func(p *Prompt) error {
		p.renderer.selectedDescriptionBGColor = x
		return nil
	}
}

func OptionMaxSuggestion(x uint16) option {
	return func(p *Prompt) error {
		p.completion.max = x
		return nil
	}
}

func OptionHistory(x []string) option {
	return func(p *Prompt) error {
		p.history.histories = x
		p.history.Clear()
		return nil
	}
}

func SwitchKeyBindMode(m KeyBindMode) option {
	return func(p *Prompt) error {
		p.keyBindMode = m
		return nil
	}
}

func OptionAddKeyBind(b ...KeyBind) option {
	return func(p *Prompt) error {
		p.keyBindings = append(p.keyBindings, b...)
		return nil
	}
}

func New(executor Executor, completer Completer, opts ...option) *Prompt {
	pt := &Prompt{
		in: &VT100Parser{fd: syscall.Stdin},
		renderer: &Render{
			prefix:                       "> ",
			out:                          &VT100Writer{fd: syscall.Stdout},
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
