package main

import (
	"fmt"

	"github.com/c-bata/go-prompt-toolkit/prompt"
)

func executor(b *prompt.Buffer) string {
	r := "\n>>> Your input: '" + b.Text() + "' <<<\n"
	return r
}

func completer(b *prompt.Buffer) []string {
	return []string{"foo", "bar", "baz"}
}

func main() {
	pt := prompt.NewPrompt(executor)
	defer fmt.Println("\nGoodbye!")
	pt.Run()
}
