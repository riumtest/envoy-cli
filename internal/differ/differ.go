package differ

import (
	"github.com/envoy-cli/envoy-cli/internal/envfile"
)

// DiffKind describes the type of change between two env files.
type DiffKind string

const (
	Added   DiffKind = "added"
	Removed DiffKind = "removed"
	Changed DiffKind = "changed"
	Unchanged DiffKind = "unchanged"
)

// Change represents a single key-level difference.
type Change struct {
	Key      string
	Kind     DiffKind
	OldValue string
	NewValue string
}

// Result holds the full diff between two env files.
type Result struct {
	Changes []Change
}

// Summary returns counts of each change kind.
func (r *Result) Summary() map[DiffKind]int {
	m := map[DiffKind]int{}
	for _, c := range r.Changes {
		m[c.Kind]++
	}
	return m
}

// Compare diffs two slices of envfile.Entry, returning a Result.
func Compare(base, target []envfile.Entry) Result {
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

	return Result{Changes: changes}
}

func toMap(entries []envfile.Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}
