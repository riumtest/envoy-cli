package differ_test

import (
	"testing"

	"github.com/envoy-cli/internal/differ"
	"github.com/envoy-cli/internal/envfile"
)

func TestCompare_NoChanges(t *testing.T) {
	base := []envfile.Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "BAZ", Value: "qux"},
	}
	result := differ.Compare(base, base)
	for _, d := range result.Diffs {
		if d.Type != differ.Unchanged {
			t.Errorf("expected unchanged, got %s for key %s", d.Type, d.Key)
		}
	}
}

func TestCompare_AddedKeys(t *testing.T) {
	base := []envfile.Entry{{Key: "FOO", Value: "bar"}}
	target := []envfile.Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "NEW", Value: "value"},
	}
	result := differ.Compare(base, target)
	found := false
	for _, d := range result.Diffs {
		if d.Key == "NEW" && d.Type == differ.Added {
			found = true
		}
	}
	if !found {
		t.Error("expected NEW key to be marked as added")
	}
}

func TestCompare_RemovedKeys(t *testing.T) {
	base := []envfile.Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "OLD", Value: "gone"},
	}
	target := []envfile.Entry{{Key: "FOO", Value: "bar"}}
	result := differ.Compare(base, target)
	found := false
	for _, d := range result.Diffs {
		if d.Key == "OLD" && d.Type == differ.Removed {
			found = true
		}
	}
	if !found {
		t.Error("expected OLD key to be marked as removed")
	}
}

func TestCompare_ChangedKeys(t *testing.T) {
	base := []envfile.Entry{{Key: "FOO", Value: "old"}}
	target := []envfile.Entry{{Key: "FOO", Value: "new"}}
	result := differ.Compare(base, target)
	if len(result.Diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(result.Diffs))
	}
	if result.Diffs[0].Type != differ.Changed {
		t.Errorf("expected changed, got %s", result.Diffs[0].Type)
	}
	if result.Diffs[0].OldValue != "old" || result.Diffs[0].NewValue != "new" {
		t.Errorf("unexpected values: old=%s new=%s", result.Diffs[0].OldValue, result.Diffs[0].NewValue)
	}
}

func TestSummary(t *testing.T) {
	base := []envfile.Entry{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
	}
	target := []envfile.Entry{
		{Key: "A", Value: "changed"},
		{Key: "C", Value: "new"},
	}
	result := differ.Compare(base, target)
	summary := result.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
	expected := "1 added, 1 removed, 1 changed"
	if summary != expected {
		t.Errorf("expected %q, got %q", expected, summary)
	}
}
