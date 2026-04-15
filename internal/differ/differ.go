// Package differ provides functionality for comparing two sets of env entries
// and producing structured diffs.
package differ

import (
	"fmt"

	"github.com/envoy-cli/envoy/internal/envfile"
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

// Result holds the full diff outcome.
type Result struct {
	Changes []Change
}

// Summary returns a human-readable one-line summary of the diff result.
func (r *Result) Summary() string {
	var added, removed, changed int
	for _, c := range r.Changes {
		switch c.Type {
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

// Compare diffs two slices of env entries and returns a Result.
func Compare(base, target []envfile.Entry) *Result {
	baseMap := toMap(base)
	targetMap := toMap(target)

	var changes []Change

	for k, bv := range baseMap {
		if tv, ok := targetMap[k]; ok {
			if bv != tv {
				changes = append(changes, Change{Key: k, OldValue: bv, NewValue: tv, Type: Changed})
			} else {
				changes = append(changes, Change{Key: k, OldValue: bv, NewValue: tv, Type: Unchanged})
			}
		} else {
			changes = append(changes, Change{Key: k, OldValue: bv, NewValue: "", Type: Removed})
		}
	}

	for k, tv := range targetMap {
		if _, ok := baseMap[k]; !ok {
			changes = append(changes, Change{Key: k, OldValue: "", NewValue: tv, Type: Added})
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
