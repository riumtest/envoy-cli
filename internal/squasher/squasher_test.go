package squasher_test

import (
	"testing"

	"github.com/user/envoy-cli/internal/envfile"
	"github.com/user/envoy-cli/internal/squasher"
)

func entries(kvs ...string) []envfile.Entry {
	out := make([]envfile.Entry, 0, len(kvs)/2)
	for i := 0; i+1 < len(kvs); i += 2 {
		out = append(out, envfile.Entry{Key: kvs[i], Value: kvs[i+1]})
	}
	return out
}

func TestSquash_NoDuplicates(t *testing.T) {
	in := entries("A", "1", "B", "2")
	res, err := squasher.Squash(in, squasher.DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Entries) != 2 || res.Squashed != 0 {
		t.Fatalf("expected 2 entries and 0 squashed, got %d/%d", len(res.Entries), res.Squashed)
	}
}

func TestSquash_StrategyLast(t *testing.T) {
	in := entries("A", "first", "B", "only", "A", "last")
	res, err := squasher.Squash(in, squasher.Options{Strategy: squasher.StrategyLast})
	if err != nil {
		t.Fatal(err)
	}
	if res.Squashed != 1 {
		t.Fatalf("expected 1 squashed, got %d", res.Squashed)
	}
	if res.Entries[0].Value != "last" {
		t.Fatalf("expected 'last', got %q", res.Entries[0].Value)
	}
}

func TestSquash_StrategyFirst(t *testing.T) {
	in := entries("A", "first", "A", "second")
	res, err := squasher.Squash(in, squasher.Options{Strategy: squasher.StrategyFirst})
	if err != nil {
		t.Fatal(err)
	}
	if res.Entries[0].Value != "first" {
		t.Fatalf("expected 'first', got %q", res.Entries[0].Value)
	}
	if res.Squashed != 1 {
		t.Fatalf("expected 1 squashed, got %d", res.Squashed)
	}
}

func TestSquash_StrategyError(t *testing.T) {
	in := entries("A", "1", "A", "2")
	_, err := squasher.Squash(in, squasher.Options{Strategy: squasher.StrategyError})
	if err == nil {
		t.Fatal("expected error for duplicate key")
	}
}

func TestSquash_DoesNotMutateOriginal(t *testing.T) {
	in := entries("X", "a", "X", "b")
	copy := make([]envfile.Entry, len(in))
	for i, e := range in {
		copy[i] = e
	}
	squasher.Squash(in, squasher.DefaultOptions()) //nolint
	for i, e := range in {
		if e != copy[i] {
			t.Fatal("original slice was mutated")
		}
	}
}
