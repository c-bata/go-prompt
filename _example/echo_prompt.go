package main

import (
	"fmt"

	"github.com/c-bata/go-prompt-toolkit/prompt"
)

func executor(b *prompt.Buffer) string {
	r := "Your input: " + b.Text()
	return r
}

func completer(b *prompt.Buffer) []string {
	if w := b.Document().GetWordBeforeCursor(); w == "" {
		return []string{}
	}
	return []string{
		"users",
		"sites",
		"articles",
		"comments",
	}
}

func main() {
	pt := prompt.NewPrompt(
		executor,
		completer,
		prompt.OptionMaxCompletions(8),
		prompt.OptionPrefix(">>> "),
		prompt.OptionTitle("SQLITE CLI"),
	)
	defer fmt.Println("\nGoodbye!")
	pt.Run()
}
