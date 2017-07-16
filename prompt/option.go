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
		p.title = x
		return nil
	}
}

func PrefixOption(x string) option {
	return func(p *Prompt) error {
		p.renderer.Prefix = x
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
			Prefix:         "> ",
			out:            &VT100Writer{fd: syscall.Stdout},
		},
		title:     "",
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
