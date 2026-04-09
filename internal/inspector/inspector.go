// Package inspector provides functionality for analyzing .env file contents
// and producing a summary report of key types, sensitive fields, and statistics.
package inspector

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/user/envoy-cli/internal/envfile"
	"github.com/user/envoy-cli/internal/masker"
)

// Report holds the results of inspecting an env file.
type Report struct {
	TotalKeys     int
	EmptyValues   int
	SensitiveKeys []string
	URLKeys       []string
	BooleanKeys   []string
	NumericKeys   []string
	Keys          []string
}

// Inspect analyses a slice of envfile.Entry values and returns a Report.
// If m is nil, sensitive key detection is skipped.
func Inspect(entries []envfile.Entry, m *masker.Masker) Report {
	r := Report{}

	for _, e := range entries {
		if e.Key == "" {
			continue
		}

		r.TotalKeys++
		r.Keys = append(r.Keys, e.Key)

		if e.Value == "" {
			r.EmptyValues++
		}

		if m != nil && m.IsSensitive(e.Key) {
			r.SensitiveKeys = append(r.SensitiveKeys, e.Key)
		}

		switch {
		case isURL(e.Value):
			r.URLKeys = append(r.URLKeys, e.Key)
		case isBoolean(e.Value):
			r.BooleanKeys = append(r.BooleanKeys, e.Key)
		case isNumeric(e.Value):
			r.NumericKeys = append(r.NumericKeys, e.Key)
		}
	}

	return r
}

func isURL(v string) bool {
	if !strings.Contains(v, "://") {
		return false
	}
	u, err := url.ParseRequestURI(v)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func isBoolean(v string) bool {
	lower := strings.ToLower(v)
	return lower == "true" || lower == "false" || lower == "yes" || lower == "no" || lower == "1" || lower == "0"
}

func isNumeric(v string) bool {
	if v == "" {
		return false
	}
	_, err := strconv.ParseFloat(v, 64)
	return err == nil
}
