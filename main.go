package main

import (
	"fmt"

	"github.com/c-bata/go-prompt-toolkit/prompt"
)

func executor(b *prompt.Buffer) string {
	r := "\n>>> Your input: '" + b.Text() + "' <<<\n"
	return r
}

func completer(b *prompt.Buffer) []string {
	if w := b.Document().GetWordBeforeCursor(); w == "" {
		return []string{}
	} else {
		return []string{"select", "from", "insert", "where"}
	}
}

func main() {
	pt := prompt.NewPrompt(executor, completer)
	defer fmt.Println("\nGoodbye!")
	fmt.Print("Hello! This is a example appication using prompt-toolkit.\n\n")
	pt.Run()
}
