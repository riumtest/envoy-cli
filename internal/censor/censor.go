// Package censor provides entry-level censoring of sensitive environment
// variable values, replacing them with a configurable placeholder.
package censor

import (
	"strings"

	"github.com/envoy-cli/internal/envfile"
	"github.com/envoy-cli/internal/masker"
)

// DefaultPlaceholder is used when no custom placeholder is specified.
const DefaultPlaceholder = "***"

// Options configures how censoring is applied.
type Options struct {
	// Placeholder replaces sensitive values. Defaults to DefaultPlaceholder.
	Placeholder string
	// ExtraPatterns are additional key-name patterns treated as sensitive.
	ExtraPatterns []string
	// Keys is an explicit list of keys to censor regardless of pattern.
	Keys []string
}

// Censor replaces sensitive values in entries with a placeholder.
// Sensitivity is determined by the masker package plus any explicit keys
// provided in opts.
func Censor(entries []envfile.Entry, opts Options) []envfile.Entry {
	if opts.Placeholder == "" {
		opts.Placeholder = DefaultPlaceholder
	}

	m := masker.NewWithPatterns(opts.ExtraPatterns)
	explicit := buildExplicitSet(opts.Keys)

	result := make([]envfile.Entry, len(entries))
	for i, e := range entries {
		if explicit[strings.ToUpper(e.Key)] || m.IsSensitive(e.Key) {
			e.Value = opts.Placeholder
		}
		result[i] = e
	}
	return result
}

// buildExplicitSet converts a slice of key names into an uppercase lookup map.
func buildExplicitSet(keys []string) map[string]bool {
	set := make(map[string]bool, len(keys))
	for _, k := range keys {
		set[strings.ToUpper(k)] = true
	}
	return set
}
