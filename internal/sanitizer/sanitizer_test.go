package sanitizer_test

import (
	"testing"

	"github.com/user/envoy-cli/internal/envfile"
	"github.com/user/envoy-cli/internal/sanitizer"
)

func entries(pairs ...string) []envfile.Entry {
	var out []envfile.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, envfile.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestSanitize_TrimSpace(t *testing.T) {
	input := entries("  KEY  ", "  value  ")
	opts := sanitizer.Options{TrimSpace: true}
	got := sanitizer.Sanitize(input, opts)

	if got[0].Key != "KEY" {
		t.Errorf("expected key 'KEY', got %q", got[0].Key)
	}
	if got[0].Value != "value" {
		t.Errorf("expected value 'value', got %q", got[0].Value)
	}
}

func TestSanitize_RemoveInlineComments(t *testing.T) {
	input := entries("KEY", "hello world # this is a comment")
	opts := sanitizer.Options{RemoveInlineComments: true, TrimSpace: true}
	got := sanitizer.Sanitize(input, opts)

	if got[0].Value != "hello world" {
		t.Errorf("expected 'hello world', got %q", got[0].Value)
	}
}

func TestSanitize_QuotedValuePreservesComment(t *testing.T) {
	input := entries("KEY", `"hello # not a comment"`)
	opts := sanitizer.Options{RemoveInlineComments: true}
	got := sanitizer.Sanitize(input, opts)

	if got[0].Value != `"hello # not a comment"` {
		t.Errorf("quoted value should not be stripped, got %q", got[0].Value)
	}
}

func TestSanitize_UppercaseKeys(t *testing.T) {
	input := entries("db_host", "localhost", "Api_Key", "abc123")
	opts := sanitizer.Options{UppercaseKeys: true}
	got := sanitizer.Sanitize(input, opts)

	if got[0].Key != "DB_HOST" {
		t.Errorf("expected 'DB_HOST', got %q", got[0].Key)
	}
	if got[1].Key != "API_KEY" {
		t.Errorf("expected 'API_KEY', got %q", got[1].Key)
	}
}

func TestSanitize_RemoveEmpty(t *testing.T) {
	input := entries("KEY1", "value", "KEY2", "", "KEY3", "other")
	opts := sanitizer.Options{RemoveEmpty: true}
	got := sanitizer.Sanitize(input, opts)

	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if got[0].Key != "KEY1" || got[1].Key != "KEY3" {
		t.Errorf("unexpected entries: %+v", got)
	}
}

func TestSanitize_DefaultOptions(t *testing.T) {
	opts := sanitizer.DefaultOptions()
	if !opts.TrimSpace {
		t.Error("expected TrimSpace to be true by default")
	}
	if !opts.RemoveInlineComments {
		t.Error("expected RemoveInlineComments to be true by default")
	}
	if opts.UppercaseKeys {
		t.Error("expected UppercaseKeys to be false by default")
	}
	if opts.RemoveEmpty {
		t.Error("expected RemoveEmpty to be false by default")
	}
}

func TestSanitize_DoesNotMutateOriginal(t *testing.T) {
	input := entries("  key  ", "  val  ")
	orig := input[0].Key
	sanitizer.Sanitize(input, sanitizer.DefaultOptions())
	if input[0].Key != orig {
		t.Error("Sanitize must not mutate the original slice")
	}
}
