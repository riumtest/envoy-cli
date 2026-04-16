// Package swapper provides functionality to swap keys and values in env entries.
package swapper

import "github.com/envoy-cli/envoy-cli/internal/envfile"

// DefaultOptions returns a default Options value.
func DefaultOptions() Options {
	return Options{SkipEmpty: true}
}

// Options controls Swap behaviour.
type Options struct {
	// SkipEmpty skips entries whose value is empty.
	SkipEmpty bool
	// SkipDuplicates drops swapped entries whose new key already exists.
	SkipDuplicates bool
}

// Swap returns a new slice where each entry's Key and Value are exchanged.
// Entries whose value would produce an empty key are skipped when SkipEmpty
// is set. When SkipDuplicates is set, only the first occurrence of a new key
// is kept.
func Swap(entries []envfile.Entry, opts Options) []envfile.Entry {
	seen := make(map[string]bool)
	result := make([]envfile.Entry, 0, len(entries))

	for _, e := range entries {
		if opts.SkipEmpty && e.Value == "" {
			continue
		}
		newKey := e.Value
		newVal := e.Key
		if opts.SkipDuplicates {
			if seen[newKey] {
				continue
			}
			seen[newKey] = true
		}
		result = append(result, envfile.Entry{Key: newKey, Value: newVal})
	}
	return result
}
