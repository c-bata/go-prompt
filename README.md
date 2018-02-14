# go-prompt

Library for building a powerful interactive prompt, inspired by [python-prompt-toolkit](https://github.com/jonathanslenders/python-prompt-toolkit).
Easy building a multi-platform binary of the command line tools because written in Golang.

```go
package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
)

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "users", Description: "Store the username and age"},
		{Text: "articles", Description: "Store the article text posted by user"},
		{Text: "comments", Description: "Store the text commented to articles"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func main() {
	fmt.Println("Please select table.")
	t := prompt.Input("> ", completer)
	fmt.Println("You selected " + t)
}
```


#### Projects using go-prompt

* [c-bata/kube-prompt : An interactive kubernetes client featuring auto-complete written in Go.](https://github.com/c-bata/kube-prompt)
* [rancher/cli : The Rancher Command Line Interface (CLI)is a unified tool to manage your Rancher server](https://github.com/rancher/cli)
* [kris-nova/kubicorn : Simple. Cloud Native. Kubernetes. Infrastructure.](https://github.com/kris-nova/kubicorn)
* [cch123/asm-cli : Interactive shell of assembly language(X86/X64) based on unicorn and rasm2](https://github.com/cch123/asm-cli)
* [ktr0731/evans : more expressive universal gRPC client](https://github.com/ktr0731/evans)
* (If you create a CLI using go-prompt and want your own project to be listed here, Please submit a Github Issue.)

## Features

### Powerful auto-completion

[![demo](https://github.com/c-bata/assets/raw/master/go-prompt/kube-prompt.gif)](https://github.com/c-bata/kube-prompt)

(This is a GIF animation of kube-prompt.)

### Flexible options

go-prompt provides many options. All options are listed in [Developer Guide](./DEVELOPER_GUIDE.md).

[![options](https://github.com/c-bata/assets/raw/master/go-prompt/prompt-options.png)](#flexible-options)

### Keyboard Shortcuts

Emacs-like keyboard shortcut is available by default (it's also default shortcuts in Bash shell).
You can customize and expand these shortcuts.

[![keyboard shortcuts](https://github.com/c-bata/assets/raw/master/go-prompt/keyboard-shortcuts.gif)](#keyboard-shortcuts)

KeyBinding          | Description
--------------------|---------------------------------------------------------
<kbd>Ctrl + A</kbd> | Go to the beginning of the line (Home)
<kbd>Ctrl + E</kbd> | Go to the End of the line (End)
<kbd>Ctrl + P</kbd> | Previous command (Up arrow)
<kbd>Ctrl + N</kbd> | Next command (Down arrow)
<kbd>Ctrl + F</kbd> | Forward one character
<kbd>Ctrl + B</kbd> | Backward one character
<kbd>Ctrl + D</kbd> | Delete character under the cursor
<kbd>Ctrl + H</kbd> | Delete character before the cursor (Backspace)
<kbd>Ctrl + W</kbd> | Cut the Word before the cursor to the clipboard.
<kbd>Ctrl + K</kbd> | Cut the Line after the cursor to the clipboard.
<kbd>Ctrl + U</kbd> | Cut/delete the Line before the cursor to the clipboard.
<kbd>Ctrl + L</kbd> | Clear the screen

### History

You can use up-arrow and down-arrow to walk through the history of commands executed.

[![History](https://github.com/c-bata/assets/raw/master/go-prompt/history.gif)](#history)


### Multiple platform support

We confirmed following terminals

* iTerm2 (macOS)
* Terminal.app (macOS)
* Command Prompt (Windows)
* GNU Terminal (Ubuntu)


## Links

* [Developer Guide](./DEVELOPER_GUIDE.md).
* [Change Log](./CHANGELOG.md)
* [GoDoc](http://godoc.org/github.com/c-bata/go-prompt).

## Author

Masashi Shibata

* Twitter: [@c\_bata\_](https://twitter.com/c_bata_/)
* Github: [@c-bata](https://github.com/c-bata/)
* Facebook: [Masashi Shibata](https://www.facebook.com/masashi.cbata)

## LICENSE

This software is licensed under the MIT License (See [LICENSE](./LICENSE) ).

