package main

import (
	"fmt"
	"strings"

	"github.com/c-bata/go-prompt"
)

var ApplicationState struct {
	CWD string
}

func livePrompt() string {
	return fmt.Sprintf("(%s) >> ", ApplicationState.CWD)
}

func executor(in string) {
	ApplicationState.CWD = strings.Split(in, " ")[1]
}

func completer(in prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "cd [target]", Description: "change directory to"},
	}
	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func main() {
	p := prompt.New(
		executor,
		completer,
		prompt.OptionLivePrefix(livePrompt),
		prompt.OptionTitle("livePrefixPrompt"),
	)
	p.Run()
}
