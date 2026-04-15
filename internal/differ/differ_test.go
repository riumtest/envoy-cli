package differ_test

import (
	"strings"
	"testing"

	"github.com/envoy-cli/envoy-cli/internal/differ"
	"github.com/envoy-cli/envoy-cli/internal/envfile"
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
	target := mkEntries("FOO", "bar", "BAZ", "qux")
	res := differ.Compare(base, target)
	for _, c := range res.Changes {
		if c.Kind != differ.Unchanged {
			t.Errorf("expected unchanged, got %s for key %s", c.Kind, c.Key)
		}
	}
}

func TestCompare_AddedKeys(t *testing.T) {
	base := mkEntries("FOO", "bar")
	target := mkEntries("FOO", "bar", "NEW", "val")
	res := differ.Compare(base, target)
	found := false
	for _, c := range res.Changes {
		if c.Key == "NEW" && c.Kind == differ.Added {
			found = true
		}
	}
	if !found {
		t.Error("expected NEW to be marked as added")
	}
}

func TestCompare_RemovedKeys(t *testing.T) {
	base := mkEntries("FOO", "bar", "OLD", "gone")
	target := mkEntries("FOO", "bar")
	res := differ.Compare(base, target)
	found := false
	for _, c := range res.Changes {
		if c.Key == "OLD" && c.Kind == differ.Removed {
			found = true
		}
	}
	if !found {
		t.Error("expected OLD to be marked as removed")
	}
}

func TestCompare_ChangedKeys(t *testing.T) {
	base := mkEntries("FOO", "old")
	target := mkEntries("FOO", "new")
	res := differ.Compare(base, target)
	if len(res.Changes) != 1 || res.Changes[0].Kind != differ.Changed {
		t.Errorf("expected FOO to be changed, got %+v", res.Changes)
	}
	if res.Changes[0].OldValue != "old" || res.Changes[0].NewValue != "new" {
		t.Errorf("unexpected old/new values: %+v", res.Changes[0])
	}
}

func TestSummary(t *testing.T) {
	base := mkEntries("A", "1", "B", "2", "C", "3")
	target := mkEntries("A", "changed", "D", "new")
	res := differ.Compare(base, target)
	summary := res.Summary()
	if !strings.Contains(summary, "added=1") {
		t.Errorf("expected added=1 in summary, got: %s", summary)
	}
	if !strings.Contains(summary, "removed=2") {
		t.Errorf("expected removed=2 in summary, got: %s", summary)
	}
	if !strings.Contains(summary, "changed=1") {
		t.Errorf("expected changed=1 in summary, got: %s", summary)
	}
}
