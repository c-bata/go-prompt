package prompt

import "strings"

const (
	leftPrefix = " "
	leftSuffix = " "
	rightPrefix = " "
	rightSuffix = " "
	leftMargin = len(leftPrefix + leftSuffix)
	rightMargin = len(rightPrefix + rightSuffix)
	completionMargin = leftMargin + rightMargin
)

type Render struct {
	out    ConsoleWriter
	prefix string
	title  string
	row    uint16
	col    uint16
	// colors
	prefixTextColor              Color
	prefixBGColor                Color
	inputTextColor               Color
	inputBGColor                 Color
	outputTextColor              Color
	outputBGColor                Color
	previewSuggestionTextColor   Color
	previewSuggestionBGColor     Color
	suggestionTextColor          Color
	suggestionBGColor            Color
	selectedSuggestionTextColor  Color
	selectedSuggestionBGColor    Color
	descriptionTextColor         Color
	descriptionBGColor           Color
	selectedDescriptionTextColor Color
	selectedDescriptionBGColor   Color
}

func (r *Render) Setup() {
	if r.title != "" {
		r.out.SetTitle(r.title)
		r.out.Flush()
	}
}

func (r *Render) renderPrefix() {
	r.out.SetColor(r.prefixTextColor, r.prefixBGColor)
	r.out.WriteStr(r.prefix)
	r.out.SetColor(DefaultColor, DefaultColor)
}

func (r *Render) TearDown() {
	r.out.ClearTitle()
	r.out.EraseDown()
	r.out.Flush()
}

func (r *Render) prepareArea(lines int) {
	for i := 0; i < lines; i++ {
		r.out.ScrollDown()
	}
	for i := 0; i < lines; i++ {
		r.out.ScrollUp()
	}
	return
}

func (r *Render) UpdateWinSize(ws *WinSize) {
	r.row = ws.Row
	r.col = ws.Col
	return
}

func (r *Render) renderWindowTooSmall() {
	r.out.CursorGoTo(0, 0)
	r.out.EraseScreen()
	r.out.SetColor(DarkRed, White)
	r.out.WriteStr("Your console window is too small...")
	r.out.Flush()
	return
}

func (r *Render) renderCompletion(buf *Buffer, completions []Completion, max uint16, selected int) {
	if max > r.row {
		max = r.row
	}

	if l := len(completions); l == 0 {
		return
	} else if l > int(max) {
		completions = completions[:max]
	}

	formatted, width := formatCompletions(
		completions,
		int(r.col)-len(r.prefix),
	)
	l := len(formatted)
	r.prepareArea(l)

	d := (len(r.prefix) + len(buf.Document().TextBeforeCursor())) % int(r.col)
	if d == 0 { // the cursor is on right end.
		r.out.CursorBackward(width)
	} else if d+width > int(r.col) {
		r.out.CursorBackward(d + width - int(r.col))
	}

	r.out.SetColor(White, Cyan)
	for i := 0; i < l; i++ {
		r.out.CursorDown(1)
		if i == selected {
			r.out.SetColor(r.selectedSuggestionTextColor, r.selectedSuggestionBGColor)
		} else {
			r.out.SetColor(r.suggestionTextColor, r.suggestionBGColor)
		}
		r.out.WriteStr(formatted[i].Text)

		if i == selected {
			r.out.SetColor(r.selectedDescriptionTextColor, r.selectedDescriptionBGColor)
		} else {
			r.out.SetColor(r.descriptionTextColor, r.descriptionBGColor)
		}
		r.out.WriteStr(formatted[i].Description)
		r.out.CursorBackward(width)
	}
	if d == 0 { // the cursor is on right end.
		// DON'T CURSOR DOWN HERE. Because the line doesn't erase properly.
		r.out.CursorForward(width + 1)
	} else if d+width > int(r.col) {
		r.out.CursorForward(d + width - int(r.col))
	}

	r.out.CursorUp(l)
	r.out.SetColor(DefaultColor, DefaultColor)
	return
}

func (r *Render) Render(buffer *Buffer, completions []Completion, maxCompletions uint16, selected int) {
	// Erasing
	r.out.CursorBackward(int(r.col) + len(buffer.Text()) + len(r.prefix))
	r.out.EraseDown()

	// prepare area
	line := buffer.Text()
	h := ((len(r.prefix) + len(line)) / int(r.col)) + 1 + int(maxCompletions)
	if h > int(r.row) || completionMargin > int(r.col) {
		r.renderWindowTooSmall()
		return
	}

	// Rendering
	r.renderPrefix()

	r.out.SetColor(r.inputTextColor, r.inputBGColor)
	r.out.WriteStr(line)
	r.out.SetColor(DefaultColor, DefaultColor)
	r.out.CursorBackward(len(line) - buffer.CursorPosition)
	r.renderCompletion(buffer, completions, maxCompletions, selected)
	if selected != -1 {
		c := completions[selected]
		r.out.CursorBackward(len([]rune(buffer.Document().GetWordBeforeCursor())))
		r.out.SetColor(r.previewSuggestionTextColor, r.previewSuggestionBGColor)
		r.out.WriteStr(c.Text)
		r.out.SetColor(DefaultColor, DefaultColor)
	}
	r.out.Flush()
}

func (r *Render) BreakLine(buffer *Buffer) {
	// CR
	r.out.CursorBackward(int(r.col) + len(buffer.Text()) + len(r.prefix))
	// Erasing and Render
	r.out.EraseDown()
	r.renderPrefix()
	r.out.SetColor(r.inputTextColor, r.inputBGColor)
	r.out.WriteStr(buffer.Document().Text + "\n")
	r.out.SetColor(DefaultColor, DefaultColor)
	r.out.Flush()
}

func (r *Render) RenderResult(result string) {
	// Render Result
	if result != "" {
		r.out.SetColor(r.outputTextColor, r.outputBGColor)
		r.out.WriteRawStr(result)
	}
	r.out.SetColor(DefaultColor, DefaultColor)
}

func formatCompletions(completions []Completion, max int) (new []Completion, width int) {
	num := len(completions)
	new = make([]Completion, num)
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
			rightWidth -= diff
		} else if rightWidth+rightMargin > diff {
			leftWidth += rightWidth + rightMargin - diff
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
		new[i] = Completion{Text: newText, Description: newDescription}
	}
	return
}
