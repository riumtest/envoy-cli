package deduplicator_test

import (
	"testing"

	"github.com/yourusername/envoy-cli/internal/deduplicator"
	"github.com/yourusername/envoy-cli/internal/envfile"
)

func entries(pairs ...string) []envfile.Entry {
	out := make([]envfile.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, envfile.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestDeduplicate_NoDuplicates(t *testing.T) {
	input := entries("FOO", "1", "BAR", "2")
	res := deduplicator.Deduplicate(input, deduplicator.KeepFirst)
	if len(res.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(res.Entries))
	}
	if len(res.Duplicates) != 0 {
		t.Fatalf("expected no duplicates, got %d", len(res.Duplicates))
	}
}

func TestDeduplicate_KeepFirst(t *testing.T) {
	input := entries("FOO", "first", "BAR", "bar", "FOO", "second")
	res := deduplicator.Deduplicate(input, deduplicator.KeepFirst)
	if len(res.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(res.Entries))
	}
	if res.Entries[0].Value != "first" {
		t.Errorf("expected 'first', got %q", res.Entries[0].Value)
	}
	if len(res.Duplicates) != 1 || res.Duplicates[0].Key != "FOO" {
		t.Errorf("expected duplicate report for FOO")
	}
}

func TestDeduplicate_KeepLast(t *testing.T) {
	input := entries("FOO", "first", "BAR", "bar", "FOO", "second")
	res := deduplicator.Deduplicate(input, deduplicator.KeepLast)
	if len(res.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(res.Entries))
	}
	if res.Entries[0].Value != "second" {
		t.Errorf("expected 'second', got %q", res.Entries[0].Value)
	}
}

func TestDeduplicate_DuplicateCount(t *testing.T) {
	input := entries("FOO", "a", "FOO", "b", "FOO", "c")
	res := deduplicator.Deduplicate(input, deduplicator.KeepFirst)
	if len(res.Duplicates) != 1 {
		t.Fatalf("expected 1 duplicate entry, got %d", len(res.Duplicates))
	}
	if res.Duplicates[0].Count != 3 {
		t.Errorf("expected count 3, got %d", res.Duplicates[0].Count)
	}
}

func TestDeduplicate_PreservesOrder(t *testing.T) {
	input := entries("C", "3", "A", "1", "B", "2", "A", "99")
	res := deduplicator.Deduplicate(input, deduplicator.KeepFirst)
	keys := []string{"C", "A", "B"}
	for i, e := range res.Entries {
		if e.Key != keys[i] {
			t.Errorf("position %d: expected key %q, got %q", i, keys[i], e.Key)
		}
	}
}

func TestDeduplicate_EmptyInput(t *testing.T) {
	res := deduplicator.Deduplicate([]envfile.Entry{}, deduplicator.KeepFirst)
	if len(res.Entries) != 0 || len(res.Duplicates) != 0 {
		t.Error("expected empty result for empty input")
	}
}
