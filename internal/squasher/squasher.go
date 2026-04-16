// Package squasher merges duplicate keys in a slice of env entries,
// collapsing them into a single entry according to a chosen strategy.
package squasher

import (
	"errors"

	"github.com/user/envoy-cli/internal/envfile"
)

// Strategy controls which value wins when duplicate keys are squashed.
type Strategy string

const (
	StrategyFirst Strategy = "first" // keep the first occurrence
	StrategyLast  Strategy = "last"  // keep the last occurrence
	StrategyError Strategy = "error" // return an error on any duplicate
)

// Options configures Squash behaviour.
type Options struct {
	Strategy Strategy
}

// DefaultOptions returns sensible defaults (keep last).
func DefaultOptions() Options {
	return Options{Strategy: StrategyLast}
}

// Result holds the squashed entries and metadata.
type Result struct {
	Entries    []envfile.Entry
	Squashed   int // number of duplicate entries removed
}

// Squash collapses duplicate keys in entries according to opts.
func Squash(entries []envfile.Entry, opts Options) (Result, error) {
	if opts.Strategy == "" {
		opts = DefaultOptions()
	}

	seen := make(map[string]int) // key -> index in out
	out := make([]envfile.Entry, 0, len(entries))
	squashed := 0

	for _, e := range entries {
		if idx, exists := seen[e.Key]; exists {
			switch opts.Strategy {
			case StrategyError:
				return Result{}, errors.New("duplicate key: " + e.Key)
			case StrategyFirst:
				// discard the new entry
				squashed++
			case StrategyLast:
				out[idx] = e
				squashed++
			}
		} else {
			seen[e.Key] = len(out)
			out = append(out, e)
		}
	}

	return Result{Entries: out, Squashed: squashed}, nil
}
