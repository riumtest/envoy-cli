// Package sorter provides utilities for sorting and ordering
// environment variable entries by key or value.
package sorter

import (
	"sort"

	"envoy-cli/internal/envfile"
)

// Order defines the sort direction.
type Order string

const (
	// Ascending sorts keys from A to Z.
	Ascending Order = "asc"
	// Descending sorts keys from Z to A.
	Descending Order = "desc"
)

// Options configures how entries are sorted.
type Options struct {
	Order   Order
	ByValue bool
}

// DefaultOptions returns the default sort options (ascending by key).
func DefaultOptions() Options {
	return Options{
		Order:   Ascending,
		ByValue: false,
	}
}

// Sort returns a new slice of entries sorted according options.
// The original slice is not modified.
func Sort(entries []envfile.Entry, opts Options) []envfile.Entry {
	result := make([]envfile.Entry, len(entries))
	copy(result, entries)

	sort.SliceStable(result, func(i, j int) bool {
		var a, b string
		if opts.ByValue {
			a, b = result[i].Value, result[j].Value
		} else {
			a, b = result[i].Key, result[j].Key
		}
		if opts.Order == Descending {
			return a > b
		}
		return a < b
	})

	return result
}

// SortByKey sorts entries by key in ascending order.
func SortByKey(entries []envfile.Entry) []envfile.Entry {
	return Sort(entries, DefaultOptions())
}
