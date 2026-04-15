// Package differ provides utilities for comparing two sets of env entries
// and producing structured diff results.
package differ

import "github.com/your-org/envoy-cli/internal/envfile"

// ChangeType describes the kind of change detected between two env files.
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

// Result holds the full comparison output.
type Result struct {
	Changes []Change
}

// HasDiff returns true if any non-unchanged entries exist.
func (r Result) HasDiff() bool {
	for _, c := range r.Changes {
		if c.Type != Unchanged {
			return true
		}
	}
	return false
}

// Compare compares two slices of envfile.Entry and returns a Result.
func Compare(base, target []envfile.Entry) Result {
	baseMap := toMap(base)
	targetMap := toMap(target)

	seen := map[string]bool{}
	var changes []Change

	for _, e := range base {
		seen[e.Key] = true
		newVal, exists := targetMap[e.Key]
		switch {
		case !exists:
			changes = append(changes, Change{Key: e.Key, OldValue: e.Value, Type: Removed})
		case newVal != e.Value:
			changes = append(changes, Change{Key: e.Key, OldValue: e.Value, NewValue: newVal, Type: Changed})
		default:
			changes = append(changes, Change{Key: e.Key, OldValue: e.Value, NewValue: newVal, Type: Unchanged})
		}
	}

	for _, e := range target {
		if !seen[e.Key] {
			_ = baseMap
			changes = append(changes, Change{Key: e.Key, NewValue: e.Value, Type: Added})
		}
	}

	return Result{Changes: changes}
}

func toMap(entries []envfile.Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}
