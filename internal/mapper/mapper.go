// Package mapper provides utilities for transforming env entry slices
// into various map representations for downstream processing.
package mapper

import "github.com/envoy-cli/envoy-cli/internal/envfile"

// Result holds the different map views produced from a slice of entries.
type Result struct {
	// KeyValue maps each key to its raw value.
	KeyValue map[string]string
	// ValueKey is the inverse map — value to key. Duplicate values keep the
	// last key encountered.
	ValueKey map[string]string
	// KeyIndex maps each key to its zero-based position in the original slice.
	KeyIndex map[string]int
}

// ToKeyValue returns a simple key→value map from the given entries.
func ToKeyValue(entries []envfile.Entry) map[string]string {
	out := make(map[string]string, len(entries))
	for _, e := range entries {
		out[e.Key] = e.Value
	}
	return out
}

// ToValueKey returns an inverse value→key map. When multiple entries share
// the same value, the last key wins.
func ToValueKey(entries []envfile.Entry) map[string]string {
	out := make(map[string]string, len(entries))
	for _, e := range entries {
		if e.Value != "" {
			out[e.Value] = e.Key
		}
	}
	return out
}

// ToKeyIndex returns a map of key→index reflecting each entry's position in
// the original slice.
func ToKeyIndex(entries []envfile.Entry) map[string]int {
	out := make(map[string]int, len(entries))
	for i, e := range entries {
		out[e.Key] = i
	}
	return out
}

// Map builds a full Result from the provided entries.
func Map(entries []envfile.Entry) Result {
	return Result{
		KeyValue: ToKeyValue(entries),
		ValueKey: ToValueKey(entries),
		KeyIndex: ToKeyIndex(entries),
	}
}
