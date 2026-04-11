package stripper_test

import (
	"testing"

	"github.com/user/envoy-cli/internal/envfile"
	"github.com/user/envoy-cli/internal/stripper"
)

func entries(pairs ...string) []envfile.Entry {
	out := make([]envfile.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, envfile.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestStrip_RemovesComments(t *testing.T) {
	input := []envfile.Entry{
		{Key: "# this is a comment", Value: ""},
		{Key: "APP_ENV", Value: "production"},
	}
	opts := stripper.DefaultOptions()
	result := stripper.Strip(input, opts)
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if result[0].Key != "APP_ENV" {
		t.Errorf("unexpected key %q", result[0].Key)
	}
}

func TestStrip_RemovesBlanks(t *testing.T) {
	input := entries("APP_ENV", "dev", "", "", "PORT", "8080")
	opts := stripper.DefaultOptions()
	result := stripper.Strip(input, opts)
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestStrip_DoesNotMutateOriginal(t *testing.T) {
	input := entries("  KEY  ", "  val  ")
	opts := stripper.DefaultOptions()
	_ = stripper.Strip(input, opts)
	if input[0].Key != "  KEY  " {
		t.Error("original entry was mutated")
	}
}

func TestStrip_TrimWhitespace(t *testing.T) {
	input := entries("  KEY  ", "  value  ")
	opts := stripper.DefaultOptions()
	result := stripper.Strip(input, opts)
	if result[0].Key != "KEY" {
		t.Errorf("expected trimmed key, got %q", result[0].Key)
	}
	if result[0].Value != "value" {
		t.Errorf("expected trimmed value, got %q", result[0].Value)
	}
}

func TestStripRaw_RemovesCommentLines(t *testing.T) {
	input := "# comment\nAPP=1\n# another\nDEBUG=true"
	opts := stripper.DefaultOptions()
	result := stripper.StripRaw(input, opts)
	if result != "APP=1\nDEBUG=true" {
		t.Errorf("unexpected result: %q", result)
	}
}

func TestStripRaw_RemovesBlankLines(t *testing.T) {
	input := "APP=1\n\nDEBUG=true\n"
	opts := stripper.DefaultOptions()
	result := stripper.StripRaw(input, opts)
	if result != "APP=1\nDEBUG=true" {
		t.Errorf("unexpected result: %q", result)
	}
}

func TestStripRaw_NoStrip(t *testing.T) {
	input := "# comment\nAPP=1"
	opts := stripper.Options{StripComments: false, StripBlanks: false}
	result := stripper.StripRaw(input, opts)
	if result != input {
		t.Errorf("expected unchanged content, got %q", result)
	}
}
