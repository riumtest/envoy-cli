// Package tokenizer splits .env file entries into structured tokens
// suitable for syntax highlighting, linting, or analysis.
package tokenizer

import (
	"strings"

	"github.com/envoy-cli/envoy/internal/envfile"
)

// TokenKind represents the type of a token.
type TokenKind string

const (
	TokenKey     TokenKind = "key"
	TokenEquals  TokenKind = "equals"
	TokenValue   TokenKind = "value"
	TokenComment TokenKind = "comment"
	TokenBlank   TokenKind = "blank"
)

// Token represents a single lexical unit from an .env line.
type Token struct {
	Kind  TokenKind
	Raw   string
	Line  int
}

// Tokenize converts a slice of envfile.Entry values into a flat list of tokens.
// Each entry produces a key token, an equals token, and a value token.
// Comment and blank lines are represented as single tokens.
func Tokenize(entries []envfile.Entry) []Token {
	var tokens []Token
	for i, e := range entries {
		line := i + 1
		raw := strings.TrimSpace(e.Key)
		if raw == "" {
			tokens = append(tokens, Token{Kind: TokenBlank, Raw: "", Line: line})
			continue
		}
		if strings.HasPrefix(raw, "#") {
			tokens = append(tokens, Token{Kind: TokenComment, Raw: raw, Line: line})
			continue
		}
		tokens = append(tokens,
			Token{Kind: TokenKey, Raw: e.Key, Line: line},
			Token{Kind: TokenEquals, Raw: "=", Line: line},
			Token{Kind: TokenValue, Raw: e.Value, Line: line},
		)
	}
	return tokens
}

// FilterByKind returns only tokens of the specified kind.
func FilterByKind(tokens []Token, kind TokenKind) []Token {
	var out []Token
	for _, t := range tokens {
		if t.Kind == kind {
			out = append(out, t)
		}
	}
	return out
}
