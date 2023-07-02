package main

import (
	"fmt"
	"strings"

	"github.com/elk-language/go-prompt"
)

func main() {
	p := prompt.New(
		executor,
		completer,
		prompt.OptionSetLexer(prompt.NewEagerLexer(lexer)),
	)

	p.Run()
}

func lexer(line string) []prompt.Token {
	var elements []prompt.Token

	strArr := strings.Split(line, "")

	for i, value := range strArr {
		var color prompt.Color
		// every even char must be green.
		if i%2 == 0 {
			color = prompt.Green
		} else {
			color = prompt.White
		}
		element := prompt.NewSimpleToken(color, value)

		elements = append(elements, element)
	}

	return elements
}

func completer(in prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{}
}

func executor(s string) {
	fmt.Println("You printed: " + s)
}
