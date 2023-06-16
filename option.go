package prompt

// Option is the type to replace default parameters.
// prompt.New accepts any number of options (this is functional option pattern).
type Option func(prompt IPrompt) error

// OptionParser to set a custom ConsoleParser object. An argument should implement ConsoleParser interface.
func OptionParser(x ConsoleParser) Option {
	return func(p IPrompt) error {
		p.SetConsoleParser(x)
		return nil
	}
}

// OptionWriter to set a custom ConsoleWriter object. An argument should implement ConsoleWriter interface.
func OptionWriter(x ConsoleWriter) Option {
	return func(p IPrompt) error {
		registerConsoleWriter(x)
		p.Renderer().out = x
		return nil
	}
}

// OptionTitle to set title displayed at the header bar of terminal.
func OptionTitle(x string) Option {
	return func(p IPrompt) error {
		p.Renderer().title = x
		return nil
	}
}

// OptionPrefix to set prefix string.
func OptionPrefix(x string) Option {
	return func(p IPrompt) error {
		p.Renderer().prefix = x
		return nil
	}
}

// OptionInitialBufferText to set the initial buffer text
func OptionInitialBufferText(x string) Option {
	return func(p IPrompt) error {
		p.Buffer().InsertText(x, false, true)
		return nil
	}
}

// OptionCompletionWordSeparator to set word separators. Enable only ' ' if empty.
func OptionCompletionWordSeparator(x string) Option {
	return func(p IPrompt) error {
		p.CompletionManager().wordSeparator = x
		return nil
	}
}

// OptionLivePrefix to change the prefix dynamically by callback function
func OptionLivePrefix(f func() (prefix string, useLivePrefix bool)) Option {
	return func(p IPrompt) error {
		p.Renderer().livePrefixCallback = f
		return nil
	}
}

// OptionPrefixTextColor change a text color of prefix string
func OptionPrefixTextColor(x Color) Option {
	return func(p IPrompt) error {
		p.Renderer().prefixTextColor = x
		return nil
	}
}

// OptionPrefixBackgroundColor to change a background color of prefix string
func OptionPrefixBackgroundColor(x Color) Option {
	return func(p IPrompt) error {
		p.Renderer().prefixBGColor = x
		return nil
	}
}

// OptionInputTextColor to change a color of text which is input by user
func OptionInputTextColor(x Color) Option {
	return func(p IPrompt) error {
		p.Renderer().inputTextColor = x
		return nil
	}
}

// OptionInputBGColor to change a color of background which is input by user
func OptionInputBGColor(x Color) Option {
	return func(p IPrompt) error {
		p.Renderer().inputBGColor = x
		return nil
	}
}

// OptionPreviewSuggestionTextColor to change a text color which is completed
func OptionPreviewSuggestionTextColor(x Color) Option {
	return func(p IPrompt) error {
		p.Renderer().previewSuggestionTextColor = x
		return nil
	}
}

// OptionPreviewSuggestionBGColor to change a background color which is completed
func OptionPreviewSuggestionBGColor(x Color) Option {
	return func(p IPrompt) error {
		p.Renderer().previewSuggestionBGColor = x
		return nil
	}
}

// OptionSuggestionTextColor to change a text color in drop down suggestions.
func OptionSuggestionTextColor(x Color) Option {
	return func(p IPrompt) error {
		p.Renderer().suggestionTextColor = x
		return nil
	}
}

// OptionSuggestionBGColor change a background color in drop down suggestions.
func OptionSuggestionBGColor(x Color) Option {
	return func(p IPrompt) error {
		p.Renderer().suggestionBGColor = x
		return nil
	}
}

// OptionSelectedSuggestionTextColor to change a text color for completed text which is selected inside suggestions drop down box.
func OptionSelectedSuggestionTextColor(x Color) Option {
	return func(p IPrompt) error {
		p.Renderer().selectedSuggestionTextColor = x
		return nil
	}
}

// OptionSelectedSuggestionBGColor to change a background color for completed text which is selected inside suggestions drop down box.
func OptionSelectedSuggestionBGColor(x Color) Option {
	return func(p IPrompt) error {
		p.Renderer().selectedSuggestionBGColor = x
		return nil
	}
}

// OptionDescriptionTextColor to change a background color of description text in drop down suggestions.
func OptionDescriptionTextColor(x Color) Option {
	return func(p IPrompt) error {
		p.Renderer().descriptionTextColor = x
		return nil
	}
}

// OptionDescriptionBGColor to change a background color of description text in drop down suggestions.
func OptionDescriptionBGColor(x Color) Option {
	return func(p IPrompt) error {
		p.Renderer().descriptionBGColor = x
		return nil
	}
}

