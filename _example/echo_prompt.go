package main

import (
	"fmt"

	"github.com/c-bata/go-prompt-toolkit"
)

func executor(t string) string {
	r := "Your input: " + t
	return r
}

func completer(t string) []*prompt.Suggestion {
	return []*prompt.Suggestion{
		{Text: "users"},
		{Text: "sites"},
		{Text: "articles"},
		{Text: "comments"},
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
