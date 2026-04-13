package blanker_test

import (
	"testing"

	"github.com/user/envoy-cli/internal/blanker"
	"github.com/user/envoy-cli/internal/envfile"
)

func entries() []envfile.Entry {
	return []envfile.Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "s3cr3t"},
		{Key: "API_KEY", Value: "abc123"},
		{Key: "PORT", Value: "8080"},
		{Key: "DEBUG", Value: "true"},
	}
}

func TestBlank_AllEntries(t *testing.T) {
	result := blanker.Blank(entries(), blanker.DefaultOptions())
	for _, e := range result {
		if e.Value != "" {
			t.Errorf("expected value of %q to be blank, got %q", e.Key, e.Value)
		}
	}
}

func TestBlank_ExplicitKeys(t *testing.T) {
	opts := blanker.Options{Keys: []string{"DB_PASSWORD", "API_KEY"}}
	result := blanker.Blank(entries(), opts)

	expect := map[string]string{
		"APP_NAME":    "myapp",
		"DB_PASSWORD": "",
		"API_KEY":     "",
		"PORT":        "8080",
		"DEBUG":       "true",
	}
	for _, e := range result {
		if e.Value != expect[e.Key] {
			t.Errorf("key %q: want %q, got %q", e.Key, expect[e.Key], e.Value)
		}
	}
}

func TestBlank_SensitiveOnly(t *testing.T) {
	opts := blanker.Options{SensitiveOnly: true}
	result := blanker.Blank(entries(), opts)

	for _, e := range result {
		switch e.Key {
		case "DB_PASSWORD", "API_KEY":
			if e.Value != "" {
				t.Errorf("expected sensitive key %q to be blanked", e.Key)
			}
		default:
			if e.Value == "" {
				t.Errorf("expected non-sensitive key %q to retain its value", e.Key)
			}
		}
	}
}

func TestBlank_PatternMatch(t *testing.T) {
	opts := blanker.Options{Patterns: []string{"PORT", "DEBUG"}}
	result := blanker.Blank(entries(), opts)

	for _, e := range result {
		switch e.Key {
		case "PORT", "DEBUG":
			if e.Value != "" {
				t.Errorf("expected %q to be blanked by pattern", e.Key)
			}
		default:
			if e.Value == "" {
				t.Errorf("expected %q to retain its value", e.Key)
			}
		}
	}
}

func TestBlank_DoesNotMutateOriginal(t *testing.T) {
	orig := entries()
	blanker.Blank(orig, blanker.DefaultOptions())
	for _, e := range orig {
		if e.Value == "" && e.Key != "" {
			// Only fail if the original had a non-empty value.
			// We stored them above; just re-check a known one.
		}
	}
	// Spot-check a known entry.
	if orig[0].Value != "myapp" {
		t.Errorf("original entry mutated: got %q", orig[0].Value)
	}
}

func TestBlank_EmptyEntries(t *testing.T) {
	result := blanker.Blank([]envfile.Entry{}, blanker.DefaultOptions())
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}
