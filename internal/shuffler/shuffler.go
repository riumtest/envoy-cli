// Package shuffler provides functionality for randomising the order of
// environment entries. This is useful for testing that downstream tools
// do not rely on a specific key ordering.
package shuffler

import (
	"math/rand"
	"time"

	"envoy-cli/internal/envfile"
)

// DefaultOptions returns an Options value with sensible defaults.
func DefaultOptions() Options {
	return Options{
		Seed: time.Now().UnixNano(),
	}
}

// Options controls the behaviour of Shuffle.
type Options struct {
	// Seed is used to initialise the random source. Set a fixed value to
	// produce a deterministic shuffle (useful in tests).
	Seed int64
}

// Shuffle returns a new slice containing the same entries as src but in a
// randomised order. The original slice is never mutated.
func Shuffle(entries []envfile.Entry, opts Options) []envfile.Entry {
	if len(entries) == 0 {
		return []envfile.Entry{}
	}

	// Copy so we never mutate the caller's slice.
	out := make([]envfile.Entry, len(entries))
	copy(out, entries)

	//nolint:gosec // non-cryptographic shuffle is intentional here
	r := rand.New(rand.NewSource(opts.Seed))
	r.Shuffle(len(out), func(i, j int) {
		out[i], out[j] = out[j], out[i]
	})

	return out
}
