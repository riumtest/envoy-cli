package profiler

import (
	"testing"

	"github.com/user/envoy-cli/internal/envfile"
)

func entries() []envfile.Entry {
	return []envfile.Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "APP_PORT", Value: "8080"},
		{Key: "APP_DEBUG", Value: "true"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PASSWORD", Value: "s3cr3t"},
		{Key: "API_KEY", Value: "abc123"},
		{Key: "SITE_URL", Value: "https://example.com"},
		{Key: "EMPTY_VAR", Value: ""},
	}
}

func TestAnalyze_TotalKeys(t *testing.T) {
	p := Analyze(entries())
	if p.TotalKeys != 8 {
		t.Errorf("expected 8 total keys, got %d", p.TotalKeys)
	}
}

func TestAnalyze_EmptyValues(t *testing.T) {
	p := Analyze(entries())
	if p.EmptyValues != 1 {
		t.Errorf("expected 1 empty value, got %d", p.EmptyValues)
	}
}

func TestAnalyze_NumericKeys(t *testing.T) {
	p := Analyze(entries())
	if len(p.NumericKeys) != 1 || p.NumericKeys[0] != "APP_PORT" {
		t.Errorf("expected [APP_PORT], got %v", p.NumericKeys)
	}
}

func TestAnalyze_BooleanKeys(t *testing.T) {
	p := Analyze(entries())
	if len(p.BooleanKeys) != 1 || p.BooleanKeys[0] != "APP_DEBUG" {
		t.Errorf("expected [APP_DEBUG], got %v", p.BooleanKeys)
	}
}

func TestAnalyze_URLKeys(t *testing.T) {
	p := Analyze(entries())
	if len(p.URLKeys) != 1 || p.URLKeys[0] != "SITE_URL" {
		t.Errorf("expected [SITE_URL], got %v", p.URLKeys)
	}
}

func TestAnalyze_SecretKeys(t *testing.T) {
	p := Analyze(entries())
	if len(p.SecretKeys) != 2 {
		t.Errorf("expected 2 secret keys, got %d: %v", len(p.SecretKeys), p.SecretKeys)
	}
}

func TestAnalyze_PrefixGroups(t *testing.T) {
	p := Analyze(entries())
	if p.PrefixGroups["APP"] != 3 {
		t.Errorf("expected APP prefix count 3, got %d", p.PrefixGroups["APP"])
	}
	if p.PrefixGroups["DB"] != 2 {
		t.Errorf("expected DB prefix count 2, got %d", p.PrefixGroups["DB"])
	}
}

func TestAnalyze_EmptyEntries(t *testing.T) {
	p := Analyze([]envfile.Entry{})
	if p.TotalKeys != 0 {
		t.Errorf("expected 0 keys for empty input, got %d", p.TotalKeys)
	}
	if len(p.PrefixGroups) != 0 {
		t.Errorf("expected empty prefix groups, got %v", p.PrefixGroups)
	}
}
