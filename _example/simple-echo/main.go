package main

import (
	"fmt"

	prompt "github.com/elk-language/go-prompt"
)

func completer(in prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "users", Description: "Store the username and age"},
		{Text: "articles", Description: "Store the article text posted by user"},
		{Text: "comments", Description: "Store the text commented to articles"},
		{Text: "groups", Description: "Combine users with specific rules"},
	}
	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func main() {
	in := prompt.Input(
		prompt.WithPrefix(">>> "),
		prompt.WithTitle("sql-prompt"),
		prompt.WithHistory([]string{"SELECT * FROM users;"}),
		prompt.WithPrefixTextColor(prompt.Yellow),
		prompt.WithPreviewSuggestionTextColor(prompt.Blue),
		prompt.WithSelectedSuggestionBGColor(prompt.LightGray),
		prompt.WithSuggestionBGColor(prompt.DarkGray),
		prompt.WithCompleter(completer),
	)
	fmt.Println("Your input: " + in)
}
