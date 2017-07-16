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
		if []rune(w)[0] == []rune("s")[0] {
			return []string{"select"}
		} else if []rune(w)[0] == []rune("w")[0] {
			return []string{"where"}
		} else if []rune(w)[0] == []rune("d")[0] {
			return []string{"drop", "delete"}
		} else if []rune(w)[0] == []rune("f")[0] {
			return []string{"from"}
		}
	}
	return []string{
		"aaaaaaaaaabbbbbbbbbbccccccccccddddddddddeeeeeeeeeeffffffffff",
		"aaaaaaaaaabbbbbbbbbbccccccccccddddddddddeeeeeeeeeeffffffffff",
		"aaaaaaaaaabbbbbbbbbbccccccccccddddddddddeeeeeeeeeeffffffffff",
		"aaaaaaaaaabbbbbbbbbbccccccccccddddddddddeeeeeeeeeeffffffffff",
		"aaaaaaaaaabbbbbbbbbbccccccccccddddddddddeeeeeeeeeeffffffffff",
		"aaaaaaaaaabbbbbbbbbbccccccccccddddddddddeeeeeeeeeeffffffffff",
		"aaaaaaaaaabbbbbbbbbbccccccccccddddddddddeeeeeeeeeeffffffffff",
		"aaaaaaaaaabbbbbbbbbbccccccccccddddddddddeeeeeeeeeeffffffffff",
		"aaaaaaaaaabbbbbbbbbbccccccccccddddddddddeeeeeeeeeeffffffffff",
		"aaaaaaaaaabbbbbbbbbbccccccccccddddddddddeeeeeeeeeeffffffffff",
		"aaaaaaaaaabbbbbbbbbbccccccccccddddddddddeeeeeeeeeeffffffffff",
		"aaaaaaaaaabbbbbbbbbbccccccccccddddddddddeeeeeeeeeeffffffffff",
		"aaaaaaaaaabbbbbbbbbbccccccccccddddddddddeeeeeeeeeeffffffffff",
		"aaaaaaaaaabbbbbbbbbbccccccccccddddddddddeeeeeeeeeeffffffffff",
		"aaaaaaaaaabbbbbbbbbbccccccccccddddddddddeeeeeeeeeeffffffffff",
	}
}

func main() {
	pt := prompt.NewPrompt(
		executor,
		completer,
		prompt.MaxCompletionsOption(8),
		prompt.PrefixOption(">>> "),
		prompt.PrefixColorOption("blue"),
		prompt.TitleOption("Hello! this is prompt toolkit"),
	)
	defer fmt.Println("\nGoodbye!")
	fmt.Print("Hello! This is a example appication using prompt-toolkit.\n")
	pt.Run()
}
