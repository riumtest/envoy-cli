package truncator_test

import (
	"strings"
	"testing"

	"github.com/envoy-cli/envoy-cli/internal/envfile"
	"github.com/envoy-cli/envoy-cli/internal/truncator"
)

func entries(pairs ...string) []envfile.Entry {
	out := make([]envfile.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, envfile.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestTruncate_ShortValuesUnchanged(t *testing.T) {
	in := entries("KEY", "short")
	out := truncator.Truncate(in, truncator.DefaultOptions())
	if out[0].Value != "short" {
		t.Fatalf("expected 'short', got %q", out[0].Value)
	}
}

func TestTruncate_LongValueTruncated(t *testing.T) {
	long := strings.Repeat("x", 100)
	in := entries("KEY", long)
	opts := truncator.Options{MaxLength: 10, Suffix: "..."}
	out := truncator.Truncate(in, opts)
	if len([]rune(out[0].Value)) != 10 {
		t.Fatalf("expected length 10, got %d", len([]rune(out[0].Value)))
	}
	if !strings.HasSuffix(out[0].Value, "...") {
		t.Fatalf("expected suffix '...', got %q", out[0].Value)
	}
}

func TestTruncate_KeepKeySkipsTruncation(t *testing.T) {
	long := strings.Repeat("a", 100)
	in := entries("SECRET", long, "OTHER", long)
	opts := truncator.Options{MaxLength: 10, Suffix: "...", KeepKeys: []string{"SECRET"}}
	out := truncator.Truncate(in, opts)
	if out[0].Value != long {
		t.Fatalf("SECRET should be unchanged")
	}
	if out[1].Value == long {
		t.Fatalf("OTHER should have been truncated")
	}
}

func TestTruncate_DoesNotMutateOriginal(t *testing.T) {
	long := strings.Repeat("z", 80)
	in := entries("K", long)
	orig := in[0].Value
	truncator.Truncate(in, truncator.Options{MaxLength: 20, Suffix: "..."})
	if in[0].Value != orig {
		t.Fatal("original slice was mutated")
	}
}

func TestTruncate_EmptySuffix(t *testing.T) {
	long := strings.Repeat("b", 50)
	in := entries("KEY", long)
	opts := truncator.Options{MaxLength: 10, Suffix: ""}
	out := truncator.Truncate(in, opts)
	if len([]rune(out[0].Value)) != 10 {
		t.Fatalf("expected length 10, got %d", len(out[0].Value))
	}
}

func TestTruncate_ZeroMaxLengthUsesDefault(t *testing.T) {
	long := strings.Repeat("c", 100)
	in := entries("KEY", long)
	opts := truncator.Options{MaxLength: 0, Suffix: "..."}
	out := truncator.Truncate(in, opts)
	def := truncator.DefaultOptions().MaxLength
	if len([]rune(out[0].Value)) > def {
		t.Fatalf("expected at most %d runes, got %d", def, len([]rune(out[0].Value)))
	}
}
