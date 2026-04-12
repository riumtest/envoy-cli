// Package freezer provides functionality to freeze a set of env entries
// into a read-only snapshot and detect any mutations against it.
package freezer

import (
	"errors"
	"fmt"

	"github.com/user/envoy-cli/internal/envfile"
)

// FrozenEntry holds the original key/value at freeze time.
type FrozenEntry struct {
	Key   string
	Value string
}

// Violation describes a key that has changed relative to the frozen state.
type Violation struct {
	Key      string
	Frozen   string
	Current  string
	Kind     string // "mutated", "deleted", "added"
}

// Freeze captures the current state of entries as a frozen baseline.
func Freeze(entries []envfile.Entry) []FrozenEntry {
	out := make([]FrozenEntry, len(entries))
	for i, e := range entries {
		out[i] = FrozenEntry{Key: e.Key, Value: e.Value}
	}
	return out
}

// Check compares current entries against the frozen baseline and returns
// any violations. An error is returned if violations are found and
// failOnViolation is true.
func Check(frozen []FrozenEntry, current []envfile.Entry, failOnViolation bool) ([]Violation, error) {
	frozenMap := make(map[string]string, len(frozen))
	for _, f := range frozen {
		frozenMap[f.Key] = f.Value
	}

	currentMap := make(map[string]string, len(current))
	for _, e := range current {
		currentMap[e.Key] = e.Value
	}

	var violations []Violation

	for _, f := range frozen {
		curVal, exists := currentMap[f.Key]
		if !exists {
			violations = append(violations, Violation{Key: f.Key, Frozen: f.Value, Kind: "deleted"})
			continue
		}
		if curVal != f.Value {
			violations = append(violations, Violation{Key: f.Key, Frozen: f.Value, Current: curVal, Kind: "mutated"})
		}
	}

	for _, e := range current {
		if _, exists := frozenMap[e.Key]; !exists {
			violations = append(violations, Violation{Key: e.Key, Current: e.Value, Kind: "added"})
		}
	}

	if failOnViolation && len(violations) > 0 {
		return violations, errors.New(fmt.Sprintf("%d freeze violation(s) detected", len(violations)))
	}
	return violations, nil
}
