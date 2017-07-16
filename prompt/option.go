package prompt

import "syscall"

type option func(prompt *Prompt) error

func ParserOption(x ConsoleParser) option {
	return func(p *Prompt) error {
		p.in = x
		return nil
	}
}

func WriterOption(x ConsoleWriter) option {
	return func(p *Prompt) error {
		p.renderer.out = x
		return nil
	}
}

func TitleOption(x string) option {
	return func(p *Prompt) error {
		p.renderer.title = x
		return nil
	}
}

func PrefixOption(x string) option {
	return func(p *Prompt) error {
		p.renderer.prefix = x
		return nil
	}
}

func PrefixColorOption(x Color) option {
	return func(p *Prompt) error {
		p.renderer.prefixColor = x
		return nil
	}
}

func CompletionTextColor(x Color) option {
	return func(p *Prompt) error {
		p.renderer.completionTextColor = x
		return nil
	}
}

func CompletionBackgroundColor(x Color) option {
	return func(p *Prompt) error {
		p.renderer.completionBGColor = x
		return nil
	}
}

func SelectedCompletionTextColor(x Color) option {
	return func(p *Prompt) error {
		p.renderer.selectedCompletionTextColor = x
		return nil
	}
}

func SelectedCompletionBackgroundColor(x Color) option {
	return func(p *Prompt) error {
		p.renderer.selectedCompletionBGColor = x
		return nil
	}
}

func MaxCompletionsOption(x uint16) option {
	return func(p *Prompt) error {
		p.renderer.maxCompletions = x
		return nil
	}
}

func NewPrompt(executor Executor, completer Completer, opts ...option) *Prompt {
	pt := &Prompt{
		in: &VT100Parser{fd: syscall.Stdin},
		renderer: &Render{
			prefix:      "> ",
			out:         &VT100Writer{fd: syscall.Stdout},
			prefixColor: Green,
			completionTextColor: White,
			completionBGColor: Cyan,
			selectedCompletionTextColor: Black,
			selectedCompletionBGColor: Turquoise,
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
