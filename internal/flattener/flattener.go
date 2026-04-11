package flattener

import (
	"fmt"
	"strings"

	"github.com/envoy-cli/internal/envfile"
)

// Options controls how flattening is performed.
type Options struct {
	// Separator is placed between the prefix and the key (default: "_").
	Separator string
	// Uppercase forces all resulting keys to uppercase.
	Uppercase bool
}

// DefaultOptions returns sensible defaults for flattening.
func DefaultOptions() Options {
	return Options{
		Separator: "_",
		Uppercase: true,
	}
}

// Flatten takes a slice of entries and a namespace prefix, returning a new
// slice where every key is prefixed with the namespace using the separator.
// Duplicate keys after prefixing are deduplicated; the first occurrence wins.
func Flatten(entries []envfile.Entry, namespace string, opts Options) []envfile.Entry {
	if opts.Separator == "" {
		opts.Separator = "_"
	}

	seen := make(map[string]struct{}, len(entries))
	result := make([]envfile.Entry, 0, len(entries))

	for _, e := range entries {
		key := buildKey(namespace, e.Key, opts)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, envfile.Entry{Key: key, Value: e.Value})
	}
	return result
}

// Unflatten strips a namespace prefix from all keys that carry it.
// Keys that do not start with the prefix are left unchanged.
func Unflatten(entries []envfile.Entry, namespace string, opts Options) []envfile.Entry {
	if opts.Separator == "" {
		opts.Separator = "_"
	}

	prefix := buildPrefix(namespace, opts)
	result := make([]envfile.Entry, 0, len(entries))
	for _, e := range entries {
		key := e.Key
		if strings.HasPrefix(key, prefix) {
			key = key[len(prefix):]
		}
		result = append(result, envfile.Entry{Key: key, Value: e.Value})
	}
	return result
}

func buildKey(namespace, key string, opts Options) string {
	var full string
	if namespace == "" {
		full = key
	} else {
		full = fmt.Sprintf("%s%s%s", namespace, opts.Separator, key)
	}
	if opts.Uppercase {
		return strings.ToUpper(full)
	}
	return full
}

func buildPrefix(namespace string, opts Options) string {
	if namespace == "" {
		return ""
	}
	prefix := namespace + opts.Separator
	if opts.Uppercase {
		return strings.ToUpper(prefix)
	}
	return prefix
}
