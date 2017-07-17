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

func OptionOutputTextColor(x Color) option {
	return func(p *Prompt) error {
		p.renderer.outputTextColor = x
		return nil
	}
}

func OptionOutputBGColor(x Color) option {
	return func(p *Prompt) error {
		p.renderer.outputBGColor = x
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

func OptionMaxCompletions(x uint16) option {
	return func(p *Prompt) error {
		p.renderer.maxCompletions = x
		return nil
	}
}

func NewPrompt(executor Executor, completer Completer, opts ...option) *Prompt {
	pt := &Prompt{
		in: &VT100Parser{fd: syscall.Stdin},
		renderer: &Render{
			prefix:                      "> ",
			out:                         &VT100Writer{fd: syscall.Stdout},
			prefixTextColor:             Blue,
			prefixBGColor:               DefaultColor,
			inputTextColor:              DefaultColor,
			inputBGColor:                DefaultColor,
			outputTextColor:             DefaultColor,
			outputBGColor:               DefaultColor,
			previewSuggestionTextColor:  Green,
			previewSuggestionBGColor:    DefaultColor,
			suggestionTextColor:         White,
			suggestionBGColor:           Cyan,
			selectedSuggestionTextColor: Black,
			selectedSuggestionBGColor:   Turquoise,
			maxCompletions:              10,
		},
		buf:       NewBuffer(),
		executor:  executor,
		completer: completer,
		chosen: -1,
	}

	for _, opt := range opts {
		if err := opt(pt); err != nil {
			panic(err)
		}
	}
	return pt
}
