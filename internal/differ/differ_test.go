package differ_test

import (
	"testing"

	"github.com/envoy-cli/internal/differ"
	"github.com/envoy-cli/internal/envfile"
)

func mkEntries(kvs ...string) []envfile.Entry {
	var out []envfile.Entry
	for i := 0; i+1 < len(kvs); i += 2 {
		out = append(out, envfile.Entry{Key: kvs[i], Value: kvs[i+1]})
	}
	return out
}

func TestCompare_NoChanges(t *testing.T) {
	a := mkEntries("FOO", "bar", "BAZ", "qux")
	b := mkEntries("FOO", "bar", "BAZ", "qux")
	result := differ.Compare(a, b)
	if len(result) != 0 {
		t.Fatalf("expected no changes, got %d", len(result))
	}
}

func TestCompare_AddedKeys(t *testing.T) {
	a := mkEntries("FOO", "bar")
	b := mkEntries("FOO", "bar", "NEW", "val")
	result := differ.Compare(a, b)
	if len(result) != 1 || result[0].Op != differ.OpAdded {
		t.Fatalf("expected 1 added, got %+v", result)
	}
}

func TestCompare_RemovedKeys(t *testing.T) {
	a := mkEntries("FOO", "bar", "OLD", "val")
	b := mkEntries("FOO", "bar")
	result := differ.Compare(a, b)
	if len(result) != 1 || result[0].Op != differ.OpRemoved {
		t.Fatalf("expected 1 removed, got %+v", result)
	}
}

func TestCompare_ChangedKeys(t *testing.T) {
	a := mkEntries("FOO", "old")
	b := mkEntries("FOO", "new")
	result := differ.Compare(a, b)
	if len(result) != 1 || result[0].Op != differ.OpChanged {
		t.Fatalf("expected 1 changed, got %+v", result)
	}
	if result[0].OldValue != "old" || result[0].NewValue != "new" {
		t.Fatalf("unexpected values: %+v", result[0])
	}
}
