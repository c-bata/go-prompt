# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).


## [1.0.0]

This release contains a major refactoring of the codebase.
It's the first release of the [elk-language/go-prompt](https://github.com/elk-language/go-prompt) fork.

The original library has been abandoned for at least 2 years now (although serious development has stopped 5 years ago).

This release aims to make the code a bit cleaner, fix a couple of bugs and provide new, essential functionality such as syntax highlighting, dynamic <kbd>Enter</kbd> and multiline edit support.

### Added

- `prompt.New` constructor options:
  - `prompt.WithLexer` let's you set a custom lexer for providing syntax highlighting
  - `prompt.WithCompleter` for setting a custom `Completer` (completer is no longer a required argument in `prompt.New`)
  - `prompt.WithIndentSize` let's you customise how many spaces should constitute a single indentation level
  - `prompt.WithExecuteOnEnterCallback`

- `prompt.Position` -- represents the cursor's position in 2D
- `prompt.Lexer`, `prompt.Token`, `prompt.SimpleToken`, `prompt.EagerLexer`, `prompt.LexerFunc` -- new syntax highlighting functionality
- `prompt.ExecuteOnEnterCallback` -- new dynamic <kbd>Enter</kbd> functionality (decide whether to insert a newline and indent or execute the input)

- `_example/bang-executor` -- a sample program which uses the new `ExecuteOnEnterCallback`. Pressing <kbd>Enter</kbd> will insert a newline unless the input ends with an exclamation point `!` (then it gets executed).
- `_example/even-lexer` -- a sample program which shows how to use the new lexer feature. It implements a simple lexer which colours every character with an even index green.

### Changed

- Update Go from 1.16 to 1.19
- The cursor can move in 2D.
- The Up arrow key will jump to the line above if the cursor is beyond the first line, but it will replace the input with the previous history entry if it's on the first line (like in Ruby's irb)
- The Down arrow key will jump to the line below if the cursor is before the last line, but it will replace the input with the next history entry if it's on the last line (like in Ruby's irb)
- Make `Completer` optional when creating a new `prompt.Prompt`. Change the signature of `prompt.New` from `func New(Executor, Completer, ...Option) *Prompt` to `func New(Executor, ...Option) *Prompt`
- Rename `prompt.ConsoleParser` to `prompt.Reader` and make it embed `io.ReadCloser`
- Rename `prompt.ConsoleWriter` to `prompt.Writer` and make it embed `io.Writer` and `io.StringWriter`
- Rename `prompt.OptionTitle` to `prompt.WithTitle`
- Rename `prompt.OptionPrefix` to `prompt.WithPrefix`
- Rename `prompt.OptionInitialBufferText` to `prompt.WithInitialText`
- Rename `prompt.OptionCompletionWordSeparator` to `prompt.WithCompletionWordSeparator`
- Replace `prompt.OptionLivePrefix` with `prompt.WithPrefixCallback` -- `func() string`. The prefix is always determined by a callback function which should always return a `string`.
- Rename `prompt.OptionPrefixTextColor` to `prompt.WithPrefixTextColor`
- Rename `prompt.OptionPrefixBackgroundColor` to `prompt.WithPrefixBackgroundColor`
- Rename `prompt.OptionInputTextColor` to `prompt.WithInputTextColor`
- Rename `prompt.OptionInputBGColor` to `prompt.WithInputBGColor`
- Rename `prompt.OptionPreviewSuggestionTextColor` to `prompt.WithPreviewSuggestionTextColor`
- Rename `prompt.OptionSuggestionTextColor` to `prompt.WithSuggestionTextColor`
- Rename `prompt.OptionSuggestionBGColor` to `prompt.WithSuggestionBGColor`
- Rename `prompt.OptionSelectedSuggestionTextColor` to `prompt.WithSelectedSuggestionTextColor`
- Rename `prompt.OptionSelectedSuggestionBGColor` to `prompt.WithSelectedSuggestionBGColor`
- Rename `prompt.OptionDescriptionTextColor` to `prompt.WithDescriptionTextColor`
- Rename `prompt.OptionDescriptionBGColor` to `prompt.WithDescriptionBGColor`
- Rename `prompt.OptionSelectedDescriptionTextColor` to `prompt.WithSelectedDescriptionTextColor`
- Rename `prompt.OptionSelectedDescriptionBGColor` to `prompt.WithSelectedDescriptionBGColor`
- Rename `prompt.OptionScrollbarThumbColor` to `prompt.WithScrollbarThumbColor`
- Rename `prompt.OptionScrollbarBGColor` to `prompt.WithScrollbarBGColor`
- Rename `prompt.OptionMaxSuggestion` to `prompt.WithMaxSuggestion`
- Rename `prompt.OptionHistory` to `prompt.WithHistory`
- Rename `prompt.OptionSwitchKeyBindMode` to `prompt.WithKeyBindMode`
- Rename `prompt.OptionCompletionOnDown` to `prompt.WithCompletionOnDown`
- Rename `prompt.OptionAddKeyBind` to `prompt.WithKeyBind`
- Rename `prompt.OptionAddASCIICodeBind` to `prompt.WithASCIICodeBind`
- Rename `prompt.OptionShowCompletionAtStart` to `prompt.WithShowCompletionAtStart`
- Rename `prompt.OptionBreakLineCallback` to `prompt.WithBreakLineCallback`
- Rename `prompt.OptionExitChecker` to `prompt.WithExitChecker`

