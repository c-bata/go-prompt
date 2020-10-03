# go-prompt

[![Go Report Card](https://goreportcard.com/badge/github.com/c-bata/go-prompt)](https://goreportcard.com/report/github.com/c-bata/go-prompt)
![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)
[![GoDoc](https://godoc.org/github.com/c-bata/go-prompt?status.svg)](https://godoc.org/github.com/c-bata/go-prompt) 
![tests](https://github.com/c-bata/go-prompt/workflows/tests/badge.svg)

A library for building powerful interactive prompts inspired by [python-prompt-toolkit](https://github.com/jonathanslenders/python-prompt-toolkit),
making it easier to build cross-platform command line tools using Go.

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
* [kubicorn/kubicorn : Simple, cloud native infrastructure for Kubernetes.](https://github.com/kubicorn/kubicorn)
* [cch123/asm-cli : Interactive shell of assembly language(X86/X64) based on unicorn and rasm2](https://github.com/cch123/asm-cli)
* [ktr0731/evans : more expressive universal gRPC client](https://github.com/ktr0731/evans)
* [CrushedPixel/moshpit: A Command-line tool for datamoshing.](https://github.com/CrushedPixel/moshpit)
* [last-ent/testy-go: Testy Go: A tool for easy testing!](https://github.com/last-ent/testy-go)
* [tiagorlampert/CHAOS: a PoC that allow generate payloads and control remote operating systems.](https://github.com/tiagorlampert/CHAOS)
* [abs-lang/abs: ABS is a scripting language that works best on terminal. It tries to combine the elegance of languages such as Python, or Ruby, to the convenience of Bash.](https://github.com/abs-lang/abs)
* [takashabe/btcli: btcli is a CLI client for the Bigtable. Has many read options and auto-completion.](https://github.com/takashabe/btcli)
* [ysn2233/kafka-prompt: An interactive kafka-prompt(kafka-shell) built on existing kafka command client](https://github.com/ysn2233/kafka-prompt)
* [fishi0x01/vsh: HashiCorp Vault interactive shell](https://github.com/fishi0x01/vsh)
* [mstrYoda/docker-shell: A simple interactive prompt for docker](https://github.com/mstrYoda/docker-shell)
* [c-bata/gh-prompt: An interactive GitHub CLI featuring auto-complete.](https://github.com/c-bata/gh-prompt)
* [docker-slim/docker-slim: Don't change anything in your Docker container image and minify it by up to 30x (and for compiled languages even more) making it secure too! (free and open source)](https://github.com/docker-slim/docker-slim)
* [rueyaa332266/ezcron: Ezcron is a CLI tool, helping you deal with cron expression easier.](https://github.com/rueyaa332266/ezcron)
* [qingstor/qsctl: Advanced command line tool for QingStor Object Storage.](https://github.com/qingstor/qsctl)
* [segmentio/topicctl: Tool for declarative management of Kafka topics](https://github.com/segmentio/topicctl)
* [chriswalz/bit: Bit is a modern Git CLI](https://github.com/chriswalz/bit)
* (If you create a CLI utility using go-prompt and want your own project to be listed here, please submit a GitHub issue.)

## Features

### Powerful auto-completion

[![demo](https://github.com/c-bata/assets/raw/master/go-prompt/kube-prompt.gif)](https://github.com/c-bata/kube-prompt)

(This is a GIF animation of kube-prompt.)

### Flexible options

go-prompt provides many options. Please check [option section of GoDoc](https://godoc.org/github.com/c-bata/go-prompt#Option) for more details.

[![options](https://github.com/c-bata/assets/raw/master/go-prompt/prompt-options.png)](#flexible-options)

### Keyboard Shortcuts

Emacs-like keyboard shortcuts are available by default (these also are the default shortcuts in Bash shell).
You can customize and expand these shortcuts.

[![keyboard shortcuts](https://github.com/c-bata/assets/raw/master/go-prompt/keyboard-shortcuts.gif)](#keyboard-shortcuts)

Key Binding          | Description
---------------------|---------------------------------------------------------
<kbd>Ctrl + A</kbd>  | Go to the beginning of the line (Home)
<kbd>Ctrl + E</kbd>  | Go to the end of the line (End)
<kbd>Ctrl + P</kbd>  | Previous command (Up arrow)
<kbd>Ctrl + N</kbd>  | Next command (Down arrow)
<kbd>Ctrl + F</kbd>  | Forward one character
<kbd>Ctrl + B</kbd>  | Backward one character
<kbd>Ctrl + D</kbd>  | Delete character under the cursor
<kbd>Ctrl + H</kbd>  | Delete character before the cursor (Backspace)
<kbd>Ctrl + W</kbd>  | Cut the word before the cursor to the clipboard
<kbd>Ctrl + K</kbd>  | Cut the line after the cursor to the clipboard
<kbd>Ctrl + U</kbd>  | Cut the line before the cursor to the clipboard
<kbd>Ctrl + L</kbd>  | Clear the screen

### History

You can use <kbd>Up arrow</kbd> and <kbd>Down arrow</kbd> to walk through the history of commands executed.

[![History](https://github.com/c-bata/assets/raw/master/go-prompt/history.gif)](#history)

### Multiple platform support

We have confirmed go-prompt works fine in the following terminals:

* iTerm2 (macOS)
* Terminal.app (macOS)
* Command Prompt (Windows)
* gnome-terminal (Ubuntu)

## Links

* [Change Log](./CHANGELOG.md)
* [GoDoc](http://godoc.org/github.com/c-bata/go-prompt)
* [gocover.io](https://gocover.io/github.com/c-bata/go-prompt)

## Author

Masashi Shibata

* Twitter: [@c\_bata\_](https://twitter.com/c_bata_/)
* Github: [@c-bata](https://github.com/c-bata/)

## License

This software is licensed under the MIT license, see [LICENSE](./LICENSE) for more information.

