// Package differ provides functionality for comparing two sets of env entries.
package differ

import (
	"fmt"

	"github.com/envoy-cli/envoy-cli/internal/envfile"
)

// ChangeKind represents the type of change between two env files.
type ChangeKind string

const (
	Added   ChangeKind = "added"
	Removed ChangeKind = "removed"
	Changed ChangeKind = "changed"
	Unchanged ChangeKind = "unchanged"
)

// Change represents a single key-level difference.
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

// Summary returns a human-readable summary of the diff result.
func (r *Result) Summary() string {
	added, removed, changed := 0, 0, 0
	for _, c := range r.Changes {
		switch c.Kind {
		case Added:
			added++
		case Removed:
			removed++
		case Changed:
			changed++
		}
	}
	return fmt.Sprintf("added=%d removed=%d changed=%d", added, removed, changed)
}

// Compare diffs two slices of env entries, returning a Result.
func Compare(base, target []envfile.Entry) *Result {
	baseMap := toMap(base)
	targetMap := toMap(target)

	var changes []Change

	for k, bv := range baseMap {
		if tv, ok := targetMap[k]; ok {
			if bv != tv {
				changes = append(changes, Change{Key: k, Kind: Changed, OldValue: bv, NewValue: tv})
			} else {
				changes = append(changes, Change{Key: k, Kind: Unchanged, OldValue: bv, NewValue: tv})
			}
		} else {
			changes = append(changes, Change{Key: k, Kind: Removed, OldValue: bv})
		}
	}

	for k, tv := range targetMap {
		if _, ok := baseMap[k]; !ok {
			changes = append(changes, Change{Key: k, Kind: Added, NewValue: tv})
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
