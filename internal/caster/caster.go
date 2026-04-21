// Package caster provides type-casting utilities for .env entry values.
// It attempts to infer and convert string values into their natural Go types
// (bool, int, float64) and returns a typed representation alongside the original.
package caster

import (
	"strconv"
	"strings"

	"github.com/user/envoy-cli/internal/envfile"
)

// Kind represents the inferred type of an env value.
type Kind string

const (
	KindString  Kind = "string"
	KindBool    Kind = "bool"
	KindInt     Kind = "int"
	KindFloat   Kind = "float"
	KindEmpty   Kind = "empty"
)

// CastEntry holds the original entry plus its inferred type and parsed value.
type CastEntry struct {
	envfile.Entry
	Kind   Kind
	Parsed interface{}
}

// Cast iterates over entries and annotates each with its inferred Kind and
// parsed Go value. The original entries are never mutated.
func Cast(entries []envfile.Entry) []CastEntry {
	out := make([]CastEntry, 0, len(entries))
	for _, e := range entries {
		out = append(out, castEntry(e))
	}
	return out
}

func castEntry(e envfile.Entry) CastEntry {
	v := e.Value

	if v == "" {
		return CastEntry{Entry: e, Kind: KindEmpty, Parsed: nil}
	}

	// Boolean check
	lower := strings.ToLower(v)
	if lower == "true" || lower == "false" {
		b, _ := strconv.ParseBool(lower)
		return CastEntry{Entry: e, Kind: KindBool, Parsed: b}
	}

	// Integer check
	if i, err := strconv.ParseInt(v, 10, 64); err == nil {
		return CastEntry{Entry: e, Kind: KindInt, Parsed: i}
	}

	// Float check
	if f, err := strconv.ParseFloat(v, 64); err == nil {
		return CastEntry{Entry: e, Kind: KindFloat, Parsed: f}
	}

	return CastEntry{Entry: e, Kind: KindString, Parsed: v}
}
