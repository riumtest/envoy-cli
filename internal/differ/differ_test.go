package differ_test

import (
	"testing"

	"github.com/yourusername/envoy-cli/internal/differ"
	"github.com/yourusername/envoy-cli/internal/envfile"
)

func TestCompare_NoChanges(t *testing.T) {
	base := []envfile.Entry{{Key: "FOO", Value: "bar"}}
	target := []envfile.Entry{{Key: "FOO", Value: "bar"}}

	result := differ.Compare(base, target)
	summary := result.Summary()

	if summary[differ.Unchanged] != 1 {
		t.Errorf("expected 1 unchanged, got %d", summary[differ.Unchanged])
	}
	if summary[differ.Changed]+summary[differ.Added]+summary[differ.Removed] != 0 {
		t.Error("expected no other changes")
	}
}

func TestCompare_AddedKeys(t *testing.T) {
	base := []envfile.Entry{}
	target := []envfile.Entry{{Key: "NEW_KEY", Value: "value"}}

	result := differ.Compare(base, target)
	summary := result.Summary()

	if summary[differ.Added] != 1 {
		t.Errorf("expected 1 added, got %d", summary[differ.Added])
	}
}

func TestCompare_RemovedKeys(t *testing.T) {
	base := []envfile.Entry{{Key: "OLD_KEY", Value: "value"}}
	target := []envfile.Entry{}

	result := differ.Compare(base, target)
	summary := result.Summary()

	if summary[differ.Removed] != 1 {
		t.Errorf("expected 1 removed, got %d", summary[differ.Removed])
	}
}

func TestCompare_ChangedKeys(t *testing.T) {
	base := []envfile.Entry{{Key: "FOO", Value: "old"}}
	target := []envfile.Entry{{Key: "FOO", Value: "new"}}

	result := differ.Compare(base, target)
	summary := result.Summary()

	if summary[differ.Changed] != 1 {
		t.Errorf("expected 1 changed, got %d", summary[differ.Changed])
	}
}

func TestSummary(t *testing.T) {
	base := []envfile.Entry{
		{Key: "KEEP", Value: "same"},
		{Key: "MODIFY", Value: "old"},
		{Key: "REMOVE", Value: "gone"},
	}
	target := []envfile.Entry{
		{Key: "KEEP", Value: "same"},
		{Key: "MODIFY", Value: "new"},
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
