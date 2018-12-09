package main

import (
	"fmt"
	"os"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	"github.com/c-bata/go-prompt/completer"
)

var filePathCompleter = completer.FilePathCompleter{
	IgnoreCase: true,
	Filter: func(fi os.FileInfo) bool {
		return fi.IsDir() || strings.HasSuffix(fi.Name(), ".go")
	},
}

func executor(in string) {
	fmt.Println("Your input: " + in)
}

func completerFunc(d prompt.Document) []prompt.Suggest {
	t := d.GetWordBeforeCursor()
	if strings.HasPrefix(t, "--") {
		return []prompt.Suggest{
			{"--foo", ""},
			{"--bar", ""},
			{"--baz", ""},
		}
	}
	return filePathCompleter.Complete(d)
}

func main() {
	p := prompt.New(
		executor,
		completerFunc,
		prompt.OptionPrefix(">>> "),
		prompt.OptionCompletionWordSeparator(completer.FilePathCompletionSeparator),
	)
	p.Run()
}
