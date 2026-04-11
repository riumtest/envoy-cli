// Package pinner provides functionality for pinning (locking) environment
// variable values to a snapshot, producing a diff of what has drifted.
package pinner

import (
	"fmt"

	"github.com/envoy-cli/envoy/internal/envfile"
)

// DriftEntry describes a single key whose live value differs from its pinned value.
type DriftEntry struct {
	Key     string
	Pinned  string
	Current string
}

// Result holds the outcome of a pin comparison.
type Result struct {
	Drifted []DriftEntry
	Missing []string // keys present in pin but absent from current
	New     []string // keys present in current but absent from pin
}

// Summary returns a human-readable one-liner.
func (r Result) Summary() string {
	return fmt.Sprintf("%d drifted, %d missing, %d new",
		len(r.Drifted), len(r.Missing), len(r.New))
}

// Pin compares current entries against a pinned baseline and returns drift.
func Pin(pinned, current []envfile.Entry) Result {
	pinnedMap := toMap(pinned)
	currentMap := toMap(current)

	var result Result

	for _, p := range pinned {
		cv, ok := currentMap[p.Key]
		if !ok {
			result.Missing = append(result.Missing, p.Key)
			continue
		}
		if cv != p.Value {
			result.Drifted = append(result.Drifted, DriftEntry{
				Key:     p.Key,
				Pinned:  p.Value,
				Current: cv,
			})
		}
	}

	for _, c := range current {
		if _, ok := pinnedMap[c.Key]; !ok {
			result.New = append(result.New, c.Key)
		}
	}

	return result
}

func toMap(entries []envfile.Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}
