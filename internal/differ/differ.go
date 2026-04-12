package differ

import (
	"github.com/envoy-cli/envoy-cli/internal/envfile"
)

// ChangeType represents the kind of change between two env files.
type ChangeType string

const (
	Added   ChangeType = "added"
	Removed ChangeType = "removed"
	Changed ChangeType = "changed"
	Unchanged ChangeType = "unchanged"
)

// Change represents a single key-level difference.
type Change struct {
	Key      string
	OldValue string
	NewValue string
	Type     ChangeType
}

// Result holds the full diff output.
type Result struct {
	Changes []Change
}

// Summary returns counts of each change type.
func (r *Result) Summary() map[ChangeType]int {
	counts := map[ChangeType]int{
		Added:     0,
		Removed:   0,
		Changed:   0,
		Unchanged: 0,
	}
	for _, c := range r.Changes {
		counts[c.Type]++
	}
	return counts
}

// Compare diffs two slices of envfile.Entry and returns a Result.
func Compare(base, target []envfile.Entry) *Result {
	baseMap := toMap(base)
	targetMap := toMap(target)

	seen := map[string]bool{}
	var changes []Change

	for _, e := range base {
		seen[e.Key] = true
		if newVal, ok := targetMap[e.Key]; ok {
			if newVal != e.Value {
				changes = append(changes, Change{Key: e.Key, OldValue: e.Value, NewValue: newVal, Type: Changed})
			} else {
				changes = append(changes, Change{Key: e.Key, OldValue: e.Value, NewValue: newVal, Type: Unchanged})
			}
		} else {
			changes = append(changes, Change{Key: e.Key, OldValue: e.Value, NewValue: "", Type: Removed})
		}
	}

	for _, e := range target {
		if !seen[e.Key] {
			if _, exists := baseMap[e.Key]; !exists {
				changes = append(changes, Change{Key: e.Key, OldValue: "", NewValue: e.Value, Type: Added})
			}
		}
	}

	return &Result{Changes: changes}
}

func toMap(entries []envfile.Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}
