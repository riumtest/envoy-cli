// Package differ provides utilities for comparing two sets of env entries.
package differ

import (
	"fmt"

	"github.com/envoy-cli/internal/envfile"
)

// ChangeType describes the kind of change detected for a key.
type ChangeType string

const (
	Added   ChangeType = "added"
	Removed ChangeType = "removed"
	Changed ChangeType = "changed"
	Unchanged ChangeType = "unchanged"
)

// Diff represents a single key-level difference between two env files.
type Diff struct {
	Key      string
	OldValue string
	NewValue string
	Type     ChangeType
}

// Result holds the full comparison output.
type Result struct {
	Diffs []Diff
}

// Summary returns a human-readable summary string.
func (r *Result) Summary() string {
	var added, removed, changed int
	for _, d := range r.Diffs {
		switch d.Type {
		case Added:
			added++
		case Removed:
			removed++
		case Changed:
			changed++
		}
	}
	return fmt.Sprintf("%d added, %d removed, %d changed", added, removed, changed)
}

// Compare compares two slices of env entries and returns a Result.
func Compare(base, target []envfile.Entry) *Result {
	baseMap := toMap(base)
	targetMap := toMap(target)

	var diffs []Diff

	for _, e := range base {
		if newVal, ok := targetMap[e.Key]; ok {
			if newVal != e.Value {
				diffs = append(diffs, Diff{Key: e.Key, OldValue: e.Value, NewValue: newVal, Type: Changed})
			} else {
				diffs = append(diffs, Diff{Key: e.Key, OldValue: e.Value, NewValue: newVal, Type: Unchanged})
			}
		} else {
			diffs = append(diffs, Diff{Key: e.Key, OldValue: e.Value, NewValue: "", Type: Removed})
		}
	}

	for _, e := range target {
		if _, ok := baseMap[e.Key]; !ok {
			diffs = append(diffs, Diff{Key: e.Key, OldValue: "", NewValue: e.Value, Type: Added})
		}
	}

	return &Result{Diffs: diffs}
}

func toMap(entries []envfile.Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}
