package prompt

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestEagerLexerNext(t *testing.T) {
	tests := map[string]struct {
		lexer *EagerLexer
		want  Token
		ok    bool
	}{
		"return the first token when at the beginning": {
			lexer: &EagerLexer{
				tokens: []Token{
					&SimpleToken{lexeme: "foo"},
					&SimpleToken{lexeme: "bar"},
				},
				currentIndex: 0,
			},
			want: &SimpleToken{lexeme: "foo"},
			ok:   true,
		},
		"return the second token": {
			lexer: &EagerLexer{
				tokens: []Token{
					&SimpleToken{lexeme: "foo"},
					&SimpleToken{lexeme: "bar"},
					&SimpleToken{lexeme: "baz"},
				},
				currentIndex: 1,
			},
			want: &SimpleToken{lexeme: "bar"},
			ok:   true,
		},
		"return false when at the end": {
			lexer: &EagerLexer{
				tokens: []Token{
					&SimpleToken{lexeme: "foo"},
					&SimpleToken{lexeme: "bar"},
					&SimpleToken{lexeme: "baz"},
				},
				currentIndex: 3,
			},
			want: nil,
			ok:   false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, ok := tc.lexer.Next()
			opts := []cmp.Option{
				cmp.AllowUnexported(SimpleToken{}, EagerLexer{}),
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.ok, ok, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func charLex(s string) []Token {
	var result []Token
	for _, char := range s {
		result = append(result, NewSimpleToken(0, string(char)))
	}

	return result
}

func TestEagerLexerInit(t *testing.T) {
	tests := map[string]struct {
		lexer *EagerLexer
		input string
		want  *EagerLexer
	}{
		"reset the lexer's state": {
			lexer: &EagerLexer{
				lexFunc: charLex,
				tokens: []Token{
					&SimpleToken{lexeme: "foo"},
					&SimpleToken{lexeme: "bar"},
				},
				currentIndex: 2,
			},
			input: "foo",
			want: &EagerLexer{
				lexFunc: charLex,
				tokens: []Token{
					&SimpleToken{lexeme: "f"},
					&SimpleToken{lexeme: "o"},
					&SimpleToken{lexeme: "o"},
				},
				currentIndex: 0,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.lexer.Init(tc.input)
			opts := []cmp.Option{
				cmp.AllowUnexported(SimpleToken{}, EagerLexer{}),
				cmpopts.IgnoreFields(EagerLexer{}, "lexFunc"),
			}
			if diff := cmp.Diff(tc.want, tc.lexer, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
