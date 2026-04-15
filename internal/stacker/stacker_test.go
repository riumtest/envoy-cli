package stacker_test

import (
	"testing"

	"github.com/user/envoy-cli/internal/envfile"
	"github.com/user/envoy-cli/internal/stacker"
)

func entries(pairs ...string) []envfile.Entry {
	out := make([]envfile.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, envfile.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestStack_NoConflicts(t *testing.T) {
	base := entries("A", "1", "B", "2")
	overlay := entries("C", "3")
	res := stacker.Stack([][]envfile.Entry{base, overlay}, stacker.DefaultOptions())
	if len(res.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(res.Entries))
	}
	if res.Overridden != 0 {
		t.Errorf("expected 0 overrides, got %d", res.Overridden)
	}
}

func TestStack_StrategyLast(t *testing.T) {
	base := entries("A", "base", "B", "base")
	overlay := entries("A", "overlay")
	res := stacker.Stack([][]envfile.Entry{base, overlay}, stacker.DefaultOptions())
	for _, e := range res.Entries {
		if e.Key == "A" && e.Value != "overlay" {
			t.Errorf("expected overlay value, got %q", e.Value)
		}
	}
	if res.Overridden != 1 {
		t.Errorf("expected 1 override, got %d", res.Overridden)
	}
}

func TestStack_StrategyFirst(t *testing.T) {
	base := entries("A", "base")
	overlay := entries("A", "overlay", "B", "new")
	opts := stacker.Options{Strategy: stacker.StrategyFirst}
	res := stacker.Stack([][]envfile.Entry{base, overlay}, opts)
	for _, e := range res.Entries {
		if e.Key == "A" && e.Value != "base" {
			t.Errorf("expected base value, got %q", e.Value)
		}
	}
	if res.Overridden != 1 {
		t.Errorf("expected 1 override, got %d", res.Overridden)
	}
}

func TestStack_PreservesOrder(t *testing.T) {
	base := entries("Z", "1", "A", "2")
	overlay := entries("M", "3")
	res := stacker.Stack([][]envfile.Entry{base, overlay}, stacker.DefaultOptions())
	keys := []string{"Z", "A", "M"}
	for i, e := range res.Entries {
		if e.Key != keys[i] {
			t.Errorf("position %d: expected %q got %q", i, keys[i], e.Key)
		}
	}
}

func TestStack_EmptyLayers(t *testing.T) {
	res := stacker.Stack([][]envfile.Entry{}, stacker.DefaultOptions())
	if len(res.Entries) != 0 {
		t.Errorf("expected empty result")
	}
}

func TestStack_SkipsEmptyKeys(t *testing.T) {
	layer := []envfile.Entry{{Key: "", Value: "orphan"}, {Key: "VALID", Value: "ok"}}
	res := stacker.Stack([][]envfile.Entry{layer}, stacker.DefaultOptions())
	if len(res.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(res.Entries))
	}
}
