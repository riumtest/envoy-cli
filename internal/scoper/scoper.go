// Package scoper provides functionality for scoping .env entries
// to a named environment (e.g., "production", "staging") by prefixing
// or stripping environment-specific key prefixes.
package scoper

import (
	"fmt"
	"strings"

	"github.com/yourusername/envoy-cli/internal/envfile"
)

// Options controls how scoping is applied.
type Options struct {
	// Separator between the scope prefix and the key name (default: "_").
	Separator string
	// UpperCase converts the scope prefix to uppercase before applying.
	UpperCase bool
}

// DefaultOptions returns sensible defaults for scoping.
func DefaultOptions() Options {
	return Options{
		Separator: "_",
		UpperCase: true,
	}
}

// Scope adds a scope prefix to each entry's key.
// For example, with scope "prod" and separator "_", KEY becomes PROD_KEY.
func Scope(entries []envfile.Entry, scope string, opts Options) []envfile.Entry {
	if scope == "" {
		return entries
	}
	prefix := buildPrefix(scope, opts)
	result := make([]envfile.Entry, len(entries))
	for i, e := range entries {
		result[i] = envfile.Entry{
			Key:   prefix + e.Key,
			Value: e.Value,
		}
	}
	return result
}

// Unscope strips a scope prefix from each entry's key.
// Entries whose keys do not begin with the prefix are left unchanged.
func Unscope(entries []envfile.Entry, scope string, opts Options) []envfile.Entry {
	if scope == "" {
		return entries
	}
	prefix := buildPrefix(scope, opts)
	result := make([]envfile.Entry, len(entries))
	for i, e := range entries {
		key := e.Key
		if strings.HasPrefix(key, prefix) {
			key = strings.TrimPrefix(key, prefix)
		}
		result[i] = envfile.Entry{
			Key:   key,
			Value: e.Value,
		}
	}
	return result
}

// FilterByScope returns only entries whose keys start with the given scope prefix.
func FilterByScope(entries []envfile.Entry, scope string, opts Options) []envfile.Entry {
	if scope == "" {
		return entries
	}
	prefix := buildPrefix(scope, opts)
	var result []envfile.Entry
	for _, e := range entries {
		if strings.HasPrefix(e.Key, prefix) {
			result = append(result, e)
		}
	}
	return result
}

func buildPrefix(scope string, opts Options) string {
	sep := opts.Separator
	if sep == "" {
		sep = "_"
	}
	if opts.UpperCase {
		scope = strings.ToUpper(scope)
	}
	return fmt.Sprintf("%s%s", scope, sep)
}
