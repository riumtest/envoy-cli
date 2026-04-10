// Package deduplicator provides functionality for removing duplicate entries
// from a list of environment variables, with configurable strategies for
// which duplicate to retain.
package deduplicator

import "github.com/yourusername/envoy-cli/internal/envfile"

// Strategy determines which entry to keep when duplicates are found.
type Strategy int

const (
	// KeepFirst retains the first occurrence of a duplicate key.
	KeepFirst Strategy = iota
	// KeepLast retains the last occurrence of a duplicate key.
	KeepLast
)

// Result holds the output of a deduplication operation.
type Result struct {
	Entries    []envfile.Entry
	Duplicates []Duplicate
}

// Duplicate records a key that appeared more than once and how many times.
type Duplicate struct {
	Key   string
	Count int
}

// Deduplicate removes duplicate keys from entries according to the given strategy.
// It returns a Result containing the cleaned entries and a report of duplicates found.
func Deduplicate(entries []envfile.Entry, strategy Strategy) Result {
	type occurrence struct {
		index int
		count int
	}

	seen := make(map[string]*occurrence)
	order := make([]string, 0, len(entries))

	for i, e := range entries {
		if occ, exists := seen[e.Key]; exists {
			occ.count++
			if strategy == KeepLast {
				occ.index = i
			}
		} else {
			seen[e.Key] = &occurrence{index: i, count: 1}
			order = append(order, e.Key)
		}
	}

	result := Result{
		Entries:    make([]envfile.Entry, 0, len(seen)),
		Duplicates: []Duplicate{},
	}

	for _, key := range order {
		occ := seen[key]
		result.Entries = append(result.Entries, entries[occ.index])
		if occ.count > 1 {
			result.Duplicates = append(result.Duplicates, Duplicate{
				Key:   key,
				Count: occ.count,
			})
		}
	}

	return result
}
