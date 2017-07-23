package main

import (
	"fmt"
	"os/exec"
	"context"
	"time"

	"github.com/c-bata/go-prompt-toolkit"
)

func executor(ctx context.Context, t string) string {
	ctx, cancel := context.WithTimeout(ctx, 10 * time.Second)
	defer cancel()

	if t == "sleep 5s" {
		cmd := exec.CommandContext(ctx, "sleep", "5")
		cmd.Run()
		fmt.Println("Foo")
	} else if t == "sleep 20s" {
		cmd := exec.CommandContext(ctx, "sleep", "20")
		cmd.Run()
		fmt.Println("Foo")
	}
	return ""
}

func completer(t string) []prompt.Completion {
	return []prompt.Completion{
		{Text: "sleep 5s"},
		{Text: "sleep 20s"},
	}
}

func main() {
	pt := prompt.NewPrompt(
		executor,
		completer,
		prompt.OptionOutputTextColor(prompt.DarkGray),
	)
	defer fmt.Println("\nGoodbye!")
	pt.Run()
}
