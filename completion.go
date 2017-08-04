package prompt

type Suggest struct {
	Text        string
	Description string
}

type CompletionManager struct {
	selected int // -1 means nothing one is selected.
	tmp      []*Suggest
	Max      uint16
}

func (c *CompletionManager) Reset() {
	c.selected = -1
	return
}

func (c *CompletionManager) Update(new []*Suggest) {
	c.selected = -1
	c.tmp = new
	return
}

func (c *CompletionManager) Previous() {
	c.selected--
	return
}

func (c *CompletionManager) Next() {
	c.selected++
	return
}

func (c *CompletionManager) Completing() bool {
	return c.selected != -1
}

func (c *CompletionManager) update(completions []Suggest) {
	max := int(c.Max)
	if len(completions) < max {
		max = len(completions)
	}
	if c.selected >= max {
		c.Reset()
	} else if c.selected < -1 {
		c.selected = max - 1
	}
}

func NewCompletionManager(max uint16) *CompletionManager {
	return &CompletionManager{
		selected: -1,
		Max:      max,
	}
}
