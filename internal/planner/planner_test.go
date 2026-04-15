package planner_test

import (
	"testing"

	"github.com/envoy-cli/envoy-cli/internal/differ"
	"github.com/envoy-cli/envoy-cli/internal/envfile"
	"github.com/envoy-cli/envoy-cli/internal/planner"
)

func entries(kv ...string) []envfile.Entry {
	out := make([]envfile.Entry, 0, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		out = append(out, envfile.Entry{Key: kv[i], Value: kv[i+1]})
	}
	return out
}

func TestBuild_NoChanges(t *testing.T) {
	base := entries("FOO", "bar", "BAZ", "qux")
	result := differ.Compare(base, base)
	plan := planner.Build(result)
	if plan.HasChanges() {
		t.Error("expected no changes")
	}
}

func TestBuild_AddedKeys(t *testing.T) {
	base := entries("FOO", "bar")
	target := entries("FOO", "bar", "NEW", "val")
	result := differ.Compare(base, target)
	plan := planner.Build(result)
	if !plan.HasChanges() {
		t.Fatal("expected changes")
	}
	found := false
	for _, a := range plan.Actions {
		if a.Key == "NEW" && a.Kind == planner.ActionSet {
			found = true
		}
	}
	if !found {
		t.Error("expected ActionSet for NEW")
	}
}

func TestBuild_RemovedKeys(t *testing.T) {
	base := entries("FOO", "bar", "OLD", "val")
	target := entries("FOO", "bar")
	result := differ.Compare(base, target)
	plan := planner.Build(result)
	for _, a := range plan.Actions {
		if a.Key == "OLD" && a.Kind != planner.ActionDelete {
			t.Errorf("expected ActionDelete for OLD, got %s", a.Kind)
		}
	}
}

func TestBuild_UpdatedKeys(t *testing.T) {
	base := entries("FOO", "old")
	target := entries("FOO", "new")
	result := differ.Compare(base, target)
	plan := planner.Build(result)
	for _, a := range plan.Actions {
		if a.Key == "FOO" {
			if a.Kind != planner.ActionUpdate {
				t.Errorf("expected ActionUpdate, got %s", a.Kind)
			}
			if a.OldValue != "old" || a.NewValue != "new" {
				t.Errorf("unexpected values: %s -> %s", a.OldValue, a.NewValue)
			}
		}
	}
}

func TestBuild_ActionString(t *testing.T) {
	a := planner.Action{Kind: planner.ActionSet, Key: "X", NewValue: "1"}
	if s := a.String(); s != "SET X=1" {
		t.Errorf("unexpected string: %s", s)
	}
	b := planner.Action{Kind: planner.ActionDelete, Key: "Y", OldValue: "2"}
	if s := b.String(); s != "DELETE Y" {
		t.Errorf("unexpected string: %s", s)
	}
	c := planner.Action{Kind: planner.ActionUpdate, Key: "Z", OldValue: "a", NewValue: "b"}
	if s := c.String(); s != "UPDATE Z: a -> b" {
		t.Errorf("unexpected string: %s", s)
	}
}
