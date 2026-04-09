// Package filter provides functionality for filtering env file entries
// based on key patterns, prefixes, or value predicates.
package filter

import (
	"strings"

	"github.com/user/envoy-cli/internal/envfile"
)

// Options controls how entries are filtered.
type Options struct {
	// Prefix filters entries whose keys start with the given prefix.
	Prefix string
	// KeyContains filters entries whose keys contain the given substring.
	KeyContains string
	// ExcludeEmpty removes entries with empty values when true.
	ExcludeEmpty bool
	// Keys is an explicit allowlist of keys to keep. If non-empty, only
	// entries whose keys appear in this slice are retained.
	Keys []string
}

// Filter returns a new slice of entries that match all criteria in opts.
// The original slice is never mutated.
func Filter(entries []envfile.Entry, opts Options) []envfile.Entry {
	allowlist := buildAllowlist(opts.Keys)

	var result []envfile.Entry
	for _, e := range entries {
		if len(allowlist) > 0 && !allowlist[e.Key] {
			continue
		}
		if opts.Prefix != "" && !strings.HasPrefix(e.Key, opts.Prefix) {
			continue
		}
		if opts.KeyContains != "" && !strings.Contains(e.Key, opts.KeyContains) {
			continue
		}
		if opts.ExcludeEmpty && strings.TrimSpace(e.Value) == "" {
			continue
		}
		result = append(result, e)
	}
	return result
}

// buildAllowlist converts a slice of keys into a set for O(1) lookup.
func buildAllowlist(keys []string) map[string]bool {
	if len(keys) == 0 {
		return nil
	}
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		m[k] = true
	}
	return m
}
