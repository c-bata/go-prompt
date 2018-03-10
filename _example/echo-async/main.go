package main

import (
	"fmt"
	"time"

	"github.com/c-bata/go-prompt"
)

func executor(in string) {
	fmt.Println("Your input: " + in)
}

type asyncCompleter struct {
	minDelay time.Duration
	lastTime time.Time
	lastSugg []prompt.Suggest
}

func newAsyncCompleter(minDelay time.Duration) *asyncCompleter {
	return &asyncCompleter{minDelay: minDelay, lastTime: time.Now(), lastSugg: []prompt.Suggest{}}
}

func (c *asyncCompleter) completer(in prompt.Document) []prompt.Suggest {
	// Completer is called twice, but first call seems to not be ready to process input.
	// If we don't do this, the delta time will always be very fast and completion never invoked
	if in.GetWordBeforeCursor() == "" {
		now := time.Now()
		since := now.Sub(c.lastTime)
		c.lastTime = now

		s := []prompt.Suggest{}
		if since > c.minDelay {
			// Do latency sensitive population of suggest slice here
			s = []prompt.Suggest{
				{Text: "users", Description: "Store the username and age"},
				{Text: "articles", Description: "Store the article text posted by user"},
				{Text: "comments", Description: "Store the text commented to articles"},
				{Text: "groups", Description: "Combine users with specific rules"},
			}
		}
		c.lastSugg = s
		return s
	}
	return prompt.FilterHasPrefix(c.lastSugg, in.GetWordBeforeCursor(), true)
}

func main() {
	c := newAsyncCompleter(1000 * time.Millisecond)
	p := prompt.New(
		executor,
		c.completer,
		prompt.OptionPrefix(">>> "),
		prompt.OptionTitle("sql-prompt"),
	)
	p.Run()
}
