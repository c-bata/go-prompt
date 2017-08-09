# go-prompt

![demo](./_resources/demo.gif)

Library for building powerful interactive prompt in Go, inspired by python-prompt-toolkit.

#### Similar Projects

* [jonathanslenders/python-prompt-toolkit](https://github.com/jonathanslenders/python-prompt-toolkit): **go-prompt** is inspired by this library.
* [peterh/liner](https://github.com/peterh/liner): The most similar project in golang is **liner** that I've ever seen.

#### Projects using go-prompt

* [kube-prompt : An interactive kubernetes client featuring auto-complete written in Go.](https://github.com/c-bata/kube-prompt)

## Getting Started

```go
package main

import (
	"fmt"

	"github.com/c-bata/go-prompt"
)

// executor executes command and print the output.
func executor(in string) {
	fmt.Println("Your input: " + in)
}

// completer returns the completion items from user input.
func completer(in string) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "users", Description: "user table"},
		{Text: "sites", Description: "sites table"},
		{Text: "articles", Description: "articles table"},
		{Text: "comments", Description: "comments table"},
	}
	return prompt.FilterHasPrefix(s, in, true)
}

func main() {
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(">>> "),
		prompt.OptionTitle("sql-prompt"),
	)
	p.Run()
}
```

## Color Options

go-prompt has many color options.
It is difficult to describe by text. So please see below:

![options](./_resources/prompt-options.png)

* **OptionPrefixTextColor(prompt.Color)** : default `prompt.Blue`
* **OptionPrefixBackgroundColor(prompt.Color)** : default `prompt.DefaultColor`
* **OptionInputTextColor(prompt.Color)** : default `prompt.DefaultColor`
* **OptionInputBGColor(prompt.Color)** : default `prompt.DefaultColor`
* **OptionPreviewSuggestionTextColor(prompt.Color)** : default `prompt.Green`
* **OptionPreviewSuggestionBGColor(prompt.Color)** : default `prompt.DefaultColor`
* **OptionSuggestionTextColor(prompt.Color)** : default `prompt.White`
* **OptionSuggestionBGColor(prompt.Color)** : default `prompt.Cyan`
* **OptionSelectedSuggestionTextColor(prompt.Color)** : `default prompt.Black`
* **OptionSelectedSuggestionBGColor(prompt.Color)** : `default prompt.DefaultColor`
* **OptionDescriptionTextColor(prompt.Color)** : default `prompt.Black`
* **OptionDescriptionBGColor(prompt.Color)** : default `prompt.Turquoise`
* **OptionSelectedDescriptionTextColor(prompt.Color)** : default `prompt.White`
* **OptionSelectedDescriptionBGColor(prompt.Color)** : default `prompt.Cyan`

## Other Options

#### `OptionTitle(string)` : default `""`
Option to set title displayed at the header bar of terminal.

#### `OptionHistory([]string)` : default `[]string{}`
Option to set history.

#### `OptionPrefix(string)` : default `"> "`
Option to set prefix string.

#### `OptionMaxSuggestions(x uint16)` : default `6`
The max number of displayed suggestions.

#### `OptionParser(prompt.ConsoleParser)` : default `VT100Parser`
To set a custom ConsoleParser object.
An argument should implement ConsoleParser interface.

#### `OptionWriter(prompt.ConsoleWriter)` : default `VT100Writer`
To set a custom ConsoleWriter object.
An argument should implement ConsoleWriter interace.


## `go-prompt` internals

This is a short description of go-prompt implementation.
go-prompt consists of three parts.

1. Input parser
2. Emulate user input with Buffer object.
3. Render buffer object.

### Input Parser

![input-parser animation](./_resources/input-parser.gif)

Input Parser only supports only vt100 compatible console now.

* Set raw mode.
* Read standard input.
* Parse byte array

### Emulate user input

go-prompt contains Buffer class.
It represents input state by handling user input key.

`Buffer` object has text and cursor position.

```go
package main

import "github.com/c-bata/go-prompt"

func main() {
  b := prompt.NewBuffer()
  ... wip
}
```

### Render Buffer object.

`Renderer` object renders a buffer object.

```go
package main
```

the output is below:

**TODO this is screen shot**

## LICENSE

This software is licensed under the MIT License (See [LICENSE](./LICENSE) ).
