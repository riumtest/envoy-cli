package differ_test

import (
	"testing"

	"github.com/yourusername/envoy-cli/internal/differ"
	"github.com/yourusername/envoy-cli/internal/envfile"
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
		t.Errorf("expected no changes")
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
	base := []envfile.Entry{{Key: "FOO", Value: "bar"}}
	target := []envfile.Entry{{Key: "FOO", Value: "baz"}}

	result := differ.Compare(base, target)

	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result.Changes))
	}
	if result.Changes[0].Type != differ.Changed {
		t.Errorf("expected changed, got %s", result.Changes[0].Type)
	}
	if result.Changes[0].OldValue != "bar" || result.Changes[0].NewValue != "baz" {
		t.Errorf("unexpected old/new values")
	}
}

func TestSummary(t *testing.T) {
	base := []envfile.Entry{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
		{Key: "C", Value: "3"},
	}
	target := []envfile.Entry{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "changed"},
		{Key: "D", Value: "4"},
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
