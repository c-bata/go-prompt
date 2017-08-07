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


## Options

![options](./_resources/prompt-options.png)

#### `OptionParser(x ConsoleParser)`
#### `OptionWriter(x ConsoleWriter)`
#### `OptionTitle(x string)`
#### `OptionPrefix(x string)`
#### `OptionPrefixTextColor(x Color)`
#### `OptionPrefixBackgroundColor(x Color)`
#### `OptionInputTextColor(x Color)`
#### `OptionInputBGColor(x Color)`
#### `OptionPreviewSuggestionTextColor(x Color)`
#### `OptionPreviewSuggestionBGColor(x Color)`
#### `OptionSuggestionTextColor(x Color)`
#### `OptionSuggestionBGColor(x Color)`
#### `OptionSelectedSuggestionTextColor(x Color)`
#### `OptionSelectedSuggestionBGColor(x Color)`
#### `OptionMaxCompletions(x uint16)`
#### `OptionDescriptionTextColor(x Color)`
#### `OptionDescriptionBGColor(x Color)`
#### `OptionSelectedDescriptionTextColor(x Color)`
#### `OptionSelectedDescriptionBGColor(x Color)`


## LICENSE

```
MIT License

Copyright (c) 2017 Masashi SHIBATA

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
