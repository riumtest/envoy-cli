package grouper

import (
	"sort"
	"strings"

	"github.com/envoy-cli/internal/envfile"
)

// Group holds a named collection of env entries sharing a common prefix.
type Group struct {
	Name    string
	Entries []envfile.Entry
}

// Options controls grouping behaviour.
type Options struct {
	// Delimiter separates prefix from the rest of the key (default "_").
	Delimiter string
	// MinSize skips groups with fewer entries than this value (0 = include all).
	MinSize int
	// SortGroups sorts group names alphabetically when true.
	SortGroups bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Delimiter:  "_",
		MinSize:    1,
		SortGroups: true,
	}
}

// GroupByPrefix groups entries by the first segment of their key split on
// opts.Delimiter. Keys with no delimiter are placed in a group named "_other".
func GroupByPrefix(entries []envfile.Entry, opts Options) []Group {
	if opts.Delimiter == "" {
		opts.Delimiter = "_"
	}

	index := make(map[string][]envfile.Entry)

	for _, e := range entries {
		parts := strings.SplitN(e.Key, opts.Delimiter, 2)
		prefix := parts[0]
		if len(parts) == 1 {
			prefix = "_other"
		}
		index[prefix] = append(index[prefix], e)
	}

	var groups []Group
	for name, es := range index {
		if len(es) < opts.MinSize {
			continue
		}
		groups = append(groups, Group{Name: name, Entries: es})
	}

	if opts.SortGroups {
		sort.Slice(groups, func(i, j int) bool {
			return groups[i].Name < groups[j].Name
		})
	}

	return groups
}

// ToMap converts a slice of Groups into a map keyed by group name.
func ToMap(groups []Group) map[string][]envfile.Entry {
	m := make(map[string][]envfile.Entry, len(groups))
	for _, g := range groups {
		m[g.Name] = g.Entries
	}
	return m
}