// OptionSelectedDescriptionTextColor to change a text color of description which is selected inside suggestions drop down box.
func OptionSelectedDescriptionTextColor(x Color) Option {
	return func(p IPrompt) error {
		p.Renderer().selectedDescriptionTextColor = x
		return nil
	}
}

// OptionSelectedDescriptionBGColor to change a background color of description which is selected inside suggestions drop down box.
func OptionSelectedDescriptionBGColor(x Color) Option {
	return func(p IPrompt) error {
		p.Renderer().selectedDescriptionBGColor = x
		return nil
	}
}

// OptionScrollbarThumbColor to change a thumb color on scrollbar.
func OptionScrollbarThumbColor(x Color) Option {
	return func(p IPrompt) error {
		p.Renderer().scrollbarThumbColor = x
		return nil
	}
}

// OptionScrollbarBGColor to change a background color of scrollbar.
func OptionScrollbarBGColor(x Color) Option {
	return func(p IPrompt) error {
		p.Renderer().scrollbarBGColor = x
		return nil
	}
}

// OptionMaxSuggestion specify the max number of displayed suggestions.
func OptionMaxSuggestion(x uint16) Option {
	return func(p IPrompt) error {
		p.CompletionManager().max = x
		return nil
	}
}

// OptionHistory to set history expressed by string array.
func OptionHistory(x []string) Option {
	return func(p IPrompt) error {
		p.History().histories = x
		p.History().Clear()
		return nil
	}
}

// OptionSwitchKeyBindMode set a key bind mode.
func OptionSwitchKeyBindMode(m KeyBindMode) Option {
	return func(p IPrompt) error {
		p.SetKeyBindMode(m)
		return nil
	}
}

// OptionCompletionOnDown allows for Down arrow key to trigger completion.
func OptionCompletionOnDown() Option {
	return func(p IPrompt) error {
		p.SetCompletionOnDown(true)
		return nil
	}
}

// SwitchKeyBindMode to set a key bind mode.
// Deprecated: Please use OptionSwitchKeyBindMode.
var SwitchKeyBindMode = OptionSwitchKeyBindMode

// OptionAddKeyBind to set a custom key bind.
func OptionAddKeyBind(b ...KeyBind) Option {
	return func(p IPrompt) error {
		p.AddKeyBindings(b...)
		return nil
	}
}

// OptionAddASCIICodeBind to set a custom key bind.
func OptionAddASCIICodeBind(b ...ASCIICodeBind) Option {
	return func(p IPrompt) error {
		p.AddASCIICodeBindings(b...)
		return nil
	}
}

// OptionShowCompletionAtStart to set completion window is open at start.
func OptionShowCompletionAtStart() Option {
	return func(p IPrompt) error {
		p.CompletionManager().showAtStart = true
		return nil
	}
}

// OptionBreakLineCallback to run a callback at every break line
func OptionBreakLineCallback(fn func(*Document)) Option {
	return func(p IPrompt) error {
		p.Renderer().breakLineCallback = fn
		return nil
	}
}

// OptionSetExitCheckerOnInput set an exit function which checks if go-prompt exits its Run loop
func OptionSetExitCheckerOnInput(fn ExitChecker) Option {
	return func(p IPrompt) error {
		p.SetExitChecker(fn)
		return nil
	}
}

// OptionSetLexer set lexer function and enable it.
func OptionSetLexer(fn LexerFunc) Option {
	return func(p IPrompt) error {
		p.Lexer().SetLexerFunction(fn)
		return nil
	}
}

// OptionSetStatementTerminator allows you to configure a callback that tells you if statement
// has been terminated and ready to pass to exec in a (potentially) multiline buffer
func OptionSetStatementTerminator(fn StatementTerminatorCb) Option {
	return func(p IPrompt) error {
		p.SetStatementTerminatorCb(fn)
		return nil
	}
}

// New returns a Prompt with powerful auto-completion.
func New(executor Executor, completer Completer, opts ...Option) IPrompt {
	defaultWriter := NewStdoutWriter()
	registerConsoleWriter(defaultWriter)

	pt := &Prompt{
		in: NewStandardInputParser(),
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
		lexer:       NewLexer(),
		completion:  NewCompletionManager(completer, 6),
		keyBindMode: EmacsKeyBind, // All the above assume that bash is running in the default Emacs setting
		statementTerminatorCb: func(lastKeyStroke Key, buffer *Buffer) bool {
			// terminate statement on enter which is either \r or \n, based on OS
			if lastKeyStroke == ControlM || lastKeyStroke == Enter {
				return true
			}
			return false
		},
	}

	for _, opt := range opts {
		if err := opt(pt); err != nil {
			panic(err)
		}
	}
	return pt
}
