package prompt

import (
	"context"
	"log"
	"runtime"

	"github.com/mattn/go-runewidth"
)

// NewRenderer returns Renderer object.
func NewRenderer(out ConsoleWriter, initialRequest RenderRequest, opts ...RendererOption) *Renderer {
	ws := out.GetWinSize()
	renderer := &Renderer{
		out:             out,
		col:             ws.Col,
		row:             ws.Row,
		previousRequest: initialRequest,

		// Prefix
		prefix: "> ",

		// Colors
		prefixTextColor:              Blue,
		prefixBGColor:                DefaultColor,
		inputTextColor:               DefaultColor,
		inputBGColor:                 DefaultColor,
		previewSuggestionTextColor:   Green,
		previewSuggestionBGColor:     DefaultColor,
		suggestionTextColor:          White,
		suggestionBGColor:            Cyan,
		selectedSuggestionTextColor:  Black,
		selectedSuggestionBGColor:    Turquoise,
		descriptionTextColor:         Black,
		descriptionBGColor:           Turquoise,
		selectedDescriptionTextColor: White,
		selectedDescriptionBGColor:   Cyan,
		scrollbarThumbColor:          DarkGray,
		scrollbarBGColor:             Cyan,

		// Channels
		Clear:  make(chan struct{}),
		Render: make(chan RenderRequest),
	}

	for _, o := range opts {
		o(renderer)
	}
	return renderer
}

// RendererOption is the type to replace default renderer parameters.
// NewRenderer accepts any number of options.
type RendererOption func(render *Renderer)

// Renderer to render prompt information from the state of Buffer.
type Renderer struct {
	out   ConsoleWriter
	title string
	row   uint16
	col   uint16

	previousCursor  int
	previousRequest RenderRequest

	// Prefix
	prefix             string
	livePrefixCallback func() (prefix string, useLivePrefix bool)

	// Colors
	prefixTextColor              Color
	prefixBGColor                Color
	inputTextColor               Color
	inputBGColor                 Color
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
	scrollbarThumbColor          Color
	scrollbarBGColor             Color

	// Channels
	Clear  chan struct{}
	Render chan RenderRequest
}

// Run to start worker goroutine for rendering.
func (r *Renderer) Run(ctx context.Context) {
	r.SetUp()
	defer r.TearDown()
	r.render(r.previousRequest)
	sigwinch := r.out.SIGWINCH(ctx)

	for {
		select {
		case <-ctx.Done():
			r.breakLine(r.previousRequest.buffer)
			return
		case <-r.Clear:
			log.Print("renderer: catch erase request")
			r.out.EraseScreen()
			r.out.CursorGoTo(0, 0)
			r.out.Flush()
		case req := <-r.Render:
			r.render(req)
		case ws := <-sigwinch:
			r.col = ws.Col
			r.row = ws.Row
			r.render(r.previousRequest)
		}
	}
}

// RenderRequest is used when requesting render
type RenderRequest struct {
	buffer     *Buffer
	completion *CompletionManager
}

// SetUp to initialize console output.
func (r *Renderer) SetUp() {
	if r.title != "" {
		r.out.SetTitle(r.title)
		r.out.Flush()
	}
}

// getCurrentPrefix to get current prefix.
// If live-prefix is enabled, return live-prefix.
func (r *Renderer) getCurrentPrefix() string {
	if r.livePrefixCallback == nil {
		return r.prefix
	}
	prefix, ok := r.livePrefixCallback()
	if ok {
		return prefix
	}
	return r.prefix
}

func (r *Renderer) renderPrefix() {
	r.out.SetColor(r.prefixTextColor, r.prefixBGColor, false)
	r.out.WriteStr(r.getCurrentPrefix())
	r.out.SetColor(DefaultColor, DefaultColor, false)
}

// TearDown to clear title and erasing.
func (r *Renderer) TearDown() {
	r.out.ClearTitle()
	r.out.EraseDown()
	r.out.Flush()

	close(r.Clear)
	close(r.Render)
}

func (r *Renderer) prepareArea(lines int) {
	for i := 0; i < lines; i++ {
		r.out.ScrollDown()
	}
	for i := 0; i < lines; i++ {
		r.out.ScrollUp()
	}
	return
}

func (r *Renderer) renderWindowTooSmall() {
	r.out.CursorGoTo(0, 0)
	r.out.EraseScreen()
	r.out.SetColor(DarkRed, White, false)
	r.out.WriteStr("Your console window is too small...")
	return
}

