// Package splitter provides functionality for splitting a set of env entries
// into multiple groups based on key prefixes or custom rules.
package splitter

import (
	"strings"

	"github.com/user/envoy-cli/internal/envfile"
)

// Options controls how splitting is performed.
type Options struct {
	// Prefixes is the list of prefixes to split by.
	// Each prefix becomes its own output group.
	Prefixes []string
	// StripPrefix removes the matched prefix from keys in the output group.
	StripPrefix bool
	// IncludeUnmatched collects entries that match no prefix into a group keyed by "".
	IncludeUnmatched bool
}

// DefaultOptions returns a sensible default Options value.
func DefaultOptions() Options {
	return Options{
		StripPrefix:      false,
		IncludeUnmatched: true,
	}
}

// Split partitions entries into groups keyed by the matched prefix.
// Entries that match no prefix are placed in the "" group when
// Options.IncludeUnmatched is true.
func Split(entries []envfile.Entry, opts Options) map[string][]envfile.Entry {
	result := make(map[string][]envfile.Entry)

	for _, e := range entries {
		matched := false
		for _, prefix := range opts.Prefixes {
			norm := normalise(prefix)
			if strings.HasPrefix(strings.ToUpper(e.Key), norm) {
				key := e.Key
				if opts.StripPrefix {
					key = e.Key[len(prefix):]
					key = strings.TrimPrefix(key, "_")
				}
				result[prefix] = append(result[prefix], envfile.Entry{Key: key, Value: e.Value})
				matched = true
				break
			}
		}
		if !matched && opts.IncludeUnmatched {
			result[""] = append(result[""], e)
		}
	}
	return result
}

func normalise(prefix string) string {
	return strings.ToUpper(strings.TrimSuffix(prefix, "_"))
}
