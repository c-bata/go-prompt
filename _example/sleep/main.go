package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/c-bata/go-prompt"
)

func executor(t string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if t == "sleep 5s" {
		cmd := exec.CommandContext(ctx, "sleep", "5")
		cmd.Run()
	} else if t == "sleep 20s" {
		cmd := exec.CommandContext(ctx, "sleep", "20")
		cmd.Run()
	}
	return
}

func completer(t prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "sleep 5s"},
		{Text: "sleep 20s"},
	}
}

func main() {
	p := prompt.New(
		executor,
		completer,
	)
	defer fmt.Println("Goodbye!")
	p.Run()
}
