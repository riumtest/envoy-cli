package differ_test

import (
	"testing"

	"github.com/envoy-cli/envoy-cli/internal/differ"
	"github.com/envoy-cli/envoy-cli/internal/envfile"
)

func TestCompare_NoChanges(t *testing.T) {
	base := []envfile.Entry{{Key: "FOO", Value: "bar"}}
	target := []envfile.Entry{{Key: "FOO", Value: "bar"}}
	result := differ.Compare(base, target)
	for _, c := range result.Changes {
		if c.Kind != differ.Unchanged {
			t.Errorf("expected unchanged, got %s", c.Kind)
		}
	}
}

func TestCompare_AddedKeys(t *testing.T) {
	base := []envfile.Entry{}
	target := []envfile.Entry{{Key: "NEW", Value: "val"}}
	result := differ.Compare(base, target)
	if len(result.Changes) != 1 || result.Changes[0].Kind != differ.Added {
		t.Error("expected one added change")
	}
}

func TestCompare_RemovedKeys(t *testing.T) {
	base := []envfile.Entry{{Key: "OLD", Value: "val"}}
	target := []envfile.Entry{}
	result := differ.Compare(base, target)
	if len(result.Changes) != 1 || result.Changes[0].Kind != differ.Removed {
		t.Error("expected one removed change")
	}
}

func TestCompare_ChangedKeys(t *testing.T) {
	base := []envfile.Entry{{Key: "FOO", Value: "old"}}
	target := []envfile.Entry{{Key: "FOO", Value: "new"}}
	result := differ.Compare(base, target)
	if len(result.Changes) != 1 || result.Changes[0].Kind != differ.Changed {
		t.Errorf("expected changed, got %v", result.Changes)
	}
	if result.Changes[0].OldValue != "old" || result.Changes[0].NewValue != "new" {
		t.Error("unexpected old/new values")
	}
}

func TestSummary(t *testing.T) {
	base := []envfile.Entry{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
	}
	target := []envfile.Entry{
		{Key: "A", Value: "changed"},
		{Key: "C", Value: "3"},
	}
	result := differ.Compare(base, target)
	summary := result.Summary()
	if summary[differ.Changed] != 1 {
		t.Errorf("expected 1 changed, got %d", summary[differ.Changed])
	}
	if summary[differ.Removed] != 1 {
		t.Errorf("expected 1 removed, got %d", summary[differ.Removed])
	}
	if summary[differ.Added] != 1 {
		t.Errorf("expected 1 added, got %d", summary[differ.Added])
	}
}
