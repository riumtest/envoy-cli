// Package pitcher provides functionality for promoting env entries
// from one environment to another with optional key filtering and
// value transformation.
package pitcher

import (
	"fmt"
	"strings"

	"github.com/user/envoy-cli/internal/envfile"
)

// Options configures the Pitch operation.
type Options struct {
	// Keys restricts promotion to these specific keys. Empty means all.
	Keys []string
	// Overwrite controls whether existing keys in the target are replaced.
	Overwrite bool
	// Prefix is prepended to each promoted key.
	Prefix string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Overwrite: true,
	}
}

// Result holds the outcome of a Pitch operation.
type Result struct {
	Promoted []string
	Skipped  []string
	Entries  []envfile.Entry
}

// Pitch promotes entries from src into dst according to opts.
func Pitch(src, dst []envfile.Entry, opts Options) (Result, error) {
	allow := buildAllowlist(opts.Keys)
	targetMap := toMap(dst)

	var result Result
	result.Entries = make([]envfile.Entry, len(dst))
	copy(result.Entries, dst)

	for _, e := range src {
		if len(allow) > 0 && !allow[e.Key] {
			continue
		}
		destKey := e.Key
		if opts.Prefix != "" {
			destKey = strings.ToUpper(opts.Prefix) + "_" + e.Key
		}
		if _, exists := targetMap[destKey]; exists && !opts.Overwrite {
			result.Skipped = append(result.Skipped, destKey)
			continue
		}
		if _, exists := targetMap[destKey]; exists {
			for i, t := range result.Entries {
				if t.Key == destKey {
					result.Entries[i].Value = e.Value
					break
				}
			}
		} else {
			result.Entries = append(result.Entries, envfile.Entry{Key: destKey, Value: e.Value})
		}
		targetMap[destKey] = e.Value
		result.Promoted = append(result.Promoted, destKey)
	}
	if len(allow) > 0 {
		for k := range allow {
			if _, found := toMap(src)[k]; !found {
				return result, fmt.Errorf("key %q not found in source", k)
			}
		}
	}
	return result, nil
}

func buildAllowlist(keys []string) map[string]bool {
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		m[k] = true
	}
	return m
}

func toMap(entries []envfile.Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}
