package main

import (
	"fmt"

	"github.com/c-bata/go-prompt-toolkit"
)

func executor(t string) {
	fmt.Println("Your input: " + t)
}

func completer(t string) []prompt.Suggest {
	return []prompt.Suggest{
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
	)
	defer fmt.Println("\nGoodbye!")
	pt.Run()
}
