package main

import (
	"fmt"

	"github.com/c-bata/go-prompt-toolkit"
)

func executor(t string) string {
	return fmt.Sprintln("Your input: " + t)
}

func completer(t string) []prompt.Completion {
	return []prompt.Completion{
		{Text: "users", Description: "user table"},
		{Text: "sites", Description: "sites table"},
		{Text: "articles", Description: "articles table"},
		{Text: "comments", Description: "comments table"},
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
