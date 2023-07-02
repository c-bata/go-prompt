package prompt

// Lexer is a streaming lexer that takes in a piece of text
// and streams tokens with the Next() method
type Lexer interface {
	Init(string) // Reset the lexer's state and initialise it with the given input.
	// Next returns the next Token and a bool flag
	// which is false when the end of input has been reached.
	Next() (Token, bool)
}

// Token is a single unit of text returned by a Lexer.
type Token interface {
	Color() Color
	Lexeme() string // original string that matches this token
}

// SimpleToken as the default implementation of Token.
type SimpleToken struct {
	color  Color
	lexeme string
}

// Create a new SimpleToken.
func NewSimpleToken(color Color, lexeme string) *SimpleToken {
	return &SimpleToken{
		color:  color,
		lexeme: lexeme,
	}
}

// Retrieve the color of this token.
func (t *SimpleToken) Color() Color {
	return t.color
}

// Retrieve the text that this token represents.
func (t *SimpleToken) Lexeme() string {
	return t.lexeme
}

// LexerFunc is a function implementing
// a simple lexer that receives a string
// and returns a complete slice of Tokens.
type LexerFunc func(string) []Token

// EagerLexer is a wrapper around LexerFunc that
// transforms an eager lexer which produces an
// array with all tokens at once into a streaming
// lexer compatible with go-prompt.
type EagerLexer struct {
	lexFunc      LexerFunc
	tokens       []Token
	currentIndex int
}

// Create a new EagerLexer.
func NewEagerLexer(fn LexerFunc) *EagerLexer {
	return &EagerLexer{
		lexFunc: fn,
	}
}

// Initialise the lexer with the given input.
func (l *EagerLexer) Init(input string) {
	l.tokens = l.lexFunc(input)
	l.currentIndex = 0
}

// Return the next token and true if the operation
// was successful.
func (l *EagerLexer) Next() (Token, bool) {
	if l.currentIndex >= len(l.tokens) {
		return nil, false
	}

	result := l.tokens[l.currentIndex]
	l.currentIndex++
	return result, true
}
