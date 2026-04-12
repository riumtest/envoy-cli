package differ_test

import (
	"testing"

	"github.com/envoy-cli/envoy-cli/internal/differ"
	"github.com/envoy-cli/envoy-cli/internal/envfile"
)

func TestCompare_NoChanges(t *testing.T) {
	base := []envfile.Entry{{Key: "FOO", Value: "bar"}, {Key: "BAZ", Value: "qux"}}
	target := []envfile.Entry{{Key: "FOO", Value: "bar"}, {Key: "BAZ", Value: "qux"}}

	result := differ.Compare(base, target)
	summary := result.Summary()

	if summary[differ.Unchanged] != 2 {
		t.Errorf("expected 2 unchanged, got %d", summary[differ.Unchanged])
	}
	if summary[differ.Added]+summary[differ.Removed]+summary[differ.Changed] != 0 {
		t.Error("expected no additions, removals, or changes")
	}
}

func TestCompare_AddedKeys(t *testing.T) {
	base := []envfile.Entry{{Key: "FOO", Value: "bar"}}
	target := []envfile.Entry{{Key: "FOO", Value: "bar"}, {Key: "NEW_KEY", Value: "new"}}

	result := differ.Compare(base, target)
	summary := result.Summary()

	if summary[differ.Added] != 1 {
		t.Errorf("expected 1 added, got %d", summary[differ.Added])
	}
}

func TestCompare_RemovedKeys(t *testing.T) {
	base := []envfile.Entry{{Key: "FOO", Value: "bar"}, {Key: "OLD_KEY", Value: "old"}}
	target := []envfile.Entry{{Key: "FOO", Value: "bar"}}

	result := differ.Compare(base, target)
	summary := result.Summary()

	if summary[differ.Removed] != 1 {
		t.Errorf("expected 1 removed, got %d", summary[differ.Removed])
	}
}

func TestCompare_ChangedKeys(t *testing.T) {
	base := []envfile.Entry{{Key: "FOO", Value: "old_value"}}
	target := []envfile.Entry{{Key: "FOO", Value: "new_value"}}

	result := differ.Compare(base, target)
	summary := result.Summary()

	if summary[differ.Changed] != 1 {
		t.Errorf("expected 1 changed, got %d", summary[differ.Changed])
	}
	if result.Changes[0].OldValue != "old_value" {
		t.Errorf("expected OldValue=old_value, got %s", result.Changes[0].OldValue)
	}
	if result.Changes[0].NewValue != "new_value" {
		t.Errorf("expected NewValue=new_value, got %s", result.Changes[0].NewValue)
	}
}

func TestSummary(t *testing.T) {
	base := []envfile.Entry{
		{Key: "KEEP", Value: "same"},
		{Key: "CHANGE", Value: "old"},
		{Key: "REMOVE", Value: "gone"},
	}
	target := []envfile.Entry{
		{Key: "KEEP", Value: "same"},
		{Key: "CHANGE", Value: "new"},
		{Key: "ADD", Value: "fresh"},
	}

	result := differ.Compare(base, target)
	summary := result.Summary()

	if summary[differ.Unchanged] != 1 {
		t.Errorf("expected 1 unchanged, got %d", summary[differ.Unchanged])
	}
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
