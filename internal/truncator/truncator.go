// Package truncator provides utilities for truncating .env entry values
// to a maximum length, with optional suffix indicators.
package truncator

import "github.com/envoy-cli/envoy-cli/internal/envfile"

// DefaultOptions returns a sensible default Options value.
func DefaultOptions() Options {
	return Options{
		MaxLength: 64,
		Suffix:    "...",
		KeepKeys:  nil,
	}
}

// Options controls truncation behaviour.
type Options struct {
	// MaxLength is the maximum number of runes allowed in a value.
	MaxLength int

	// Suffix is appended when a value is truncated. Counts toward MaxLength.
	Suffix string

	// KeepKeys lists keys whose values must never be truncated.
	KeepKeys []string
}

// Truncate returns a new slice of entries with values truncated according to
// opts. The original slice is never mutated.
func Truncate(entries []envfile.Entry, opts Options) []envfile.Entry {
	if opts.MaxLength <= 0 {
		opts.MaxLength = DefaultOptions().MaxLength
	}

	keepSet := buildKeepSet(opts.KeepKeys)

	out := make([]envfile.Entry, len(entries))
	for i, e := range entries {
		if keepSet[e.Key] {
			out[i] = e
			continue
		}
		out[i] = envfile.Entry{
			Key:   e.Key,
			Value: truncateValue(e.Value, opts.MaxLength, opts.Suffix),
		}
	}
	return out
}

func truncateValue(v string, max int, suffix string) string {
	runes := []rune(v)
	if len(runes) <= max {
		return v
	}
	suffixRunes := []rune(suffix)
	cutAt := max - len(suffixRunes)
	if cutAt < 0 {
		cutAt = 0
	}
	return string(runes[:cutAt]) + suffix
}

func buildKeepSet(keys []string) map[string]bool {
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		m[k] = true
	}
	return m
}
