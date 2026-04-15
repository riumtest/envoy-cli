// Package anonymizer replaces identifying values in env entries
// with deterministic, opaque tokens so that structure is preserved
// without leaking real data.
package anonymizer

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/user/envoy-cli/internal/envfile"
	"github.com/user/envoy-cli/internal/masker"
)

// Options controls anonymization behaviour.
type Options struct {
	// Prefix is prepended to every generated token (default: "anon").
	Prefix string
	// SensitiveOnly limits anonymization to keys flagged as sensitive.
	SensitiveOnly bool
	// ExtraPatterns are additional key-name patterns treated as sensitive.
	ExtraPatterns []string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{Prefix: "anon"}
}

// Anonymize replaces entry values with deterministic tokens.
// If opts.SensitiveOnly is true only sensitive keys are anonymized;
// all other values are kept as-is.
func Anonymize(entries []envfile.Entry, opts Options) []envfile.Entry {
	if opts.Prefix == "" {
		opts.Prefix = "anon"
	}

	var m *masker.Masker
	if opts.SensitiveOnly {
		if len(opts.ExtraPatterns) > 0 {
			m = masker.NewWithPatterns(opts.ExtraPatterns)
		} else {
			m = masker.New()
		}
	}

	out := make([]envfile.Entry, len(entries))
	for i, e := range entries {
		copied := e
		if e.Value != "" {
			if m == nil || m.IsSensitive(e.Key) {
				copied.Value = token(opts.Prefix, e.Key, e.Value)
			}
		}
		out[i] = copied
	}
	return out
}

// token builds a short deterministic replacement for a value.
func token(prefix, key, value string) string {
	h := sha256.Sum256([]byte(key + ":" + value))
	return fmt.Sprintf("%s_%s", strings.ToLower(prefix), fmt.Sprintf("%x", h[:4]))
}
