package prompt

import (
	"strings"

	fcolor "github.com/fatih/color"
)

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

// OptionWriter to set a custom ConsoleWriter object. An argument should implement ConsoleWriter interface.
func OptionWriter(x ConsoleWriter) Option {
	return func(p *Prompt) error {
		registerConsoleWriter(x)
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

// OptionInitialBufferText to set the initial buffer text
func OptionInitialBufferText(x string) Option {
	return func(p *Prompt) error {
		p.buf.InsertText(x, false, true)
		return nil
	}
}

// OptionCompletionWordSeparator to set word separators. Enable only ' ' if empty.
func OptionCompletionWordSeparator(x string) Option {
	return func(p *Prompt) error {
		p.completion.wordSeparator = x
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

// OptionPrefixColor change a text color of prefix string
func OptionPrefixColor(x *fcolor.Color) Option {
	return func(p *Prompt) error {
		p.renderer.prefixColor = x
		return nil
	}
}

// OptionInputColor to change a color of text which is input by user
func OptionInputColor(x *fcolor.Color) Option {
	return func(p *Prompt) error {
		p.renderer.inputColor = x
		return nil
	}
}

// OptionPreviewSuggestionColor to change a text color which is completed
func OptionPreviewSuggestionColor(x *fcolor.Color) Option {
	return func(p *Prompt) error {
		p.renderer.previewSuggestionColor = x
		return nil
	}
}

// OptionSuggestionColor to change a text color in drop down suggestions.
func OptionSuggestionColor(x *fcolor.Color) Option {
	return func(p *Prompt) error {
		p.renderer.suggestionColor = x
		return nil
	}
}

// OptionSelectedSuggestionColor to change a text color for completed text which is selected inside suggestions drop down box.
func OptionSelectedSuggestionColor(x *fcolor.Color) Option {
	return func(p *Prompt) error {
		p.renderer.selectedSuggestionColor = x
		return nil
	}
}

// OptionDescriptionColor to change a background color of description text in drop down suggestions.
func OptionDescriptionColor(x *fcolor.Color) Option {
	return func(p *Prompt) error {
		p.renderer.descriptionColor = x
		return nil
	}
}

// OptionSelectedDescriptionColor to change a text color of description which is selected inside suggestions drop down box.
func OptionSelectedDescriptionColor(x *fcolor.Color) Option {
	return func(p *Prompt) error {
		p.renderer.selectedDescriptionColor = x
		return nil
	}
}

// OptionScrollbarThumbColor to change a thumb color on scrollbar.
func OptionScrollbarThumbColor(x *fcolor.Color) Option {
	return func(p *Prompt) error {
		p.renderer.scrollbarThumbColor = x
		return nil
	}
}

// OptionScrollbarBGColor to change a background color of scrollbar.
func OptionScrollbarBGColor(x *fcolor.Color) Option {
	return func(p *Prompt) error {
		p.renderer.scrollbarColor = x
		return nil
	}
}

// OptionStatusBarColor sets the color of the status bar
func OptionStatusBarColor(color *fcolor.Color) Option {
	return func(p *Prompt) error {
		p.renderer.statusBarColor = color
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

// OptionCompletionOnDown allows for Down arrow key to trigger completion.
func OptionCompletionOnDown() Option {
	return func(p *Prompt) error {
		p.completionOnDown = true
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

// OptionAddASCIICodeBind to set a custom key bind.
func OptionAddASCIICodeBind(b ...ASCIICodeBind) Option {
	return func(p *Prompt) error {
		p.ASCIICodeBindings = append(p.ASCIICodeBindings, b...)
		return nil
	}
}

// OptionShowCompletionAtStart to set completion window is open at start.
func OptionShowCompletionAtStart() Option {
	return func(p *Prompt) error {
		p.completion.showAtStart = true
		return nil
	}
}

// OptionBreakLineCallback to run a callback at every break line
func OptionBreakLineCallback(fn func(*Document)) Option {
	return func(p *Prompt) error {
		p.renderer.breakLineCallback = fn
		return nil
	}
}

// OptionSetExitCheckerOnInput set an exit function which checks if go-prompt exits its Run loop
func OptionSetExitCheckerOnInput(fn ExitChecker) Option {
	return func(p *Prompt) error {
		p.exitChecker = fn
		return nil
	}
}

// OptionStatusBarCallback sets a callback function that returns the status bar message and color to display it in
func OptionStatusBarCallback(cb func(*Buffer, *CompletionManager) (string, bool)) Option {
	return func(p *Prompt) error {
		p.renderer.statusBarCallback = cb
		return nil
	}
}

// OptionKeywordColor sets colors for keywords in displaying
func OptionKeywordColor(kwCol *fcolor.Color) Option {
	return func(p *Prompt) error {
		p.renderer.keywordColor = kwCol
		return nil
	}
}

// OptionKeywords sets the list of words to consider as keywords
func OptionKeywords(kw []string) Option {
	return func(p *Prompt) error {
		p.renderer.keywords = make(map[string]bool)
		for _, w := range kw {
			lw := strings.ToLower(w)
			p.renderer.keywords[lw] = true
		}
		return nil
	}
}

// New returns a Prompt with powerful auto-completion.
func New(executor Executor, completer Completer, opts ...Option) *Prompt {
	defaultWriter := NewStdoutWriter()
	registerConsoleWriter(defaultWriter)
	defStatus := func(buf *Buffer, comp *CompletionManager) (string, bool) {
		return "", false
	}

	pt := &Prompt{
		in: NewStandardInputParser(),
		renderer: &Render{
			prefix:                   "> ",
			out:                      defaultWriter,
			livePrefixCallback:       func() (string, bool) { return "", false },
			statusBarCallback:        defStatus,
			prefixColor:              fcolor.New(fcolor.FgBlue),
			inputColor:               nil,
			previewSuggestionColor:   fcolor.New(fcolor.FgGreen),
			suggestionColor:          fcolor.New(fcolor.FgWhite, fcolor.BgCyan),
			selectedSuggestionColor:  fcolor.New(fcolor.FgBlack, fcolor.BgHiCyan),
			descriptionColor:         fcolor.New(fcolor.FgBlack, fcolor.BgHiCyan),
			selectedDescriptionColor: fcolor.New(fcolor.FgWhite, fcolor.BgCyan),
			scrollbarColor:           fcolor.New(fcolor.FgCyan),
			scrollbarThumbColor:      fcolor.New(fcolor.FgHiBlack),
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
