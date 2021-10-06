package prompt

// LexerFunc is a callback from render.
type LexerFunc = func(line string) []LexerElement

// LexerElement is a element of lexer.
type LexerElement struct {
	Color Color
	Text  string
}

// Lexer is a struct with lexer param and function.
type Lexer struct {
	IsEnabled bool
	fn        LexerFunc
}

// NewLexer returns new Lexer.
func NewLexer() *Lexer {
	return &Lexer{
		IsEnabled: false,
	}
}

// SetLexerFunction in lexer struct.
func (l *Lexer) SetLexerFunction(fn LexerFunc) {
	l.IsEnabled = true
	l.fn = fn
}

// Process line with a custom function.
func (l *Lexer) Process(line string) []LexerElement {
	return l.fn(line)
}
