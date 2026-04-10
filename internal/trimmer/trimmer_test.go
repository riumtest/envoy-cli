package trimmer_test

import (
	"testing"

	"github.com/user/envoy-cli/internal/envfile"
	"github.com/user/envoy-cli/internal/trimmer"
)

func entries(pairs ...string) []envfile.Entry {
	out := make([]envfile.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, envfile.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestTrim_TrimValues(t *testing.T) {
	in := entries("KEY", "  hello  ", "OTHER", "\tworld\t")
	out := trimmer.Trim(in, trimmer.Options{TrimValues: true})
	if out[0].Value != "hello" {
		t.Errorf("expected 'hello', got %q", out[0].Value)
	}
	if out[1].Value != "world" {
		t.Errorf("expected 'world', got %q", out[1].Value)
	}
}

func TestTrim_RemoveEmpty(t *testing.T) {
	in := entries("A", "value", "B", "", "C", "other")
	out := trimmer.Trim(in, trimmer.Options{RemoveEmpty: true})
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out[0].Key != "A" || out[1].Key != "C" {
		t.Errorf("unexpected keys: %v", out)
	}
}

func TestTrim_DeduplicateKeys(t *testing.T) {
	in := entries("KEY", "first", "OTHER", "x", "KEY", "second")
	out := trimmer.Trim(in, trimmer.Options{DeduplicateKeys: true})
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out[0].Value != "second" {
		t.Errorf("expected last value 'second', got %q", out[0].Value)
	}
}

func TestTrim_DoesNotMutateOriginal(t *testing.T) {
	in := entries("KEY", "  padded  ")
	orig := in[0].Value
	trimmer.Trim(in, trimmer.DefaultOptions())
	if in[0].Value != orig {
		t.Errorf("original entry was mutated")
	}
}

func TestTrim_DefaultOptions(t *testing.T) {
	opts := trimmer.DefaultOptions()
	if !opts.TrimValues {
		t.Error("expected TrimValues to be true by default")
	}
	if opts.RemoveEmpty {
		t.Error("expected RemoveEmpty to be false by default")
	}
	if opts.DeduplicateKeys {
		t.Error("expected DeduplicateKeys to be false by default")
	}
}

func TestTrim_EmptyInput(t *testing.T) {
	out := trimmer.Trim(nil, trimmer.DefaultOptions())
	if len(out) != 0 {
		t.Errorf("expected empty result, got %d entries", len(out))
	}
}
