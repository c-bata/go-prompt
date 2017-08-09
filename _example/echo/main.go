package main

import (
	"fmt"

	"github.com/c-bata/go-prompt"
	"strings"
)

func executor(in string) {
	fmt.Println("Your input: " + in)
}

func completer(in string) []prompt.Suggest {
	args := strings.Split(in, " ")
	last := args[0]
	if len(args) > 0 {
		last = args[len(args) - 1]
	}
	s := []prompt.Suggest{
		{Text: "users", Description: "user table"},
		{Text: "sites", Description: "sites table"},
		{Text: "articles", Description: "articles table"},
		{Text: "comments", Description: "comments table"},
	}
	return prompt.FilterHasPrefix(s, last, true)
}

func main() {
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(">>> "),
		prompt.OptionTitle("sql-prompt"),
	)
	p.Run()
}
