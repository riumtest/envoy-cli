// Package shrinker reduces a set of env entries by removing keys whose values
// exceed a maximum length or fall below a minimum length threshold.
package shrinker

import "github.com/your-org/envoy-cli/internal/envfile"

// DefaultOptions returns a sensible default configuration.
func DefaultOptions() Options {
	return Options{
		MinLen: 0,
		MaxLen: 0, // 0 means no upper limit
	}
}

// Options controls which entries are kept or removed.
type Options struct {
	// MinLen drops entries whose value length is strictly less than this.
	MinLen int
	// MaxLen drops entries whose value length is strictly greater than this.
	// A value of 0 disables the upper-limit check.
	MaxLen int
	// KeepKeys is an explicit allowlist of keys that are never dropped.
	KeepKeys []string
}

// Shrink returns a new slice containing only entries that satisfy the length
// constraints defined in opts. Entries whose keys appear in KeepKeys are
// always retained regardless of value length.
func Shrink(entries []envfile.Entry, opts Options) []envfile.Entry {
	keep := buildKeepSet(opts.KeepKeys)
	result := make([]envfile.Entry, 0, len(entries))

	for _, e := range entries {
		if keep[e.Key] {
			result = append(result, e)
			continue
		}
		vl := len(e.Value)
		if opts.MinLen > 0 && vl < opts.MinLen {
			continue
		}
		if opts.MaxLen > 0 && vl > opts.MaxLen {
			continue
		}
		result = append(result, e)
	}
	return result
}

func buildKeepSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
