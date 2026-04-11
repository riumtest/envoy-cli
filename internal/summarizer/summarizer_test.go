package summarizer_test

import (
	"testing"

	"github.com/user/envoy-cli/internal/envfile"
	"github.com/user/envoy-cli/internal/summarizer"
)

func entries(kvs ...string) []envfile.Entry {
	var out []envfile.Entry
	for i := 0; i+1 < len(kvs); i += 2 {
		out = append(out, envfile.Entry{Key: kvs[i], Value: kvs[i+1]})
	}
	return out
}

func TestSummarize_TotalKeys(t *testing.T) {
	e := entries("FOO", "bar", "BAZ", "qux")
	r := summarizer.Summarize(e)
	if r.TotalKeys != 2 {
		t.Errorf("expected 2 total keys, got %d", r.TotalKeys)
	}
}

func TestSummarize_EmptyValues(t *testing.T) {
	e := entries("FOO", "", "BAR", "hello")
	r := summarizer.Summarize(e)
	if r.EmptyValues != 1 {
		t.Errorf("expected 1 empty value, got %d", r.EmptyValues)
	}
}

func TestSummarize_SensitiveKeys(t *testing.T) {
	e := entries("DB_PASSWORD", "secret123", "API_KEY", "abc", "FOO", "bar")
	r := summarizer.Summarize(e)
	if len(r.SensitiveKeys) != 2 {
		t.Errorf("expected 2 sensitive keys, got %d: %v", len(r.SensitiveKeys), r.SensitiveKeys)
	}
}

func TestSummarize_NumericValues(t *testing.T) {
	e := entries("PORT", "8080", "TIMEOUT", "30.5", "NAME", "app")
	r := summarizer.Summarize(e)
	if r.NumericValues != 2 {
		t.Errorf("expected 2 numeric values, got %d", r.NumericValues)
	}
}

func TestSummarize_BooleanValues(t *testing.T) {
	e := entries("DEBUG", "true", "VERBOSE", "false", "PORT", "3000")
	r := summarizer.Summarize(e)
	if r.BooleanValues != 2 {
		t.Errorf("expected 2 boolean values, got %d", r.BooleanValues)
	}
}

func TestSummarize_URLValues(t *testing.T) {
	e := entries("ENDPOINT", "https://api.example.com", "CALLBACK", "http://localhost/cb")
	r := summarizer.Summarize(e)
	if r.URLValues != 2 {
		t.Errorf("expected 2 URL values, got %d", r.URLValues)
	}
}

func TestSummarize_UniqueKeys(t *testing.T) {
	e := entries("FOO", "a", "FOO", "b", "BAR", "c")
	r := summarizer.Summarize(e)
	if r.TotalKeys != 3 {
		t.Errorf("expected 3 total keys, got %d", r.TotalKeys)
	}
	if r.UniqueKeys != 2 {
		t.Errorf("expected 2 unique keys, got %d", r.UniqueKeys)
	}
}

func TestSummarize_EmptyEntries(t *testing.T) {
	r := summarizer.Summarize([]envfile.Entry{})
	if r.TotalKeys != 0 {
		t.Errorf("expected 0 total keys, got %d", r.TotalKeys)
	}
	if len(r.SensitiveKeys) != 0 {
		t.Errorf("expected empty sensitive keys slice, got %v", r.SensitiveKeys)
	}
}
