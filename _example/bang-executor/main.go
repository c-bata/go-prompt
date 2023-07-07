package main

import (
	"fmt"
	"unicode/utf8"

	"github.com/elk-language/go-prompt"
)

func main() {
	p := prompt.New(
		executor,
		prompt.WithExecuteOnEnterCallback(ExecuteOnEnter),
	)

	p.Run()
}

func ExecuteOnEnter(input string) bool {
	char, _ := utf8.DecodeLastRuneInString(input)
	return char == '!'
}

func executor(s string) {
	fmt.Println("You printed: " + s)
}
