// Package merger provides functionality for merging multiple .env files
// into a single unified set of key-value pairs, with configurable conflict
// resolution strategies.
package merger

import (
	"fmt"

	"github.com/user/envoy-cli/internal/envfile"
)

// Strategy defines how conflicts are resolved when the same key appears
// in multiple files.
type Strategy int

const (
	// StrategyFirst keeps the value from the first file that defines the key.
	StrategyFirst Strategy = iota
	// StrategyLast overwrites with the value from the last file that defines the key.
	StrategyLast
	// StrategyError returns an error when a duplicate key is encountered.
	StrategyError
)

// Result holds the merged entries and metadata about the merge operation.
type Result struct {
	Entries    []envfile.Entry
	Conflicts  []Conflict
	SourceFiles []string
}

// Conflict describes a key that appeared in more than one source file.
type Conflict struct {
	Key      string
	Sources  []string
	Values   []string
	Resolved string
}

// Merge combines entries from multiple parsed env files according to the
// given strategy. Files are processed in the order they are provided.
func Merge(files map[string][]envfile.Entry, order []string, strategy Strategy) (*Result, error) {
	seen := make(map[string]int) // key -> index in result.Entries
	result := &Result{
		SourceFiles: order,
	}

	for _, name := range order {
		entries, ok := files[name]
		if !ok {
			continue
		}
		for _, entry := range entries {
			idx, exists := seen[entry.Key]
			if !exists {
				seen[entry.Key] = len(result.Entries)
				result.Entries = append(result.Entries, entry)
				continue
			}

			// Record the conflict.
			conflict := findOrCreateConflict(result, entry.Key, result.Entries[idx].Value, name, entry.Value)
			_ = conflict

			switch strategy {
			case StrategyFirst:
				// keep existing, do nothing
			case StrategyLast:
				result.Entries[idx].Value = entry.Value
			case StrategyError:
				return nil, fmt.Errorf("merge conflict: key %q defined in multiple files", entry.Key)
			}
		}
	}

	return result, nil
}

func findOrCreateConflict(r *Result, key, existingVal, newSource, newVal string) *Conflict {
	for i, c := range r.Conflicts {
		if c.Key == key {
			r.Conflicts[i].Sources = append(r.Conflicts[i].Sources, newSource)
			r.Conflicts[i].Values = append(r.Conflicts[i].Values, newVal)
			return &r.Conflicts[i]
		}
	}
	r.Conflicts = append(r.Conflicts, Conflict{
		Key:      key,
		Sources:  []string{"<previous>", newSource},
		Values:   []string{existingVal, newVal},
		Resolved: existingVal,
	})
	return &r.Conflicts[len(r.Conflicts)-1]
}
