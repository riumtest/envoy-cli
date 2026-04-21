package shrinker_test

import (
	"testing"

	"github.com/your-org/envoy-cli/internal/envfile"
	"github.com/your-org/envoy-cli/internal/shrinker"
)

func entries() []envfile.Entry {
	return []envfile.Entry{
		{Key: "SHORT", Value: "hi"},
		{Key: "MEDIUM", Value: "hello"},
		{Key: "LONG", Value: "this-is-a-very-long-value"},
		{Key: "EMPTY", Value: ""},
	}
}

func TestShrink_NoConstraints(t *testing.T) {
	opts := shrinker.DefaultOptions()
	out := shrinker.Shrink(entries(), opts)
	if len(out) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(out))
	}
}

func TestShrink_MinLen(t *testing.T) {
	opts := shrinker.DefaultOptions()
	opts.MinLen = 3
	out := shrinker.Shrink(entries(), opts)
	// "hi" (len 2) and "" (len 0) should be dropped
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out[0].Key != "MEDIUM" || out[1].Key != "LONG" {
		t.Errorf("unexpected keys: %v", out)
	}
}

func TestShrink_MaxLen(t *testing.T) {
	opts := shrinker.DefaultOptions()
	opts.MaxLen = 5
	out := shrinker.Shrink(entries(), opts)
	// "this-is-a-very-long-value" (len 25) should be dropped
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
}

func TestShrink_MinAndMaxLen(t *testing.T) {
	opts := shrinker.DefaultOptions()
	opts.MinLen = 3
	opts.MaxLen = 10
	out := shrinker.Shrink(entries(), opts)
	// Only "hello" (len 5) survives
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0].Key != "MEDIUM" {
		t.Errorf("expected MEDIUM, got %s", out[0].Key)
	}
}

func TestShrink_KeepKeys(t *testing.T) {
	opts := shrinker.DefaultOptions()
	opts.MaxLen = 3
	opts.KeepKeys = []string{"LONG"}
	out := shrinker.Shrink(entries(), opts)
	// Only "hi" and "LONG" (pinned) survive; "hello" and "" are dropped
	keys := map[string]bool{}
	for _, e := range out {
		keys[e.Key] = true
	}
	if !keys["SHORT"] {
		t.Error("expected SHORT to be kept")
	}
	if !keys["LONG"] {
		t.Error("expected LONG to be kept via KeepKeys")
	}
}

func TestShrink_DoesNotMutateOriginal(t *testing.T) {
	orig := entries()
	opts := shrinker.DefaultOptions()
	opts.MinLen = 10
	_ = shrinker.Shrink(orig, opts)
	if len(orig) != 4 {
		t.Error("original slice was mutated")
	}
}
