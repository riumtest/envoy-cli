// Package aliaser provides functionality for creating key aliases in .env files.
// An alias maps one or more new key names to the value of an existing key,
// allowing the same value to be referenced under multiple names.
package aliaser

import (
	"fmt"

	"github.com/envoy-cli/internal/envfile"
)

// Rule defines a single alias mapping: From is the source key whose value
// is copied, To is the new key name to create.
type Rule struct {
	From string
	To   string
}

// Result holds the outcome of an alias operation.
type Result struct {
	Entries  []envfile.Entry
	Aliased  []Rule
	Skipped  []Rule
	Conflicts []Rule
}

// Options controls the behaviour of the Alias function.
type Options struct {
	// Overwrite allows existing keys to be overwritten by an alias.
	Overwrite bool
	// SkipMissing silently ignores rules whose source key does not exist.
	SkipMissing bool
}

// DefaultOptions returns sensible defaults for aliasing.
func DefaultOptions() Options {
	return Options{
		Overwrite:   false,
		SkipMissing: false,
	}
}

// Alias applies the given rules to entries, producing new entries that copy
// the value from the source key into the destination key.
func Alias(entries []envfile.Entry, rules []Rule, opts Options) (Result, error) {
	index := make(map[string]string, len(entries))
	for _, e := range entries {
		index[e.Key] = e.Value
	}

	existing := make(map[string]struct{}, len(entries))
	for _, e := range entries {
		existing[e.Key] = struct{}{}
	}

	out := make([]envfile.Entry, len(entries))
	copy(out, entries)

	result := Result{}

	for _, rule := range rules {
		val, ok := index[rule.From]
		if !ok {
			if opts.SkipMissing {
				result.Skipped = append(result.Skipped, rule)
				continue
			}
			return Result{}, fmt.Errorf("aliaser: source key %q not found", rule.From)
		}

		if _, exists := existing[rule.To]; exists {
			if !opts.Overwrite {
				result.Conflicts = append(result.Conflicts, rule)
				continue
			}
			// overwrite: update value in-place
			for i, e := range out {
				if e.Key == rule.To {
					out[i].Value = val
					break
				}
			}
		} else {
			out = append(out, envfile.Entry{Key: rule.To, Value: val})
			existing[rule.To] = struct{}{}
		}
		result.Aliased = append(result.Aliased, rule)
	}

	result.Entries = out
	return result, nil
}
