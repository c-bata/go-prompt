# go-prompt-toolkit

![demo](./_resources/demo.gif)

Library for building powerful interactive command lines in Golang.

#### Similar Projects

* [jonathanslenders/python-prompt-toolkit](https://github.com/jonathanslenders/python-prompt-toolkit): **go-prompt-toolkit** is inspired by this library.
* [peterh/liner](https://github.com/peterh/liner): The most similar project in golang is **liner** that I've ever seen.

#### Projects using go-prompt-toolkit

* [kube-prompt : An interactive kubernetes client featuring autocomplete using prompt-toolkit.](https://github.com/c-bata/kube-prompt)

## Getting Started

```go
package main

import (
    "context"
    "fmt"

    "github.com/c-bata/go-prompt-toolkit/prompt"
)

// executor executes command and return the output string.
// 1. Execute sql
// 2. Get response and return output
func executor(ctx context.Context, sql string) string {
    res := "something response from db."
    return res // this is printed in console.
}

// completer returns the completion items from user input.
func completer(sql string) []string {
    return []string{"users", "articles", "comments", "groups", "tags"}
}

func main() {
    pt := prompt.NewPrompt(
        executor,
        completer,
        prompt.OptionTitle("sqlite3-prompt"),
        prompt.OptionPrefix(">>> "),
        prompt.OptionPrefixColor(prompt.Blue),
    )
    defer fmt.Println("\nGoodbye!")
    pt.Run()
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
#### `OptionOutputTextColor(x Color)`
#### `OptionOutputBGColor(x Color)`
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
