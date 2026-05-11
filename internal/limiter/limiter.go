// Package limiter provides functionality for limiting the number of entries
// returned from a set, with optional offset support for pagination-style slicing.
package limiter

import "github.com/your-org/envoy-cli/internal/envfile"

// DefaultOptions returns a Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		Limit:  0,
		Offset: 0,
	}
}

// Options controls how Limit behaves.
type Options struct {
	// Limit is the maximum number of entries to return. 0 means no limit.
	Limit int
	// Offset is the number of entries to skip before collecting results.
	Offset int
}

// Limit returns a slice of entries respecting the given offset and limit.
// If Offset exceeds the length of entries, an empty slice is returned.
// If Limit is 0, all entries after the offset are returned.
func Limit(entries []envfile.Entry, opts Options) []envfile.Entry {
	if len(entries) == 0 {
		return []envfile.Entry{}
	}

	offset := opts.Offset
	if offset < 0 {
		offset = 0
	}
	if offset >= len(entries) {
		return []envfile.Entry{}
	}

	sliced := entries[offset:]

	if opts.Limit <= 0 || opts.Limit >= len(sliced) {
		result := make([]envfile.Entry, len(sliced))
		copy(result, sliced)
		return result
	}

	result := make([]envfile.Entry, opts.Limit)
	copy(result, sliced[:opts.Limit])
	return result
}
