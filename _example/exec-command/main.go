package main

import (
	"os"
	"os/exec"

	prompt "github.com/elk-language/go-prompt"
)

func executor(t string) {
	if t != "bash" {
		return
	}

	cmd := exec.Command("bash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func completer(t prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "bash"},
	}
}

func main() {
	p := prompt.New(
		executor,
		prompt.WithCompleter(completer),
	)
	p.Run()
}
