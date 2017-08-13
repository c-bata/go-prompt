# go-prompt

Library for building a powerful interactive prompt, inspired by python-prompt-toolkit.
Easy building a multi-platform binary of the command line tools because built with Golang.

#### Similar Projects

* [jonathanslenders/python-prompt-toolkit](https://github.com/jonathanslenders/python-prompt-toolkit): go-prompt is inspired by this library.
* [peterh/liner](https://github.com/peterh/liner): The most similar project in golang is **liner** that I've ever seen.

#### Projects using go-prompt

* [kube-prompt : An interactive kubernetes client featuring auto-complete written in Go.](https://github.com/c-bata/kube-prompt)

## Features

#### Powerful auto completion

![demo](./_resources/kube-prompt.gif)

(This is a GIF animation of kube-prompt.)

#### Keyboard Shortcuts

![Keyboard shortcuts](./_resources/keyboard-shortcuts.gif)

You can customize keyboard shortcuts. More details are available from 'KeyBoard Shortcuts' section in Developer Guide.

#### Easy to use

Usage is like this:

```go
package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
)

func completer(buf prompt.Buffer) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "users", Description: "user table"},
		{Text: "sites", Description: "sites table"},
		{Text: "articles", Description: "articles table"},
	}
	return prompt.FilterHasPrefix(s, buf.Text(), true)
}

func main() {
	fmt.Println("Please select table.")
	t := prompt.Input("> ", completer)
	fmt.Println("You selected " + t)
}
```

More practical example is avairable from `_example` directory or [kube-prompt](https://github.com/c-bata/kube-prompt).

#### Flexible customization

![options](./_resources/prompt-options.png)
go-prompt has many color options. All options are listed in [Developer Guide](./example/README.md).

#### History
**up-arrow** and **down-arrow** to walk through the command line history.

![History](./_resources/history.gif)

## Other Information

* If you want to create projects using go-prompt, you might want to look at the [Getting Started](./example/README.md).
* If you want to contribute go-prompt, you might want to look at the [Developer Guide](./_tools/README.md).
* If you want to know internal API, you might want to look at the [GoDoc](http://godoc.org/github.com/c-bata/go-prompt).


## Author

Masashi Shibata

* Twitter: [@c\_bata\_](https://twitter.com/c_bata_/)
* Github: [@c-bata](https://github.com/c-bata/)

## LICENSE

This software is licensed under the MIT License (See [LICENSE](./LICENSE) ).
