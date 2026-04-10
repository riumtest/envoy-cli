// Package encoder provides utilities for encoding env entries
// into various formats suitable for shell evaluation or config injection.
package encoder

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/user/envoy-cli/internal/envfile"
)

// Format represents the target encoding format.
type Format string

const (
	FormatBase64 Format = "base64"
	FormatHex    Format = "hex"
	FormatEscape Format = "escape"
)

// Result holds an encoded entry.
type Result struct {
	Key      string
	Original string
	Encoded  string
	Format   Format
}

// Encode encodes the values of the provided entries using the given format.
// Only non-empty values are encoded; empty values are passed through unchanged.
func Encode(entries []envfile.Entry, format Format) ([]Result, error) {
	results := make([]Result, 0, len(entries))
	for _, e := range entries {
		encoded, err := encodeValue(e.Value, format)
		if err != nil {
			return nil, fmt.Errorf("encoding key %q: %w", e.Key, err)
		}
		results = append(results, Result{
			Key:      e.Key,
			Original: e.Value,
			Encoded:  encoded,
			Format:   format,
		})
	}
	return results, nil
}

// ToEntries converts a slice of Results back into envfile.Entry values
// using the encoded representation as the value.
func ToEntries(results []Result) []envfile.Entry {
	out := make([]envfile.Entry, len(results))
	for i, r := range results {
		out[i] = envfile.Entry{Key: r.Key, Value: r.Encoded}
	}
	return out
}

func encodeValue(value string, format Format) (string, error) {
	if value == "" {
		return "", nil
	}
	switch format {
	case FormatBase64:
		return base64.StdEncoding.EncodeToString([]byte(value)), nil
	case FormatHex:
		return toHex(value), nil
	case FormatEscape:
		return escapeValue(value), nil
	default:
		return "", fmt.Errorf("unsupported format: %q", format)
	}
}

func toHex(s string) string {
	var sb strings.Builder
	for _, b := range []byte(s) {
		fmt.Fprintf(&sb, "%02x", b)
	}
	return sb.String()
}

func escapeValue(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	s = strings.ReplaceAll(s, "\t", `\t`)
	return s
}
