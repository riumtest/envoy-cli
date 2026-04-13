// Package chainer provides support for applying a sequence of named
// transformation steps to a slice of env entries in a defined order.
package chainer

import (
	"fmt"

	"github.com/envoy-cli/envoy/internal/envfile"
)

// StepFn is a transformation function that receives entries and returns
// transformed entries or an error.
type StepFn func(entries []envfile.Entry) ([]envfile.Entry, error)

// Step pairs a human-readable name with a transformation function.
type Step struct {
	Name string
	Fn   StepFn
}

// Result holds the outcome of a single pipeline step.
type Result struct {
	Step    string
	Entries []envfile.Entry
	Err     error
}

// Chain executes a sequence of Steps in order, passing the output of each
// step as the input to the next. Execution stops on the first error.
// All intermediate Results are returned regardless of early termination.
func Chain(entries []envfile.Entry, steps []Step) ([]envfile.Entry, []Result, error) {
	results := make([]Result, 0, len(steps))
	current := entries

	for _, s := range steps {
		if s.Fn == nil {
			return current, results, fmt.Errorf("step %q has a nil function", s.Name)
		}
		out, err := s.Fn(current)
		results = append(results, Result{Step: s.Name, Entries: out, Err: err})
		if err != nil {
			return current, results, fmt.Errorf("step %q failed: %w", s.Name, err)
		}
		current = out
	}

	return current, results, nil
}
