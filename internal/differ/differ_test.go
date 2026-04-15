package differ_test

import (
	"testing"

	"github.com/user/envoy-cli/internal/differ"
	"github.com/user/envoy-cli/internal/envfile"
)

func mkEntries(pairs ...string) []envfile.Entry {
	var out []envfile.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, envfile.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestCompare_NoChanges(t *testing.T) {
	base := mkEntries("A", "1", "B", "2")
	result := differ.Compare(base, base)
	if result.HasDiff() {
		t.Fatal("expected no diff")
	}
	for _, c := range result.Changes {
		if c.Kind != differ.Unchanged {
			t.Errorf("expected Unchanged, got %s for key %s", c.Kind, c.Key)
		}
	}
}

func TestCompare_AddedKeys(t *testing.T) {
	base := mkEntries("A", "1")
	target := mkEntries("A", "1", "B", "2")
	result := differ.Compare(base, target)
	if !result.HasDiff() {
		t.Fatal("expected diff")
	}
	found := false
	for _, c := range result.Changes {
		if c.Key == "B" && c.Kind == differ.Added {
			found = true
		}
	}
	if !found {
		t.Error("expected key B to be marked as Added")
	}
}

func TestCompare_RemovedKeys(t *testing.T) {
	base := mkEntries("A", "1", "B", "2")
	target := mkEntries("A", "1")
	result := differ.Compare(base, target)
	if !result.HasDiff() {
		t.Fatal("expected diff")
	}
	for _, c := range result.Changes {
		if c.Key == "B" && c.Kind != differ.Removed {
			t.Errorf("expected B to be Removed, got %s", c.Kind)
		}
	}
}

func TestCompare_ChangedKeys(t *testing.T) {
	base := mkEntries("A", "old")
	target := mkEntries("A", "new")
	result := differ.Compare(base, target)
	if !result.HasDiff() {
		t.Fatal("expected diff")
	}
	for _, c := range result.Changes {
		if c.Key == "A" {
			if c.Kind != differ.Changed {
				t.Errorf("expected Changed, got %s", c.Kind)
			}
			if c.OldValue != "old" || c.NewValue != "new" {
				t.Errorf("unexpected values: old=%s new=%s", c.OldValue, c.NewValue)
			}
		}
	}
}

func TestCompare_HasDiff_ReturnsFalseWhenClean(t *testing.T) {
	entries := mkEntries("X", "1", "Y", "2")
	result := differ.Compare(entries, entries)
	if result.HasDiff() {
		t.Error("HasDiff should return false for identical inputs")
	}
}
