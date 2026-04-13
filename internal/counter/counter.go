// Package counter provides utilities for counting and aggregating
// statistics across a set of env file entries.
package counter

import (
	"sort"

	"envoy-cli/internal/envfile"
)

// Result holds the aggregated counts for a set of entries.
type Result struct {
	Total       int
	Empty       int
	NonEmpty    int
	Unique      int
	Duplicated  int
	ByPrefix    map[string]int
}

// Count analyses the provided entries and returns a Result with
// aggregated statistics. Prefix grouping uses the first segment
// before the first underscore in each key.
func Count(entries []envfile.Entry) Result {
	seen := make(map[string]int)
	prefixes := make(map[string]int)

	for _, e := range entries {
		seen[e.Key]++

		prefix := extractPrefix(e.Key)
		prefixes[prefix]++
	}

	var empty, nonEmpty, unique, duplicated int
	for _, e := range entries {
		if e.Value == "" {
			empty++
		} else {
			nonEmpty++
		}
	}

	for _, count := range seen {
		if count == 1 {
			unique++
		} else {
			duplicated++
		}
	}

	return Result{
		Total:      len(entries),
		Empty:      empty,
		NonEmpty:   nonEmpty,
		Unique:     unique,
		Duplicated: duplicated,
		ByPrefix:   prefixes,
	}
}

// TopPrefixes returns up to n prefixes sorted by descending count.
func TopPrefixes(r Result, n int) []string {
	type kv struct {
		Key   string
		Count int
	}

	pairs := make([]kv, 0, len(r.ByPrefix))
	for k, v := range r.ByPrefix {
		pairs = append(pairs, kv{k, v})
	}

	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].Count != pairs[j].Count {
			return pairs[i].Count > pairs[j].Count
		}
		return pairs[i].Key < pairs[j].Key
	})

	result := make([]string, 0, n)
	for i, p := range pairs {
		if i >= n {
			break
		}
		result = append(result, p.Key)
	}
	return result
}

func extractPrefix(key string) string {
	for i, ch := range key {
		if ch == '_' && i > 0 {
			return key[:i]
		}
	}
	return key
}
