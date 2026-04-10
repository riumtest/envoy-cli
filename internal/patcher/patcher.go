// Package patcher provides functionality to apply key-value patches
// to a set of env entries, supporting set, delete, and rename operations.
package patcher

import (
	"fmt"

	"github.com/user/envoy-cli/internal/envfile"
)

// Op represents a patch operation type.
type Op string

const (
	OpSet    Op = "set"
	OpDelete Op = "delete"
	OpRename Op = "rename"
)

// Patch describes a single mutation to apply to an env entry set.
type Patch struct {
	Op      Op
	Key     string
	Value   string // used by OpSet
	NewKey  string // used by OpRename
}

// Result holds the patched entries and a summary of applied changes.
type Result struct {
	Entries  []envfile.Entry
	Applied  []string
	Skipped  []string
}

// Apply applies a list of patches to the provided entries and returns a Result.
func Apply(entries []envfile.Entry, patches []Patch) (Result, error) {
	result := Result{}
	working := make([]envfile.Entry, len(entries))
	copy(working, entries)

	for _, p := range patches {
		switch p.Op {
		case OpSet:
			working = applySet(working, p.Key, p.Value)
			result.Applied = append(result.Applied, fmt.Sprintf("set %s", p.Key))
		case OpDelete:
			var removed bool
			working, removed = applyDelete(working, p.Key)
			if removed {
				result.Applied = append(result.Applied, fmt.Sprintf("delete %s", p.Key))
			} else {
				result.Skipped = append(result.Skipped, fmt.Sprintf("delete %s (not found)", p.Key))
			}
		case OpRename:
			var renamed bool
			working, renamed = applyRename(working, p.Key, p.NewKey)
			if renamed {
				result.Applied = append(result.Applied, fmt.Sprintf("rename %s -> %s", p.Key, p.NewKey))
			} else {
				result.Skipped = append(result.Skipped, fmt.Sprintf("rename %s (not found)", p.Key))
			}
		default:
			return Result{}, fmt.Errorf("unknown patch op: %q", p.Op)
		}
	}

	result.Entries = working
	return result, nil
}

func applySet(entries []envfile.Entry, key, value string) []envfile.Entry {
	for i, e := range entries {
		if e.Key == key {
			entries[i].Value = value
			return entries
		}
	}
	return append(entries, envfile.Entry{Key: key, Value: value})
}

func applyDelete(entries []envfile.Entry, key string) ([]envfile.Entry, bool) {
	for i, e := range entries {
		if e.Key == key {
			return append(entries[:i], entries[i+1:]...), true
		}
	}
	return entries, false
}

func applyRename(entries []envfile.Entry, oldKey, newKey string) ([]envfile.Entry, bool) {
	for i, e := range entries {
		if e.Key == oldKey {
			entries[i].Key = newKey
			return entries, true
		}
	}
	return entries, false
}
