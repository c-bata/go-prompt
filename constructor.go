package prompt

import (
	"fmt"
	"os"
)

// Option is the type to replace default parameters.
// prompt.New accepts any number of options (this is functional option pattern).
type Option func(prompt *Prompt) error

// Callback function that returns a prompt prefix.
type PrefixCallback func() (prefix string)

// Callback function that returns a prompt prefix.
type ExtendedPrefixCallback func() (prefix *Prefix)

type Prefix struct {
	Content         string
	BackgroundColor Color
	TextColor       Color
}

const DefaultIndentSize = 2

// WithIndentSize is an option that sets the amount of spaces
// that constitute a single indentation level.
func WithIndentSize(i int) Option {
	return func(p *Prompt) error {
		p.renderer.indentSize = i
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

type InterruptCallback func(prompt *Prompt, code os.Signal)

func WithInterruptCallback(fn InterruptCallback) Option {
	return func(p *Prompt) error {
		p.interruptCallback = fn
		return nil
	}
}

// WithExecuteOnEnterCallback can be used to set
// a custom callback function that determines whether an Enter key
// should trigger the Executor or add a newline to the user input buffer.
func WithExecuteOnEnterCallback(fn ExecuteOnEnterCallback) Option {
	return func(p *Prompt) error {
		p.executeOnEnterCallback = fn
		return nil
	}
}

// WithCompleter is an option that sets a custom Completer object.
func WithCompleter(c Completer) Option {
	return func(p *Prompt) error {
		p.completion.completer = c
		return nil
	}
}

// WithReader can be used to set a custom Reader object.
func WithReader(r Reader) Option {
	return func(p *Prompt) error {
		p.reader = r
		return nil
	}
}

// WithWriter can be used to set a custom Writer object.
func WithWriter(w Writer) Option {
	return func(p *Prompt) error {
		registerWriter(w)
		p.renderer.out = w
		return nil
	}
}

// WithTitle can be used to set the title displayed at the header bar of the terminal.
func WithTitle(t string) Option {
	return func(p *Prompt) error {
		p.renderer.title = t
		return nil
	}
}

// WithPrefix can be used to set a prefix string for the prompt.
func WithPrefix(prefix string) Option {
	return func(p *Prompt) error {
		p.renderer.prefixCallback = func() string { return prefix }
		return nil
	}
}

// WithInitialText can be used to set the initial buffer text.
func WithInitialText(text string) Option {
	return func(p *Prompt) error {
		p.buffer.InsertTextMoveCursor(text, p.renderer.col, int(p.renderer.row), true)
		return nil
	}
}

// WithCompletionWordSeparator can be used to set word separators.
// An empty separator splits on identifier and punctuation parts.
func WithCompletionWordSeparator(sep string) Option {
	return func(p *Prompt) error {
		p.completion.wordSeparator = sep
		return nil
	}
}

// WithPrefixCallback can be used to change the prefix dynamically by a callback function.
func WithPrefixCallback(f PrefixCallback) Option {
	return func(p *Prompt) error {
		p.renderer.prefixCallback = f
		return nil
	}
}

// WithExtendedPrefixCallback can be used to change the prefix dynamically by a callback function
// with more options like colors etc.
func WithExtendedPrefixCallback(f ExtendedPrefixCallback) Option {
	return func(p *Prompt) error {
		p.renderer.extendedPrefixCallback = f
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
		if history, ok := p.history.(*History); ok {
			history.histories = x
			history.Clear()
			return nil
		}

		return fmt.Errorf("cannot use WithHistory on a custom HistoryInterface implementation")
	}
}

// WithHistorySize sets the maximum number of history entries
func WithHistorySize(size int) Option {
	return func(p *Prompt) error {
		history, ok := p.history.(*History)
		if !ok {
			return fmt.Errorf("cannot use WithHistorySize on a custom HistoryInterface implementation")
		}
		if size < 0 {
			return fmt.Errorf("history size should be greater than or equal to 0, but got %d", size)
		}

		history.size = size
		return nil
	}
}

// WithFSHistory stores history in the file at path.
func WithFSHistory(path string) Option {
	return func(p *Prompt) error {
		history, err := NewFSHistory(path)
		if err != nil {
			return err
		}
		p.history = history
		return nil
	}
}

// WithCustomHistory to set use of own history implementation.
func WithCustomHistory(history HistoryInterface) Option {
	return func(p *Prompt) error {
		p.history = history
		p.history.Clear()
		return nil
	}
}

// WithKeyBindMode set a key bind mode.
func WithKeyBindMode(m KeyBindMode) Option {
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

// WithKeyBind to set a custom key bind.
func WithKeyBind(b ...KeyBind) Option {
	return func(p *Prompt) error {
		p.keyBindings = append(p.keyBindings, b...)
		return nil
	}
}

// WithASCIICodeBind to set a custom key bind.
func WithASCIICodeBind(b ...ASCIICodeBind) Option {
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

func DefaultExecuteOnEnterCallback(p *Prompt, indentSize int) (int, bool) {
	return 0, true
}

func DefaultPrefixCallback() string {
	return "> "
}

// New returns a Prompt with powerful auto-completion.
func New(executor Executor, opts ...Option) *Prompt {
	pt := &Prompt{
		reader:                 NewStdinReader(),
		renderer:               NewRenderer(),
		buffer:                 NewBuffer(),
		executor:               executor,
		history:                NewHistory(),
		completion:             NewCompletionManager(6),
		executeOnEnterCallback: DefaultExecuteOnEnterCallback,
		keyBindMode:            EmacsKeyBind, // All the above assume that bash is running in the default Emacs setting
	}

	for _, opt := range opts {
		if err := opt(pt); err != nil {
			panic(err)
		}
	}
	return pt
}
