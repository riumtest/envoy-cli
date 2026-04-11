// Package normalizer provides utilities for normalizing .env file entries
// by standardizing key formatting, value quoting, and line endings.
package normalizer

import (
	"strings"

	"envoy-cli/internal/envfile"
)

// Options controls the behaviour of the Normalize function.
type Options struct {
	// UppercaseKeys converts all keys to UPPER_SNAKE_CASE.
	UppercaseKeys bool
	// QuoteValues wraps all non-empty values in double quotes.
	QuoteValues bool
	// TrimValues removes leading and trailing whitespace from values.
	TrimValues bool
	// StripExport removes the "export " prefix from keys if present.
	StripExport bool
}

// DefaultOptions returns a sensible default configuration.
func DefaultOptions() Options {
	return Options{
		UppercaseKeys: true,
		TrimValues:    true,
		StripExport:   true,
		QuoteValues:   false,
	}
}

// Normalize applies the given Options to each entry and returns a new slice
// of normalized entries. The original slice is never mutated.
func Normalize(entries []envfile.Entry, opts Options) []envfile.Entry {
	out := make([]envfile.Entry, len(entries))
	for i, e := range entries {
		out[i] = normalizeEntry(e, opts)
	}
	return out
}

func normalizeEntry(e envfile.Entry, opts Options) envfile.Entry {
	key := e.Key
	val := e.Value

	if opts.StripExport {
		key = strings.TrimPrefix(key, "export ")
		key = strings.TrimPrefix(key, "export\t")
	}

	if opts.UppercaseKeys {
		key = strings.ToUpper(key)
	}

	if opts.TrimValues {
		val = strings.TrimSpace(val)
	}

	if opts.QuoteValues && val != "" && !isQuoted(val) {
		val = "\"" + val + "\""
	}

	return envfile.Entry{Key: key, Value: val}
}

// isQuoted reports whether s is already wrapped in matching single or double
// quotes.
func isQuoted(s string) bool {
	if len(s) < 2 {
		return false
	}
	return (s[0] == '"' && s[len(s)-1] == '"') ||
		(s[0] == '\'' && s[len(s)-1] == '\'')
}
