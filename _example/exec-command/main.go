package main

import (
	"fmt"
	"os/exec"
	"os"

	"github.com/c-bata/go-prompt"
)

func executor(t string) {
	if t == "bash" {
		cmd := exec.Command("bash")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
	return
}

func completer(t string) []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "bash"},
	}
}

func main() {
	p := prompt.New(
		executor,
		completer,
	)
	p.Run()
}
