// Package trimmer provides functionality for trimming and normalising
// env file entries by removing excess whitespace, empty lines, and
// optionally deduplicating keys.
package trimmer

import "github.com/user/envoy-cli/internal/envfile"

// Options controls the behaviour of the Trim function.
type Options struct {
	// RemoveEmpty drops entries whose value is an empty string.
	RemoveEmpty bool
	// DeduplicateKeys keeps only the last occurrence of each key.
	DeduplicateKeys bool
	// TrimValues strips leading and trailing whitespace from values.
	TrimValues bool
}

// DefaultOptions returns a sensible default Options configuration.
func DefaultOptions() Options {
	return Options{
		TrimValues:      true,
		RemoveEmpty:     false,
		DeduplicateKeys: false,
	}
}

// Trim applies the given Options to a slice of entries and returns a
// cleaned copy without mutating the original slice.
func Trim(entries []envfile.Entry, opts Options) []envfile.Entry {
	seen := make(map[string]int) // key -> last index in result
	result := make([]envfile.Entry, 0, len(entries))

	for _, e := range entries {
		val := e.Value
		if opts.TrimValues {
			val = trimSpace(val)
		}
		if opts.RemoveEmpty && val == "" {
			continue
		}
		e.Value = val

		if opts.DeduplicateKeys {
			if idx, exists := seen[e.Key]; exists {
				result[idx] = e
				continue
			}
			seen[e.Key] = len(result)
		}
		result = append(result, e)
	}

	return result
}

// trimSpace removes leading and trailing ASCII spaces and tabs.
func trimSpace(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}
