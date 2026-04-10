package planner_test

import (
	"testing"

	"github.com/envoy-cli/envoy-cli/internal/envfile"
	"github.com/envoy-cli/envoy-cli/internal/planner"
)

func entries(pairs ...string) []envfile.Entry {
	var out []envfile.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, envfile.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestBuild_NoChanges(t *testing.T) {
	src := entries("FOO", "bar", "BAZ", "qux")
	plan := planner.Build(src, src)
	if plan.KeepCount() != 2 {
		t.Fatalf("expected 2 keep steps, got %d", plan.KeepCount())
	}
	if plan.AddCount()+plan.RemoveCount()+plan.UpdateCount() != 0 {
		t.Fatal("expected no mutations")
	}
}

func TestBuild_AddedKeys(t *testing.T) {
	src := entries("FOO", "bar")
	dst := entries("FOO", "bar", "NEW", "val")
	plan := planner.Build(src, dst)
	if plan.AddCount() != 1 {
		t.Fatalf("expected 1 add, got %d", plan.AddCount())
	}
	if plan.Steps[1].Key != "NEW" {
		t.Errorf("expected added key NEW, got %s", plan.Steps[1].Key)
	}
}

func TestBuild_RemovedKeys(t *testing.T) {
	src := entries("FOO", "bar", "OLD", "gone")
	dst := entries("FOO", "bar")
	plan := planner.Build(src, dst)
	if plan.RemoveCount() != 1 {
		t.Fatalf("expected 1 remove, got %d", plan.RemoveCount())
	}
}

func TestBuild_UpdatedKeys(t *testing.T) {
	src := entries("FOO", "old")
	dst := entries("FOO", "new")
	plan := planner.Build(src, dst)
	if plan.UpdateCount() != 1 {
		t.Fatalf("expected 1 update, got %d", plan.UpdateCount())
	}
	if plan.Steps[0].OldValue != "old" || plan.Steps[0].NewValue != "new" {
		t.Errorf("unexpected values: %+v", plan.Steps[0])
	}
}

func TestBuild_Mixed(t *testing.T) {
	src := entries("KEEP", "v", "UPDATE", "old", "REMOVE", "gone")
	dst := entries("KEEP", "v", "UPDATE", "new", "ADD", "fresh")
	plan := planner.Build(src, dst)
	if plan.KeepCount() != 1 {
		t.Errorf("keep: want 1, got %d", plan.KeepCount())
	}
	if plan.UpdateCount() != 1 {
		t.Errorf("update: want 1, got %d", plan.UpdateCount())
	}
	if plan.RemoveCount() != 1 {
		t.Errorf("remove: want 1, got %d", plan.RemoveCount())
	}
	if plan.AddCount() != 1 {
		t.Errorf("add: want 1, got %d", plan.AddCount())
	}
}

func TestBuild_EmptySrc(t *testing.T) {
	dst := entries("A", "1", "B", "2")
	plan := planner.Build(nil, dst)
	if plan.AddCount() != 2 {
		t.Fatalf("expected 2 adds, got %d", plan.AddCount())
	}
}

func TestBuild_EmptyDst(t *testing.T) {
	src := entries("A", "1", "B", "2")
	plan := planner.Build(src, nil)
	if plan.RemoveCount() != 2 {
		t.Fatalf("expected 2 removes, got %d", plan.RemoveCount())
	}
}
