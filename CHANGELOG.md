# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

### Changed
- Improved word positioning. Punctuation is not treated as a part of a word.

### Added
- `(p *Prompt) CursorLeftBytes(count strings.ByteNumber) bool` - moves the cursor the given amount of bytes to the left
- `(p *Prompt) CursorRightBytes(count strings.ByteNumber) bool`  - moves the cursor the given amount of bytes to the right
- `CustomOverrideCommonKeyBind` - key bind mode that makes custom key bindings override common ones
- `CustomOverrideEmacsKeyBind` - key bind mode that makes custom key bindings override emacs and common ones

## [1.4.0] - 12.05.2026

[Diff](https://github.com/elk-language/go-prompt/compare/v1.3.0...elk-language:go-prompt:v1.4.0)

### Added
- styled prefixes - [Issue#27](https://github.com/elk-language/go-prompt/issues/27)
  - `Prefix` a struct that represents a styled prefix (background and text color)
  - `type ExtendedPrefixCallback func() (prefix *Prefix)` - a callback type that returns a styled prefix
  - `func WithExtendedPrefixCallback(f ExtendedPrefixCallback) Option` - a new option that registers a callback that creates a dynamic prefix with a custom background and text color
- custom interrupt handlers - [Issue#30](https://github.com/elk-language/go-prompt/issues/30)
  - `type InterruptCallback func(prompt *Prompt, code os.Signal)` - a callback type that handles interrupts
  - `func WithInterruptCallback(fn InterruptCallback) Option` - a new option that registers a callback that captures and handles interrupts
### Removed
- Sleep between reads - [Issue#29](https://github.com/elk-language/go-prompt/issues/29)

## [1.3.0] - 21.05.2025

[Diff](https://github.com/elk-language/go-prompt/compare/v1.2.0...elk-language:go-prompt:v1.3.0)

### Added
- [PR#26](https://github.com/elk-language/go-prompt/pull/26) - add `HistoryInterface` to let users implement their own custom history
- Fix multiline editing
- Add maximum size to history
- Add new methods to History
  - `Get`
  - `Entries`
  - `DeleteAll`

## [1.2.0] - 28.03.2025

[Diff](https://github.com/elk-language/go-prompt/compare/v1.1.5...elk-language:go-prompt:v1.2.0)

### Added
- [PR#24](https://github.com/elk-language/go-prompt/pull/24) - the project now builds for other environments than windows or unix (WASM)
- Update Go to 1.24

## [1.1.5] - 15.08.2023

[Diff](https://github.com/elk-language/go-prompt/compare/v1.1.4...elk-language:go-prompt:v1.1.5)

### Fixed
- [PR#18](https://github.com/elk-language/go-prompt/pull/18) - fix the `http-prompt` example and refactor the example building script (make it fail when an example doesn't compile)
- [PR#17](https://github.com/elk-language/go-prompt/pull/17) - correctly restore terminal parameters
- [PR#16](https://github.com/elk-language/go-prompt/pull/16) - propagate `getOriginalTermios` error when called multiple times

## [1.1.4] - 29.07.2023

[Diff](https://github.com/elk-language/go-prompt/compare/v1.1.3...elk-language:go-prompt:v1.1.4)

### Added
- `func (*prompt.Document).PreviousLineIndentSpaces() int`
- `func (*prompt.Document).PreviousLineIndentLevel(indentSize int) int`
- `func (*prompt.Document).PreviousLine() (s string, ok bool)`

## [1.1.3] - 29.07.2023

[Diff](https://github.com/elk-language/go-prompt/compare/v1.1.2...elk-language:go-prompt:v1.1.3)

### Added
- `func (*prompt.Document).IndentSpaces(input string) int`
- `func (*prompt.Document).IndentLevel(input string, indentSize int) int`
- `func (*prompt.Document).CurrentLineIndentSpaces() int`
- `func (*prompt.Document).CurrentLineIndentLevel(indentSize int) int`

## [1.1.2] - 29.07.2023

[Diff](https://github.com/elk-language/go-prompt/compare/v1.1.1...elk-language:go-prompt:v1.1.2)

### Added
- `func (*prompt.Prompt).DeleteBeforeCursor(count strings.GraphemeNumber) string`
- `func (*prompt.Prompt).DeleteBeforeCursorRunes(count strings.RuneNumber) string`
- `func (*prompt.Prompt).Delete(count strings.GraphemeNumber) string`
- `func (*prompt.Prompt).DeleteRunes(count strings.RuneNumber) string`
- `func (*prompt.Prompt).InsertText(text string, overwrite bool)`
- `func (*prompt.Prompt).InsertTextMoveCursor(text string, overwrite bool)`
- `func (*prompt.Prompt).UserInputColumns() strings.Width`

## [1.1.1] - 28.07.2023

[Diff](https://github.com/elk-language/go-prompt/compare/v1.1.0...elk-language:go-prompt:v1.1.1)

### Added
- `func (*prompt.Prompt).TerminalColumns() strings.Width`
- `func (*prompt.Prompt).TerminalRows() int`

### Changed
- Change signatures:
  - `prompt.ExecuteOnEnterCallback`
    - from `func(buffer *prompt.Buffer, indentSize int) (indent int, execute bool)`
    - to `func(p *prompt.Prompt, indentSize int) (indent int, execute bool)`
- Change `(*Prompt).Buffer` from a public field to a public getter method


## [1.1.0] - 28.07.2023

[Diff](https://github.com/elk-language/go-prompt/compare/v1.0.3...elk-language:go-prompt:v1.1.0)

### Fixed
- Fix cursor movement for text with grapheme clusters like 🇵🇱, 🙆🏿‍♂️

### Added
- Add `strings.GraphemeNumber`, a type that represents the amount of grapheme clusters in a string (or an offset of a grapheme cluster in a string)
  - `type strings.GraphemeNumber int`
- `func strings.GraphemeCountInString(text string) strings.GraphemeNumber`
- `func strings.RuneCountInString(s string) strings.RuneNumber`
- `func strings.RuneIndexNthGrapheme(text string, n strings.GraphemeNumber) strings.RuneNumber`
- `func strings.RuneIndexNthColumn(text string, n strings.Width) strings.RuneNumber`
- `func (*prompt.Document).GetCursorLeftPositionRunes(count strings.RuneNumber) strings.RuneNumber`
- `func (*prompt.Document).GetCursorRightPositionRunes(count strings.RuneNumber) strings.RuneNumber`
- `func (*prompt.Document).LastLineIndentLevel(indentSize int) int`
- `func (*prompt.Document).LastLineIndentSpaces() int`
- `func (*prompt.Buffer).DeleteRunes(count strings.RuneNumber, col strings.Width, row int) string`
- `func (*prompt.Buffer).DeleteBeforeCursorRunes(count strings.RuneNumber, col strings.Width, row int) string`

### Changed
- Change signatures:
  - `strings.RuneCount`
    - from `func strings.RuneCount(s string) strings.RuneNumber`
    - to `func strings.RuneCount(b []byte) strings.RuneNumber`
  - `prompt.ExecuteOnEnterCallback`
    - from `func(input string, indentSize int) (indent int, execute bool)`
    - to `func(buffer *prompt.Buffer, indentSize int) (indent int, execute bool)`
  - `(*prompt.Document).CursorPositionCol`
    - from `func (*prompt.Document).CursorPositionCol() (col strings.RuneNumber)`
    - to `func (*prompt.Document).CursorPositionCol() (col strings.Width)`
  - `(*prompt.Document).GetCursorRightPosition`
    - from `func (*prompt.Document).GetCursorRightPosition(count strings.RuneNumber) strings.RuneNumber`
    - to `func (*prompt.Document).GetCursorRightPosition(count strings.Width) strings.RuneNumber`
  - `(*prompt.Document).GetCursorLeftPosition`
    - from `func (*prompt.Document).GetCursorLeftPosition(count strings.RuneNumber) strings.RuneNumber`
    - to `func (*prompt.Document).GetCursorLeftPosition(count strings.Width) strings.RuneNumber`
  - `(*prompt.Document).GetCursorUpPosition`
    - from `func (*prompt.Document).GetCursorUpPosition(count int, preferredColumn strings.RuneNumber) strings.RuneNumber`
    - to `func (*prompt.Document).GetCursorUpPosition(count int, preferredColumn strings.Width) strings.RuneNumber`
  - `(*prompt.Document).GetCursorDownPosition`
    - from `func (*prompt.Document).GetCursorDownPosition(count int, preferredColumn strings.RuneNumber) strings.RuneNumber`
    - to `func (*prompt.Document).GetCursorDownPosition(count int, preferredColumn strings.Width) strings.RuneNumber`
  - `(*prompt.Document).TranslateRowColToIndex`
    - from `func (*prompt.Document).TranslateRowColToIndex(row int, column strings.RuneNumber) strings.RuneNumber`
    - to `func (*prompt.Document).TranslateRowColToIndex(row int, column strings.Width) strings.RuneNumber`
  - `(*prompt.Buffer).Delete`
    - from `func (*Buffer).Delete(count istrings.RuneNumber, col istrings.Width, row int) string`
    - to `func (*Buffer).Delete(count istrings.GraphemeNumber, col istrings.Width, row int) string`
  - `(*prompt.Buffer).DeleteBeforeCursor`
    - from `func (*Buffer).DeleteBeforeCursor(count istrings.RuneNumber, col istrings.Width, row int) string`
    - to `func (*Buffer).DeleteBeforeCursor(count istrings.GraphemeNumber, col istrings.Width, row int) string`

## [1.0.3] - 25.07.2023

[Diff](https://github.com/elk-language/go-prompt/compare/v1.0.2...elk-language:go-prompt:v1.0.3)

### Fixed
- Reset formatting after rendering a styled token
- Fix unit tests

## [1.0.2] - 25.07.2023

[Diff](https://github.com/elk-language/go-prompt/compare/v1.0.1...elk-language:go-prompt:v1.0.2)

### Added

- `prompt.Token` has new methods:
  - `BackgroundColor() prompt.Color` - define the background color for the token
  - `DisplayAttributes() []prompt.DisplayAttribute` - define the font eg. bold, italic, underline
- `prompt.NewSimpleToken` has new options:
  - `prompt.SimpleTokenWithColor(c Color) SimpleTokenOption`
  - `prompt.SimpleTokenWithBackgroundColor(c Color) SimpleTokenOption`
  - `prompt.SimpleTokenWithDisplayAttributes(attrs ...DisplayAttribute) SimpleTokenOption`
- `prompt.Writer` has new methods:
  - `prompt.SetDisplayAttributes(fg, bg Color, attrs ...DisplayAttribute)`

### Changed

- change the signature of `prompt.NewSimpleToken` from `func NewSimpleToken(color Color, firstIndex, lastIndex istrings.ByteNumber) *SimpleToken` to `func NewSimpleToken(firstIndex, lastIndex istrings.ByteNumber, opts ...SimpleTokenOption) *SimpleToken`

## [1.0.1] - 25.07.2023

[Diff](https://github.com/elk-language/go-prompt/compare/v1.0.0...elk-language:go-prompt:v1.0.1)

### Added

- `prompt.Token` has a new method `FirstByteIndex() strings.ByteNumber`


## [1.0.0] - 25.07.2023

[Diff](https://github.com/elk-language/go-prompt/compare/v0.2.6...elk-language:go-prompt:v1.0.0)

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

- examples:
  - `_example/bang-executor` -- a sample program which uses the new `ExecuteOnEnterCallback`. Pressing <kbd>Enter</kbd> will insert a newline unless the input ends with an exclamation point `!` (then it gets executed).
  - `_example/even-lexer` -- a sample program which shows how to use the new lexer feature. It implements a simple lexer which colours every character with an even index green.

### Changed

- Update Go from 1.16 to 1.19
- The cursor can move in 2D (left-right, up-down)
- The Up arrow key will jump to the line above if the cursor is beyond the first line, but it will replace the input with the previous history entry if it's on the first line (like in Ruby's irb)
- The Down arrow key will jump to the line below if the cursor is before the last line, but it will replace the input with the next history entry if it's on the last line (like in Ruby's irb)
- <kbd>Tab</kbd> will insert a single indentation level when there are no suggestions
- <kbd>Shift</kbd> + <kbd>Tab</kbd> will delete a single indentation level when there are no suggestions and the line before the cursor consists only of indentation (spaces)
- Make `Completer` optional when creating a new `prompt.Prompt`. Change the signature of `prompt.New` from `func New(Executor, Completer, ...Option) *Prompt` to `func New(Executor, ...Option) *Prompt`
- Make `prefix` and `completer` optional in `prompt.Input`. Change the signature of `prompt.Input` from `func Input(string, Completer, ...Option) string` to `func Input(...Option) string`.
- Rename `prompt.ConsoleParser` to `prompt.Reader` and make it embed `io.ReadCloser`
- Rename `prompt.ConsoleWriter` to `prompt.Writer` and make it embed `io.Writer` and `io.StringWriter`
- Rename `prompt.Render` to `prompt.Renderer`
- Rename `prompt.OptionTitle` to `prompt.WithTitle`
- Rename `prompt.OptionPrefix` to `prompt.WithPrefix`
- Rename `prompt.OptionInitialBufferText` to `prompt.WithInitialText`
- Rename `prompt.OptionCompletionWordSeparator` to `prompt.WithCompletionWordSeparator`
- Replace `prompt.OptionLivePrefix` with `prompt.WithPrefixCallback` -- `func() string`. The prefix is always determined by a callback function which should always return a `string`.
- Rename `prompt.OptionPrefixTextColor` to `prompt.WithPrefixTextColor`
- Rename `prompt.OptionPrefixBackgroundColor` to `prompt.WithPrefixBackgroundColor`
- Rename `prompt.OptionInputTextColor` to `prompt.WithInputTextColor`
- Rename `prompt.OptionInputBGColor` to `prompt.WithInputBGColor`
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
- Change the signature of `Completer` from `func(Document) []Suggest` to `func(Document) (suggestions []Suggest, startChar, endChar istrings.RuneNumber)`
- Change the signature of `KeyBindFunc` from `func(*Buffer)` to `func(p *Prompt) (rerender bool)`

### Fixed

- Make pasting multiline text work properly
- Make pasting text with tabs work properly (tabs get replaced with indentation -- spaces)
- Introduce `strings.ByteNumber`, `strings.RuneNumber`, `strings.Width` to reduce the ambiguity of when to use which of the three main units used by this library to measure string length and index parts of strings. Several subtle bugs (using the wrong unit) causing panics have been fixed this way.
- Remove a `/dev/tty` leak in `prompt.PosixReader` (old `prompt.PosixParser`)

### Removed

- `prompt.SwitchKeyBindMode`
- `prompt.OptionPreviewSuggestionTextColor`
- `prompt.OptionPreviewSuggestionBGColor`

## [0.2.6] - 2021-03-03

[Diff](https://github.com/elk-language/go-prompt/compare/v0.2.5...elk-language:go-prompt:v0.2.6)

### Changed

- Update pkg/term to 1.2.0


## [0.2.5] - 2020-09-19

[Diff](https://github.com/elk-language/go-prompt/compare/v0.2.4...elk-language:go-prompt:v0.2.5)

### Changed

- Upgrade all dependencies to latest


## [0.2.4] - 2020-09-18

[Diff](https://github.com/elk-language/go-prompt/compare/v0.2.3...elk-language:go-prompt:v0.2.4)

### Changed

- Update pkg/term module to latest and use unix.Termios


## [0.2.3] - 2018-10-25

[Diff](https://github.com/elk-language/go-prompt/compare/v0.2.2...elk-language:go-prompt:v0.2.3)

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

[Diff](https://github.com/elk-language/go-prompt/compare/v0.2.1...elk-language:go-prompt:v0.2.2)

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

[Diff](https://github.com/elk-language/go-prompt/compare/v0.2.0...elk-language:go-prompt:v0.2.1)

### Added

* ~~It seems that windows support is almost perfect.~~
    * A critical bug is found :( When you change a terminal window size, the layout will be broken because current implementation cannot catch signal for updating window size on Windows.

### Fixed

* <kbd>Shift</kbd> + <kbd>Tab</kbd> handling on Windows.
* 4-dimension arrow keys handling on Windows.


## [0.2.0] - 2018-02-13

[Diff](https://github.com/elk-language/go-prompt/compare/v0.1.0...elk-language:go-prompt:v0.2.0)

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