### Fixed

- Make pasting multiline text work properly
- Make pasting text with tabs work properly (tabs get replaced with spaces)
- Introduce `strings.ByteNumber`, `strings.RuneNumber`, `strings.StringWidth` to reduce the ambiguity of when to use which of the three main units used by this library to measure string length and index parts of strings. Several subtle bugs (using the wrong unit) causing panics have been fixed this way.
- Remove a `/dev/tty` leak in `PosixReader` (old `PosixParser`)

### Removed

- `prompt.SwitchKeyBindMode`

## [0.2.6] - 2021-03-03

### Changed

- Update pkg/term to 1.2.0


## [0.2.5] - 2020-09-19

### Changed

- Upgrade all dependencies to latest


## [0.2.4] - 2020-09-18

### Changed

- Update pkg/term module to latest and use unix.Termios


## [0.2.3] - 2018-10-25

### Added

* `prompt.FuzzyFilter` for fuzzy matching at [#92](https://github.com/c-bata/go-prompt/pull/92).
* `OptionShowCompletionAtStart` to show completion at start at [#100](https://github.com/c-bata/go-prompt/pull/100).
* `prompt.NewStderrWriter` at [#102](https://github.com/c-bata/go-prompt/pull/102).

### Fixed

* reset display attributes (please see [pull #104](https://github.com/c-bata/go-prompt/pull/104) for more details).
* handle errors of Flush function in ConsoleWriter (please see [pull #97](https://github.com/c-bata/go-prompt/pull/97) for more details).
* don't panic problem when reading from stdin before starting the prompt (please see [issue #88](https://github.com/c-bata/go-prompt/issues/88) for more details).

### Deprecated

* `prompt.NewStandardOutputWriter` -- please use `prompt.NewStdoutWriter`.


## [0.2.2] - 2018-06-28

### Added

* Support CJK (Chinese, Japanese and Korean) and Cyrillic characters.
* `OptionCompletionWordSeparator(x string)` to customize insertion points for completions.
    * To support this, text query functions by arbitrary word separator are added in `Document` (please see [here](https://github.com/c-bata/go-prompt/pull/79) for more details).
* `FilePathCompleter` to complete file path on your system.
* `option` to customize ascii code key bindings.
* `GetWordAfterCursor` method in `Document`.

### Deprecated

* `prompt.Choose` shortcut function is deprecated.


## [0.2.1] - 2018-02-14

### Added

* ~~It seems that windows support is almost perfect.~~
    * A critical bug is found :( When you change a terminal window size, the layout will be broken because current implementation cannot catch signal for updating window size on Windows.

### Fixed

* <kbd>Shift</kbd> + <kbd>Tab</kbd> handling on Windows.
* 4-dimension arrow keys handling on Windows.


## [0.2.0] - 2018-02-13

### Added

* Support scrollbar when there are too many matched suggestions
* Support Windows (but please caution because this is still not perfect).
* `OptionLivePrefix` to update the prefix dynamically
* Clear screen by <kbd>Ctrl</kbd> + <kbd>L</kbd>.

### Fixed

* Improve the <kbd>Ctrl</kbd> + <kbd>W</kbd> keybind.
* Don't panic because when running in a docker container (please see [here](https://github.com/c-bata/go-prompt/pull/32) for details).
* Don't panic when making terminal window small size after input 2 lines of texts. See [here](https://github.com/c-bata/go-prompt/issues/37) for details).
* Get rid of many bugs that layout is broken when using Terminal.app, GNU Terminal and a Goland(IntelliJ).


## [0.1.0] - 2017-08-15

Initial Release
