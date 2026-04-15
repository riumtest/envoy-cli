// Package differ compares two sets of env entries and produces structured diffs.
package differ

import (
	"sort"

	"envoy-cli/internal/envfile"
)

// ChangeType represents the kind of change detected between two env files.
type ChangeType string

const (
	Added   ChangeType = "added"
	Removed ChangeType = "removed"
	Changed ChangeType = "changed"
	Unchanged ChangeType = "unchanged"
)

// Diff represents a single key-level difference.
type Diff struct {
	Key      string
	OldValue string
	NewValue string
	Change   ChangeType
}

// Result holds the full comparison output.
type Result struct {
	Diffs []Diff
}

// HasChanges returns true if any non-unchanged diffs exist.
func (r Result) HasChanges() bool {
	for _, d := range r.Diffs {
		if d.Change != Unchanged {
			return true
		}
	}
	return false
}

// Compare computes the diff between a base and head set of env entries.
func Compare(base, head []envfile.Entry) Result {
	baseMap := toMap(base)
	headMap := toMap(head)

	keys := make(map[string]struct{})
	for k := range baseMap {
		keys[k] = struct{}{}
	}
	for k := range headMap {
		keys[k] = struct{}{}
	}

	sorted := make([]string, 0, len(keys))
	for k := range keys {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)

	var diffs []Diff
	for _, k := range sorted {
		oldVal, inBase := baseMap[k]
		newVal, inHead := headMap[k]

		switch {
		case inBase && !inHead:
			diffs = append(diffs, Diff{Key: k, OldValue: oldVal, Change: Removed})
		case !inBase && inHead:
			diffs = append(diffs, Diff{Key: k, NewValue: newVal, Change: Added})
		case oldVal != newVal:
			diffs = append(diffs, Diff{Key: k, OldValue: oldVal, NewValue: newVal, Change: Changed})
		default:
			diffs = append(diffs, Diff{Key: k, OldValue: oldVal, NewValue: newVal, Change: Unchanged})
		}
	}

	return Result{Diffs: diffs}
}

func toMap(entries []envfile.Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}
