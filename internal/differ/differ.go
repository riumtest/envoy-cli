package differ

import (
	"github.com/yourusername/envoy-cli/internal/envfile"
)

// ChangeType represents the kind of change detected between two env files.
type ChangeType string

const (
	Added   ChangeType = "added"
	Removed ChangeType = "removed"
	Changed ChangeType = "changed"
	Unchanged ChangeType = "unchanged"
)

// Change describes a single key-level difference.
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
	summary := map[ChangeType]int{
		Added:     0,
		Removed:   0,
		Changed:   0,
		Unchanged: 0,
	}
	for _, c := range r.Changes {
		summary[c.Type]++
	}
	return summary
}

// Compare diffs two slices of env entries and returns a Result.
func Compare(base, target []envfile.Entry) *Result {
	baseMap := toMap(base)
	targetMap := toMap(target)

	var changes []Change

	for _, e := range base {
		if e.Key == "" {
			continue
		}
		newVal, exists := targetMap[e.Key]
		if !exists {
			changes = append(changes, Change{Key: e.Key, OldValue: e.Value, Type: Removed})
		} else if newVal != e.Value {
			changes = append(changes, Change{Key: e.Key, OldValue: e.Value, NewValue: newVal, Type: Changed})
		} else {
			changes = append(changes, Change{Key: e.Key, OldValue: e.Value, NewValue: newVal, Type: Unchanged})
		}
	}

	for _, e := range target {
		if e.Key == "" {
			continue
		}
		if _, exists := baseMap[e.Key]; !exists {
			changes = append(changes, Change{Key: e.Key, NewValue: e.Value, Type: Added})
		}
	}

	return &Result{Changes: changes}
}

func toMap(entries []envfile.Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		if e.Key != "" {
			m[e.Key] = e.Value
		}
	}
	return m
}
