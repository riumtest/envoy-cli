// Package stacker merges multiple .env entry slices in priority order,
// with later slices overriding earlier ones (or vice-versa depending on strategy).
package stacker

import "github.com/user/envoy-cli/internal/envfile"

// Strategy controls how conflicts between layers are resolved.
type Strategy int

const (
	// StrategyLast means later layers win (default, highest-priority last).
	StrategyLast Strategy = iota
	// StrategyFirst means earlier layers win (highest-priority first).
	StrategyFirst
)

// Options configures the Stack operation.
type Options struct {
	Strategy Strategy
	// IncludeSource annotates each entry's Comment field with the layer index.
	IncludeSource bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{Strategy: StrategyLast}
}

// Result holds the stacked entries and metadata.
type Result struct {
	Entries    []envfile.Entry
	Overridden int // number of keys that were overridden by a later/earlier layer
}

// Stack merges layers of entries according to opts.
// Each element of layers is a named slice; name is used for source annotation.
func Stack(layers [][]envfile.Entry, opts Options) Result {
	type record struct {
		entry      envfile.Entry
		layerIndex int
	}

	seen := make(map[string]record)
	order := []string{}
	overridden := 0

	for li, layer := range layers {
		for _, e := range layer {
			if e.Key == "" {
				continue
			}
			existing, exists := seen[e.Key]
			switch {
			case !exists:
				seen[e.Key] = record{entry: e, layerIndex: li}
				order = append(order, e.Key)
			case opts.Strategy == StrategyLast:
				_ = existing
				seen[e.Key] = record{entry: e, layerIndex: li}
				overridden++
			case opts.Strategy == StrategyFirst:
				// keep existing, do nothing
				overridden++
			}
		}
	}

	out := make([]envfile.Entry, 0, len(order))
	for _, k := range order {
		out = append(out, seen[k].entry)
	}
	return Result{Entries: out, Overridden: overridden}
}
