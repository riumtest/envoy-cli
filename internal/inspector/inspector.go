// Package inspector provides functionality for analyzing .env file contents
// and producing a summary report of key statistics.
package inspector

import (
	"net/url"
	"strconv"
	"strings"

	"envoy-cli/internal/envfile"
	"envoy-cli/internal/masker"
)

// Report holds the result of inspecting a set of env entries.
type Report struct {
	TotalKeys    int
	SensitiveKeys []string
	EmptyValues  int
	URLValues    int
	BooleanValues int
	NumericValues int
	Keys         []string
}

// Inspect analyses the provided entries and returns a Report.
// If m is nil, no sensitive-key detection is performed.
func Inspect(entries []envfile.Entry, m *masker.Masker) Report {
	r := Report{}
	for _, e := range entries {
		r.TotalKeys++
		r.Keys = append(r.Keys, e.Key)

		if e.Value == "" {
			r.EmptyValues++
		}
		if isURL(e.Value) {
			r.URLValues++
		}
		if isBoolean(e.Value) {
			r.BooleanValues++
		}
		if isNumeric(e.Value) {
			r.NumericValues++
		}
		if m != nil && m.IsSensitive(e.Key) {
			r.SensitiveKeys = append(r.SensitiveKeys, e.Key)
		}
	}
	return r
}

func isURL(v string) bool {
	if v == "" {
		return false
	}
	u, err := url.ParseRequestURI(v)
	return err == nil && strings.HasPrefix(u.Scheme, "http")
}

func isBoolean(v string) bool {
	switch strings.ToLower(v) {
	case "true", "false", "yes", "no", "1", "0":
		return true
	}
	return false
}

func isNumeric(v string) bool {
	if v == "" {
		return false
	}
	_, err := strconv.ParseFloat(v, 64)
	return err == nil
}
