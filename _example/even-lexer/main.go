package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"strings"
)

func main() {
	p := prompt.New(
		executor,
		completer,
		prompt.OptionSetLexer(lexer),
	)

	p.Run()
}

func lexer(line string) []prompt.LexerElement {
	var elements []prompt.LexerElement

	strArr := strings.Split(line, "")

	for k, v := range strArr {
		element := prompt.LexerElement{
			Text: v,
		}

		// every even char must be green.
		if k%2 == 0 {
			element.Color = prompt.Green
		} else {
			element.Color = prompt.White
		}

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
