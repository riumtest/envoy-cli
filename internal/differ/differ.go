package differ

import (
	"github.com/yourusername/envoy-cli/internal/envfile"
)

// ChangeType represents the type of change between two env files.
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
	m := map[ChangeType]int{
		Added:     0,
		Removed:   0,
		Changed:   0,
		Unchanged: 0,
	}
	for _, c := range r.Changes {
		m[c.Type]++
	}
	return m
}

// Compare diffs two slices of envfile.Entry and returns a Result.
func Compare(base, target []envfile.Entry) *Result {
	baseMap := toMap(base)
	targetMap := toMap(target)

	var changes []Change

	for key, baseVal := range baseMap {
		if targetVal, ok := targetMap[key]; ok {
			if baseVal == targetVal {
				changes = append(changes, Change{Key: key, OldValue: baseVal, NewValue: targetVal, Type: Unchanged})
			} else {
				changes = append(changes, Change{Key: key, OldValue: baseVal, NewValue: targetVal, Type: Changed})
			}
		} else {
			changes = append(changes, Change{Key: key, OldValue: baseVal, NewValue: "", Type: Removed})
		}
	}

	for key, targetVal := range targetMap {
		if _, ok := baseMap[key]; !ok {
			changes = append(changes, Change{Key: key, OldValue: "", NewValue: targetVal, Type: Added})
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
