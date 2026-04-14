package tokenizer_test

import (
	"testing"

	"github.com/envoy-cli/envoy/internal/envfile"
	"github.com/envoy-cli/envoy/internal/tokenizer"
)

func entries(pairs ...string) []envfile.Entry {
	var out []envfile.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, envfile.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestTokenize_BasicEntries(t *testing.T) {
	result := tokenizer.Tokenize(entries("FOO", "bar", "BAZ", "qux"))
	if len(result) != 6 {
		t.Fatalf("expected 6 tokens, got %d", len(result))
	}
	if result[0].Kind != tokenizer.TokenKey || result[0].Raw != "FOO" {
		t.Errorf("unexpected first token: %+v", result[0])
	}
	if result[1].Kind != tokenizer.TokenEquals {
		t.Errorf("expected equals token, got %+v", result[1])
	}
	if result[2].Kind != tokenizer.TokenValue || result[2].Raw != "bar" {
		t.Errorf("unexpected value token: %+v", result[2])
	}
}

func TestTokenize_BlankEntry(t *testing.T) {
	e := []envfile.Entry{{Key: "", Value: ""}}
	tokens := tokenizer.Tokenize(e)
	if len(tokens) != 1 || tokens[0].Kind != tokenizer.TokenBlank {
		t.Errorf("expected blank token, got %+v", tokens)
	}
}

func TestTokenize_CommentEntry(t *testing.T) {
	e := []envfile.Entry{{Key: "# this is a comment", Value: ""}}
	tokens := tokenizer.Tokenize(e)
	if len(tokens) != 1 || tokens[0].Kind != tokenizer.TokenComment {
		t.Errorf("expected comment token, got %+v", tokens)
	}
	if tokens[0].Raw != "# this is a comment" {
		t.Errorf("unexpected raw value: %s", tokens[0].Raw)
	}
}

func TestTokenize_LineNumbers(t *testing.T) {
	tokens := tokenizer.Tokenize(entries("A", "1", "B", "2"))
	for _, tok := range tokens {
		if tok.Line == 0 {
			t.Errorf("token has zero line number: %+v", tok)
		}
	}
}

func TestFilterByKind_Keys(t *testing.T) {
	tokens := tokenizer.Tokenize(entries("X", "a", "Y", "b"))
	keys := tokenizer.FilterByKind(tokens, tokenizer.TokenKey)
	if len(keys) != 2 {
		t.Fatalf("expected 2 key tokens, got %d", len(keys))
	}
	if keys[0].Raw != "X" || keys[1].Raw != "Y" {
		t.Errorf("unexpected key tokens: %+v", keys)
	}
}

func TestFilterByKind_Empty(t *testing.T) {
	tokens := tokenizer.Tokenize(entries("Z", "val"))
	blanks := tokenizer.FilterByKind(tokens, tokenizer.TokenBlank)
	if len(blanks) != 0 {
		t.Errorf("expected no blank tokens, got %d", len(blanks))
	}
}
