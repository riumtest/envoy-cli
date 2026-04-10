// Package sanitizer provides utilities for cleaning and normalizing
// environment variable entries, such as trimming whitespace, removing
// inline comments, and standardizing key formatting.
package sanitizer

import (
	"strings"

	"github.com/user/envoy-cli/internal/envfile"
)

// Options controls which sanitization steps are applied.
type Options struct {
	// TrimSpace removes leading/trailing whitespace from keys and values.
	TrimSpace bool
	// RemoveInlineComments strips inline comments (e.g. VALUE=foo # comment).
	RemoveInlineComments bool
	// UppercaseKeys converts all keys to UPPER_SNAKE_CASE.
	UppercaseKeys bool
	// RemoveEmpty drops entries with empty values.
	RemoveEmpty bool
}

// DefaultOptions returns a sensible default sanitization configuration.
func DefaultOptions() Options {
	return Options{
		TrimSpace:            true,
		RemoveInlineComments: true,
		UppercaseKeys:        false,
		RemoveEmpty:          false,
	}
}

// Sanitize applies the given Options to a slice of envfile.Entry values,
// returning a new slice of cleaned entries.
func Sanitize(entries []envfile.Entry, opts Options) []envfile.Entry {
	result := make([]envfile.Entry, 0, len(entries))

	for _, e := range entries {
		key := e.Key
		value := e.Value

		if opts.TrimSpace {
			key = strings.TrimSpace(key)
			value = strings.TrimSpace(value)
		}

		if opts.RemoveInlineComments {
			value = stripInlineComment(value)
			if opts.TrimSpace {
				value = strings.TrimSpace(value)
			}
		}

		if opts.UppercaseKeys {
			key = strings.ToUpper(key)
		}

		if opts.RemoveEmpty && value == "" {
			continue
		}

		result = append(result, envfile.Entry{Key: key, Value: value})
	}

	return result
}

// stripInlineComment removes an unquoted inline comment from a value string.
// Values enclosed in single or double quotes are left untouched.
func stripInlineComment(value string) string {
	if len(value) == 0 {
		return value
	}

	// Preserve quoted values entirely.
	if (strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`)) ||
		(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
		return value
	}

	// Strip everything from the first unquoted '#' onward.
	if idx := strings.Index(value, " #"); idx != -1 {
		return value[:idx]
	}

	return value
}
