// Package cloner provides functionality for cloning and overriding .env entries
// from a source environment into a target, with optional key filtering and masking.
package cloner

import (
	"envoy-cli/internal/envfile"
)

// Options configures the cloning behaviour.
type Options struct {
	// Keys restricts cloning to only these keys. If empty, all keys are cloned.
	Keys []string
	// Overwrite controls whether existing keys in the target are overwritten.
	Overwrite bool
}

// DefaultOptions returns sensible defaults for cloning.
func DefaultOptions() Options {
	return Options{
		Keys:      nil,
		Overwrite: true,
	}
}

// Result holds the outcome of a clone operation.
type Result struct {
	Entries  []envfile.Entry
	Cloned   int
	Skipped  int
	Conflict int
}

// Clone merges entries from src into dst according to opts.
// It returns a Result describing what happened.
func Clone(dst, src []envfile.Entry, opts Options) (Result, error) {
	allowlist := buildAllowlist(opts.Keys)

	dstMap := make(map[string]int, len(dst))
	for i, e := range dst {
		dstMap[e.Key] = i
	}

	out := make([]envfile.Entry, len(dst))
	copy(out, dst)

	var cloned, skipped, conflict int

	for _, e := range src {
		if len(allowlist) > 0 && !allowlist[e.Key] {
			skipped++
			continue
		}

		if idx, exists := dstMap[e.Key]; exists {
			if !opts.Overwrite {
				conflict++
				continue
			}
			out[idx] = e
			cloned++
		} else {
			out = append(out, e)
			dstMap[e.Key] = len(out) - 1
			cloned++
		}
	}

	return Result{
		Entries:  out,
		Cloned:   cloned,
		Skipped:  skipped,
		Conflict: conflict,
	}, nil
}

// Summary returns a human-readable description of the clone result.
func (r Result) Summary() string {
	return fmt.Sprintf("cloned %d, skipped %d, conflicts %d", r.Cloned, r.Skipped, r.Conflict)
}

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
