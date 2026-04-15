// Package prefixer provides utilities for adding or removing a
// common prefix from environment variable keys.
package prefixer

import (
	"strings"

	"github.com/user/envoy-cli/internal/envfile"
)

// DefaultOptions returns a set of safe default options for the prefixer.
func DefaultOptions() Options {
	return Options{
		Separator: "_",
		Uppercase: true,
	}
}

// Options controls how prefix addition and removal behaves.
type Options struct {
	// Separator is placed between the prefix and the original key.
	Separator string
	// Uppercase converts the final key to upper-case when adding a prefix.
	Uppercase bool
}

// Add prepends prefix to every entry key, returning a new slice.
// The original slice is never mutated.
func Add(entries []envfile.Entry, prefix string, opts Options) []envfile.Entry {
	if prefix == "" {
		return append([]envfile.Entry(nil), entries...)
	}

	p := prefix
	if opts.Uppercase {
		p = strings.ToUpper(p)
	}

	out := make([]envfile.Entry, len(entries))
	for i, e := range entries {
		key := p + opts.Separator + e.Key
		if opts.Uppercase {
			key = strings.ToUpper(key)
		}
		out[i] = envfile.Entry{Key: key, Value: e.Value, Comment: e.Comment}
	}
	return out
}

// Remove strips prefix (and the separator) from every entry key that
// carries it, returning a new slice. Keys that do not carry the prefix
// are left unchanged.
func Remove(entries []envfile.Entry, prefix string, opts Options) []envfile.Entry {
	if prefix == "" {
		return append([]envfile.Entry(nil), entries...)
	}

	p := strings.ToUpper(prefix) + strings.ToUpper(opts.Separator)

	out := make([]envfile.Entry, len(entries))
	for i, e := range entries {
		key := e.Key
		upper := strings.ToUpper(key)
		if strings.HasPrefix(upper, p) {
			key = e.Key[len(p):]
		}
		out[i] = envfile.Entry{Key: key, Value: e.Value, Comment: e.Comment}
	}
	return out
}
