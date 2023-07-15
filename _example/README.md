# Examples of go-prompt

This directory includes some examples using go-prompt.
These examples are useful to know the usage of go-prompt and check behavior for development.

## even-lexer

Uses a custom lexer that colours every character with an even index green.

Shows you how to hook up a custom lexer for syntax highlighting.

## bang-executor

Inserts a newline when the <kbd>Enter</kbd> key is pressed unless the input ends with an exclamation point `!` (then it gets printed).

Shows you how to define a custom callback which determines whether the input is complete and should be executed or a newline should be inserted (after <kbd>Enter</kbd> has been pressed).

## simple-echo

![simple-input](https://github.com/c-bata/assets/raw/master/go-prompt/examples/input.gif)

A simple echo example using `prompt.Input`.

## http-prompt

![http-prompt](https://github.com/c-bata/assets/raw/master/go-prompt/examples/http-prompt.gif)

A simple [http-prompt](https://github.com/eliangcs/http-prompt) implementation using go-prompt in less than 200 lines of Go.

## live-prefix

![live-prefix](https://github.com/c-bata/assets/raw/master/go-prompt/examples/live-prefix.gif)

A example application which changes the prefix string dynamically.
This feature is used like [ktr0731/evans](https://github.com/ktr0731/evans) which is interactive gRPC client using go-prompt.

## exec-command

Run another CLI tool via `os/exec` package.
More practical example is [a source code of kube-prompt](https://github.com/c-bata/kube-prompt).
I recommend you to look this if you want to create tools like kube-prompt.

