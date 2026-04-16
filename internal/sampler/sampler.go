// Package sampler provides functionality for sampling a subset of env entries.
package sampler

import (
	"math/rand"

	"github.com/envoy-cli/internal/envfile"
)

// DefaultOptions returns a Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		N:    5,
		Seed: 42,
	}
}

// Options controls sampling behaviour.
type Options struct {
	// N is the number of entries to sample. If N >= len(entries), all entries are returned.
	N int
	// Seed is the random seed for reproducible sampling.
	Seed int64
	// SensitiveOnly restricts sampling to sensitive-looking keys.
	SensitiveOnly bool
}

// Sample returns up to opts.N randomly selected entries from the input slice.
// The original slice is never mutated.
func Sample(entries []envfile.Entry, opts Options) []envfile.Entry {
	pool := make([]envfile.Entry, len(entries))
	copy(pool, entries)

	if opts.SensitiveOnly {
		pool = filterSensitive(pool)
	}

	if opts.N <= 0 || opts.N >= len(pool) {
		return pool
	}

	r := rand.New(rand.NewSource(opts.Seed)) //nolint:gosec
	r.Shuffle(len(pool), func(i, j int) {
		pool[i], pool[j] = pool[j], pool[i]
	})

	return pool[:opts.N]
}

var sensitivePatterns = []string{"SECRET", "PASSWORD", "TOKEN", "KEY", "PASS", "PRIVATE", "AUTH"}

func filterSensitive(entries []envfile.Entry) []envfile.Entry {
	var out []envfile.Entry
	for _, e := range entries {
		if isSensitive(e.Key) {
			out = append(out, e)
		}
	}
	return out
}

func isSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range sensitivePatterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}
