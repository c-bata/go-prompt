# Getting started

#### Download

```
$ go get -u github.com/c-bata/go-prompt
```

#### Usage

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

![options](../_resources/prompt-options.png)

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

