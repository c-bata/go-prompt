package main

import (
	"fmt"

	"github.com/c-bata/go-prompt-toolkit/prompt"
)

func executor(t string) string {
	r := "Your input: " + t
	return r
}

func completer(t string) []string {
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
		prompt.OptionPrefix(">>> "),
		prompt.OptionTitle("sqlite3-cli"),
		prompt.OptionOutputTextColor(prompt.DarkGray),
	)
	defer fmt.Println("\nGoodbye!")
	pt.Run()
}
