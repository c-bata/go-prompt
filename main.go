package main

import (
	"fmt"

	"github.com/c-bata/go-prompt-toolkit/prompt"
)

func executor(b *prompt.Buffer) string {
	r := "Your input: " + b.Text()
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
	pt := prompt.NewPrompt(executor, completer, 8)
	defer fmt.Println("\nGoodbye!")
	fmt.Print("Hello! This is a example appication using prompt-toolkit.\n")
	pt.Run()
}
