package redactor

import (
	"strings"

	"github.com/yourusername/envoy-cli/internal/envfile"
	"github.com/yourusername/envoy-cli/internal/masker"
)

// RedactedEntry holds an env entry with its value optionally redacted.
type RedactedEntry struct {
	Key       string
	Value     string
	Redacted  bool
}

// Options controls redaction behaviour.
type Options struct {
	// Placeholder replaces sensitive values. Defaults to "***".
	Placeholder string
	// ExtraPatterns are additional key patterns to treat as sensitive.
	ExtraPatterns []string
}

// Redact applies secret masking to a slice of parsed entries and returns
// RedactedEntry values, marking which keys were redacted.
func Redact(entries []envfile.Entry, opts Options) []RedactedEntry {
	placeholder := opts.Placeholder
	if placeholder == "" {
		placeholder = "***"
	}

	var m *masker.Masker
	if len(opts.ExtraPatterns) > 0 {
		m = masker.NewWithPatterns(opts.ExtraPatterns)
	} else {
		m = masker.New()
	}

	out := make([]RedactedEntry, 0, len(entries))
	for _, e := range entries {
		re := RedactedEntry{Key: e.Key, Value: e.Value}
		if m.IsSensitive(e.Key) {
			re.Value = placeholder
			re.Redacted = true
		}
		out = append(out, re)
	}
	return out
}

// RedactString replaces sensitive values inside a raw .env file string and
// returns the sanitised content.
func RedactString(raw string, opts Options) string {
	placeholder := opts.Placeholder
	if placeholder == "" {
		placeholder = "***"
	}

	var m *masker.Masker
	if len(opts.ExtraPatterns) > 0 {
		m = masker.NewWithPatterns(opts.ExtraPatterns)
	} else {
		m = masker.New()
	}

	var sb strings.Builder
	for _, line := range strings.Split(raw, "\n") {
		if idx := strings.IndexByte(line, '='); idx > 0 {
			key := strings.TrimSpace(line[:idx])
			if m.IsSensitive(key) {
				sb.WriteString(key + "=" + placeholder + "\n")
				continue
			}
		}
		sb.WriteString(line + "\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}
