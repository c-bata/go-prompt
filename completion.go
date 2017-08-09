package prompt

import (
	"log"
	"strings"
)

type Suggest struct {
	Text        string
	Description string
}

type CompletionManager struct {
	selected  int // -1 means nothing one is selected.
	tmp       []Suggest
	Max       uint16
	completer Completer
}

func (c *CompletionManager) GetSelectedSuggestion() (s Suggest, ok bool) {
	if c.selected == -1 {
		return Suggest{}, false
	} else if c.selected < -1 {
		log.Printf("[ERROR] shoud be reached here, selected=%d", c.selected)
		return Suggest{}, false
	}
	return c.tmp[c.selected], true
}

func (c *CompletionManager) GetSuggestions() []Suggest {
	return c.tmp
}

func (c *CompletionManager) Reset() {
	c.selected = -1
	c.Update("")
	return
}

func (c *CompletionManager) Update(in string) {
	c.tmp = c.completer(in)
	return
}

func (c *CompletionManager) Previous() {
	c.selected--
	c.update()
	return
}

func (c *CompletionManager) Next() {
	c.selected++
	c.update()
	return
}

func (c *CompletionManager) Completing() bool {
	return c.selected != -1
}

func (c *CompletionManager) update() {
	max := int(c.Max)
	if len(c.tmp) < max {
		max = len(c.tmp)
	}
	if c.selected >= max {
		c.Reset()
	} else if c.selected < -1 {
		c.selected = max - 1
	}
}

func formatCompletions(completions []Suggest, max int) (new []Suggest, width int) {
	num := len(completions)
	new = make([]Suggest, num)
	leftWidth := 0
	rightWidth := 0

	for i := 0; i < num; i++ {
		if leftWidth < len([]rune(completions[i].Text)) {
			leftWidth = len([]rune(completions[i].Text))
		}
		if rightWidth < len([]rune(completions[i].Description)) {
			rightWidth = len([]rune(completions[i].Description))
		}
	}

	if diff := max - completionMargin - leftWidth - rightWidth; diff < 0 {
		if rightWidth > diff {
			rightWidth += diff
		} else if rightWidth+rightMargin > -diff {
			leftWidth += rightWidth + rightMargin + diff
			rightWidth = 0
		}
	}
	if rightWidth == 0 {
		width = leftWidth + leftMargin
	} else {
		width = leftWidth + leftMargin + rightWidth + rightMargin
	}

	for i := 0; i < num; i++ {
		var newText string
		var newDescription string
		if l := len(completions[i].Text); l > leftWidth {
			newText = leftPrefix + completions[i].Text[:leftWidth-len("...")] + "..." + leftSuffix
		} else if l < width {
			spaces := strings.Repeat(" ", leftWidth-len([]rune(completions[i].Text)))
			newText = leftPrefix + completions[i].Text + spaces + leftSuffix
		} else {
			newText = leftPrefix + completions[i].Text + leftSuffix
		}

		if rightWidth == 0 {
			newDescription = ""
		} else if l := len(completions[i].Description); l > rightWidth {
			newDescription = rightPrefix + completions[i].Description[:rightWidth-len("...")] + "..." + rightSuffix
		} else if l < width {
			spaces := strings.Repeat(" ", rightWidth-len([]rune(completions[i].Description)))
			newDescription = rightPrefix + completions[i].Description + spaces + rightSuffix
		} else {
			newDescription = rightPrefix + completions[i].Description + rightSuffix
		}
		new[i] = Suggest{Text: newText, Description: newDescription}
	}
	return
}

func NewCompletionManager(completer Completer, max uint16) *CompletionManager {
	return &CompletionManager{
		selected:  -1,
		Max:       max,
		completer: completer,
	}
}
