package main

import (
	"fmt"
	"os"
	"strings"

	prompt "github.com/elk-language/go-prompt"
	"github.com/elk-language/go-prompt/completer"
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
		prompt.WithPrefix(">>> "),
		prompt.WithCompletionWordSeparator(completer.FilePathCompletionSeparator),
	)
	p.Run()
}
