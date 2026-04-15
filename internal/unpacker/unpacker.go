// Package unpacker expands packed or base64-encoded .env entries back into
// plain key=value pairs, reversing any encoding applied by the encoder package.
package unpacker

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"envoy-cli/internal/envfile"
)

// Format describes the encoding that was applied to values.
type Format string

const (
	FormatBase64 Format = "base64"
	FormatHex    Format = "hex"
	FormatEscape Format = "escape"
)

// Options controls unpacker behaviour.
type Options struct {
	Format   Format
	// Keys restricts unpacking to specific keys; empty means all keys.
	Keys     []string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{Format: FormatBase64}
}

// Unpack decodes values in entries according to opts.
// Entries whose keys are not in opts.Keys (when non-empty) are returned unchanged.
func Unpack(entries []envfile.Entry, opts Options) ([]envfile.Entry, error) {
	allowlist := buildAllowlist(opts.Keys)
	out := make([]envfile.Entry, 0, len(entries))
	for _, e := range entries {
		if len(allowlist) > 0 && !allowlist[e.Key] {
			out = append(out, e)
			continue
		}
		decoded, err := decodeValue(e.Value, opts.Format)
		if err != nil {
			return nil, fmt.Errorf("unpacker: key %q: %w", e.Key, err)
		}
		out = append(out, envfile.Entry{Key: e.Key, Value: decoded})
	}
	return out, nil
}

func decodeValue(v string, f Format) (string, error) {
	switch f {
	case FormatBase64:
		b, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return "", fmt.Errorf("base64 decode: %w", err)
		}
		return string(b), nil
	case FormatHex:
		b, err := hex.DecodeString(v)
		if err != nil {
			return "", fmt.Errorf("hex decode: %w", err)
		}
		return string(b), nil
	case FormatEscape:
		return strings.NewReplacer(
			`\n`, "\n",
			`\t`, "\t",
			`\r`, "\r",
			`\\`, "\\",
		).Replace(v), nil
	default:
		return "", fmt.Errorf("unknown format %q", f)
	}
}

func buildAllowlist(keys []string) map[string]bool {
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		m[k] = true
	}
	return m
}
