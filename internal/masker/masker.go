// Package masker provides utilities for detecting and masking sensitive
// environment variable keys based on configurable patterns.
package masker

import (
	"regexp"
	"strings"
)

// DefaultSensitivePatterns is a list of patterns that commonly indicate
// a sensitive environment variable key.
var DefaultSensitivePatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"APIKEY",
	"PRIVATE_KEY",
	"CREDENTIALS",
	"AUTH",
	"ACCESS_KEY",
}

// Masker determines whether a key is sensitive and can mask its value.
type Masker struct {
	patterns []*regexp.Regexp
	maskStr  string
}

// New creates a Masker using the default sensitive patterns and "****" as
// the mask string.
func New() *Masker {
	return NewWithPatterns(DefaultSensitivePatterns, "****")
}

// NewWithPatterns creates a Masker with custom patterns and a custom mask string.
func NewWithPatterns(patterns []string, maskStr string) *Masker {
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re := regexp.MustCompile("(?i)" + regexp.QuoteMeta(p))
		compiled = append(compiled, re)
	}
	return &Masker{patterns: compiled, maskStr: maskStr}
}

// IsSensitive reports whether the given key matches any sensitive pattern.
func (m *Masker) IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, re := range m.patterns {
		if re.MatchString(upper) {
			return true
		}
	}
	return false
}

// Mask returns the masked representation of a value if the key is sensitive,
// otherwise it returns the original value unchanged.
func (m *Masker) Mask(key, value string) string {
	if m.IsSensitive(key) {
		return m.maskStr
	}
	return value
}
