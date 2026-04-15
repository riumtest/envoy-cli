package differ

import (
	"testing"

	"github.com/envoy-cli/envoy/internal/envfile"
)

func TestCompare_NoChanges(t *testing.T) {
	entries := []envfile.Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "BAZ", Value: "qux"},
	}
	result := Compare(entries, entries)
	for _, c := range result.Changes {
		if c.Type != Unchanged {
			t.Errorf("expected unchanged, got %s for key %s", c.Type, c.Key)
		}
	}
}

func TestCompare_AddedKeys(t *testing.T) {
	base := []envfile.Entry{{Key: "FOO", Value: "bar"}}
	target := []envfile.Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "NEW_KEY", Value: "new"},
	}
	result := Compare(base, target)
	found := false
	for _, c := range result.Changes {
		if c.Key == "NEW_KEY" && c.Type == Added {
			found = true
		}
	}
	if !found {
		t.Error("expected NEW_KEY to be marked as added")
	}
}

func TestCompare_RemovedKeys(t *testing.T) {
	base := []envfile.Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "OLD_KEY", Value: "old"},
	}
	target := []envfile.Entry{{Key: "FOO", Value: "bar"}}
	result := Compare(base, target)
	found := false
	for _, c := range result.Changes {
		if c.Key == "OLD_KEY" && c.Type == Removed {
			found = true
		}
	}
	if !found {
		t.Error("expected OLD_KEY to be marked as removed")
	}
}

func TestCompare_ChangedKeys(t *testing.T) {
	base := []envfile.Entry{{Key: "FOO", Value: "old"}}
	target := []envfile.Entry{{Key: "FOO", Value: "new"}}
	result := Compare(base, target)
	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result.Changes))
	}
	c := result.Changes[0]
	if c.Type != Changed {
		t.Errorf("expected Changed, got %s", c.Type)
	}
	if c.OldValue != "old" || c.NewValue != "new" {
		t.Errorf("unexpected values: old=%s new=%s", c.OldValue, c.NewValue)
	}
}

func TestSummary(t *testing.T) {
	result := &Result{
		Changes: []Change{
			{Type: Added},
			{Type: Added},
			{Type: Removed},
			{Type: Changed},
			{Type: Unchanged},
		},
	}
	summary := result.Summary()
	expected := "2 added, 1 removed, 1 changed"
	if summary != expected {
		t.Errorf("expected %q, got %q", expected, summary)
	}
}
