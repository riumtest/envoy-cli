package differ_test

import (
	"testing"

	"github.com/your-org/envoy-cli/internal/differ"
	"github.com/your-org/envoy-cli/internal/envfile"
)

func mkEntries(pairs ...string) []envfile.Entry {
	var out []envfile.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, envfile.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestCompare_NoChanges(t *testing.T) {
	base := mkEntries("FOO", "bar", "BAZ", "qux")
	result := differ.Compare(base, base)
	if result.HasDiff() {
		t.Fatal("expected no diff")
	}
}

func TestCompare_AddedKeys(t *testing.T) {
	base := mkEntries("FOO", "bar")
	target := mkEntries("FOO", "bar", "NEW", "val")
	result := differ.Compare(base, target)
	if !result.HasDiff() {
		t.Fatal("expected diff")
	}
	found := false
	for _, c := range result.Changes {
		if c.Key == "NEW" && c.Type == differ.Added {
			found = true
		}
	}
	if !found {
		t.Error("expected Added change for NEW")
	}
}

func TestCompare_RemovedKeys(t *testing.T) {
	base := mkEntries("FOO", "bar", "OLD", "gone")
	target := mkEntries("FOO", "bar")
	result := differ.Compare(base, target)
	if !result.HasDiff() {
		t.Fatal("expected diff")
	}
	for _, c := range result.Changes {
		if c.Key == "OLD" && c.Type != differ.Removed {
			t.Errorf("expected Removed, got %s", c.Type)
		}
	}
}

func TestCompare_ChangedKeys(t *testing.T) {
	base := mkEntries("FOO", "old")
	target := mkEntries("FOO", "new")
	result := differ.Compare(base, target)
	if !result.HasDiff() {
		t.Fatal("expected diff")
	}
	for _, c := range result.Changes {
		if c.Key == "FOO" {
			if c.Type != differ.Changed {
				t.Errorf("expected Changed, got %s", c.Type)
			}
			if c.OldValue != "old" || c.NewValue != "new" {
				t.Errorf("unexpected values: %s -> %s", c.OldValue, c.NewValue)
			}
		}
	}
}

func TestCompare_MixedChanges(t *testing.T) {
	base := mkEntries("A", "1", "B", "2", "C", "3")
	target := mkEntries("A", "1", "B", "99", "D", "4")
	result := differ.Compare(base, target)
	types := map[string]differ.ChangeType{}
	for _, c := range result.Changes {
		types[c.Key] = c.Type
	}
	if types["A"] != differ.Unchanged {
		t.Errorf("A should be unchanged")
	}
	if types["B"] != differ.Changed {
		t.Errorf("B should be changed")
	}
	if types["C"] != differ.Removed {
		t.Errorf("C should be removed")
	}
	if types["D"] != differ.Added {
		t.Errorf("D should be added")
	}
}
