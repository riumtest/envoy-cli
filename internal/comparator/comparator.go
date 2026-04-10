// Package comparator provides functionality for comparing two sets of
// environment entries and producing a structured overlap report, including
// shared keys, unique keys, and value mismatches.
package comparator

import "github.com/user/envoy-cli/internal/envfile"

// Result holds the outcome of comparing two env entry sets.
type Result struct {
	SharedKeys    []string          // Keys present in both files with identical values
	MismatchedKeys []MismatchEntry  // Keys present in both files but with different values
	OnlyInLeft    []string          // Keys present only in the left file
	OnlyInRight   []string          // Keys present only in the right file
}

// MismatchEntry represents a key whose value differs between two env files.
type MismatchEntry struct {
	Key        string
	LeftValue  string
	RightValue string
}

// Compare takes two slices of envfile.Entry and produces a Result
// describing their similarities and differences.
func Compare(left, right []envfile.Entry) Result {
	leftMap := toMap(left)
	rightMap := toMap(right)

	result := Result{}

	for key, lv := range leftMap {
		rv, exists := rightMap[key]
		if !exists {
			result.OnlyInLeft = append(result.OnlyInLeft, key)
		} else if lv == rv {
			result.SharedKeys = append(result.SharedKeys, key)
		} else {
			result.MismatchedKeys = append(result.MismatchedKeys, MismatchEntry{
				Key:        key,
				LeftValue:  lv,
				RightValue: rv,
			})
		}
	}

	for key := range rightMap {
		if _, exists := leftMap[key]; !exists {
			result.OnlyInRight = append(result.OnlyInRight, key)
		}
	}

	return result
}

// OverlapRatio returns the fraction of total unique keys that are shared
// (with identical values) between the two files. Returns 0 if no keys exist.
func OverlapRatio(r Result) float64 {
	total := len(r.SharedKeys) + len(r.MismatchedKeys) + len(r.OnlyInLeft) + len(r.OnlyInRight)
	if total == 0 {
		return 0
	}
	return float64(len(r.SharedKeys)) / float64(total)
}

func toMap(entries []envfile.Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}