func (r *Renderer) renderCompletion(buf *Buffer, completions *CompletionManager) {
	suggestions := completions.GetSuggestions()
	if len(completions.GetSuggestions()) == 0 {
		return
	}
	prefix := r.getCurrentPrefix()
	formatted, width := formatSuggestions(
		suggestions,
		int(r.col)-runewidth.StringWidth(prefix)-1, // -1 means a width of scrollbar
	)
	// +1 means a width of scrollbar.
	width++

	windowHeight := len(formatted)
	if windowHeight > int(completions.max) {
		windowHeight = int(completions.max)
	}
	formatted = formatted[completions.verticalScroll : completions.verticalScroll+windowHeight]
	r.prepareArea(windowHeight)

	cursor := runewidth.StringWidth(prefix) + runewidth.StringWidth(buf.Document().TextBeforeCursor())
	x, _ := r.toPos(cursor)
	if x+width >= int(r.col) {
		cursor = r.backward(cursor, x+width-int(r.col))
	}

	contentHeight := len(completions.tmp)

	fractionVisible := float64(windowHeight) / float64(contentHeight)
	fractionAbove := float64(completions.verticalScroll) / float64(contentHeight)

	scrollbarHeight := int(clamp(float64(windowHeight), 1, float64(windowHeight)*fractionVisible))
	scrollbarTop := int(float64(windowHeight) * fractionAbove)

	isScrollThumb := func(row int) bool {
		return scrollbarTop <= row && row <= scrollbarTop+scrollbarHeight
	}

	selected := completions.selected - completions.verticalScroll
	r.out.SetColor(White, Cyan, false)
	for i := 0; i < windowHeight; i++ {
		r.out.CursorDown(1)
		if i == selected {
			r.out.SetColor(r.selectedSuggestionTextColor, r.selectedSuggestionBGColor, true)
		} else {
			r.out.SetColor(r.suggestionTextColor, r.suggestionBGColor, false)
		}
		r.out.WriteStr(formatted[i].Text)

		if i == selected {
			r.out.SetColor(r.selectedDescriptionTextColor, r.selectedDescriptionBGColor, false)
		} else {
			r.out.SetColor(r.descriptionTextColor, r.descriptionBGColor, false)
		}
		r.out.WriteStr(formatted[i].Description)

		if isScrollThumb(i) {
			r.out.SetColor(DefaultColor, r.scrollbarThumbColor, false)
		} else {
			r.out.SetColor(DefaultColor, r.scrollbarBGColor, false)
		}
		r.out.WriteStr(" ")
		r.out.SetColor(DefaultColor, DefaultColor, false)

		r.lineWrap(cursor + width)
		r.backward(cursor+width, width)
	}

	if x+width >= int(r.col) {
		r.out.CursorForward(x + width - int(r.col))
	}

	r.out.CursorUp(windowHeight)
	r.out.SetColor(DefaultColor, DefaultColor, false)
	return
}

// render renders to the console.
func (r *Renderer) render(req RenderRequest) {
	// In situations where a pseudo tty is allocated (e.g. within a docker container),
	// window size via TIOCGWINSZ is not immediately available and will result in 0,0 dimensions.
	if r.col == 0 {
		return
	}
	defer r.out.Flush()
	r.move(r.previousCursor, 0)

	line := req.buffer.Text()
	prefix := r.getCurrentPrefix()
	cursor := runewidth.StringWidth(prefix) + runewidth.StringWidth(line)

	// prepare area
	_, y := r.toPos(cursor)

	h := y + 1 + int(req.completion.max)
	if h > int(r.row) || completionMargin > int(r.col) {
		r.renderWindowTooSmall()
		return
	}

	// Rendering
	r.out.HideCursor()
	defer r.out.ShowCursor()

	r.renderPrefix()
	r.out.SetColor(r.inputTextColor, r.inputBGColor, false)
	r.out.WriteStr(line)
	r.out.SetColor(DefaultColor, DefaultColor, false)
	r.lineWrap(cursor)

	r.out.EraseDown()

	cursor = r.backward(cursor, runewidth.StringWidth(line)-req.buffer.DisplayCursorPosition())

	r.renderCompletion(req.buffer, req.completion)
	if suggest, ok := req.completion.GetSelectedSuggestion(); ok {
		cursor = r.backward(cursor, runewidth.StringWidth(req.buffer.Document().GetWordBeforeCursorUntilSeparator(req.completion.wordSeparator)))

		r.out.SetColor(r.previewSuggestionTextColor, r.previewSuggestionBGColor, false)
		r.out.WriteStr(suggest.Text)
		r.out.SetColor(DefaultColor, DefaultColor, false)
		cursor += runewidth.StringWidth(suggest.Text)

		rest := req.buffer.Document().TextAfterCursor()
		r.out.WriteStr(rest)
		cursor += runewidth.StringWidth(rest)
		r.lineWrap(cursor)

		cursor = r.backward(cursor, runewidth.StringWidth(rest))
	}
	r.previousCursor = cursor
	r.previousRequest = req
}

// breakLine to break line.
func (r *Renderer) breakLine(buffer *Buffer) {
	// Erasing and render
	cursor := runewidth.StringWidth(buffer.Document().TextBeforeCursor()) + runewidth.StringWidth(r.getCurrentPrefix())
	r.clear(cursor)
	r.renderPrefix()
	r.out.SetColor(r.inputTextColor, r.inputBGColor, false)
	r.out.WriteStr(buffer.Document().Text + "\n")
	r.out.SetColor(DefaultColor, DefaultColor, false)
	r.out.Flush()
	r.previousCursor = 0
	r.previousRequest.buffer = buffer
}

// clear erases the screen from a beginning of input
// even if there is line break which means input length exceeds a window's width.
func (r *Renderer) clear(cursor int) {
	r.move(cursor, 0)
	r.out.EraseDown()
}

// backward moves cursor to backward from a current cursor position
// regardless there is a line break.
func (r *Renderer) backward(from, n int) int {
	return r.move(from, from-n)
}

// move moves cursor to specified position from the beginning of input
// even if there is a line break.
func (r *Renderer) move(from, to int) int {
	fromX, fromY := r.toPos(from)
	toX, toY := r.toPos(to)

	r.out.CursorUp(fromY - toY)
	r.out.CursorBackward(fromX - toX)
	return to
}

// toPos returns the relative position from the beginning of the string.
func (r *Renderer) toPos(cursor int) (x, y int) {
	col := int(r.col)
	return cursor % col, cursor / col
}

func (r *Renderer) lineWrap(cursor int) {
	if runtime.GOOS != "windows" && cursor > 0 && cursor%int(r.col) == 0 {
		r.out.WriteRaw([]byte{'\n'})
	}
}

func clamp(high, low, x float64) float64 {
	switch {
	case high < x:
		return high
	case x < low:
		return low
	default:
		return x
	}
}
