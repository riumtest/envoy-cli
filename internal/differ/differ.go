// Package differ compares two sets of env entries and reports differences.
package differ

import (
	"github.com/user/envoy-cli/internal/envfile"
)

// ChangeKind describes the type of change detected between two env files.
type ChangeKind string

const (
	Added   ChangeKind = "added"
	Removed ChangeKind = "removed"
	Changed ChangeKind = "changed"
	Unchanged ChangeKind = "unchanged"
)

// Change represents a single key-level diff between two env files.
type Change struct {
	Key      string
	Kind     ChangeKind
	OldValue string
	NewValue string
}

// Result holds the full diff output.
type Result struct {
	Changes []Change
}

// HasDiff returns true if any non-unchanged entries exist.
func (r Result) HasDiff() bool {
	for _, c := range r.Changes {
		if c.Kind != Unchanged {
			return true
		}
	}
	return false
}

// Compare diffs two slices of env entries and returns a Result.
func Compare(base, target []envfile.Entry) Result {
	baseMap := toMap(base)
	targetMap := toMap(target)

	seen := map[string]bool{}
	var changes []Change

	for _, e := range base {
		if seen[e.Key] {
			continue
		}
		seen[e.Key] = true

		newVal, exists := targetMap[e.Key]
		switch {
		case !exists:
			changes = append(changes, Change{Key: e.Key, Kind: Removed, OldValue: e.Value})
		case newVal != e.Value:
			changes = append(changes, Change{Key: e.Key, Kind: Changed, OldValue: e.Value, NewValue: newVal})
		default:
			changes = append(changes, Change{Key: e.Key, Kind: Unchanged, OldValue: e.Value, NewValue: newVal})
		}
	}

	for _, e := range target {
		if seen[e.Key] {
			continue
		}
		seen[e.Key] = true
		if _, exists := baseMap[e.Key]; !exists {
			changes = append(changes, Change{Key: e.Key, Kind: Added, NewValue: e.Value})
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
