package scoper_test

import (
	"testing"

	"github.com/yourusername/envoy-cli/internal/envfile"
	"github.com/yourusername/envoy-cli/internal/scoper"
)

func entries() []envfile.Entry {
	return []envfile.Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "API_KEY", Value: "secret"},
	}
}

func TestScope_AddsPrefix(t *testing.T) {
	opts := scoper.DefaultOptions()
	result := scoper.Scope(entries(), "prod", opts)
	expected := []string{"PROD_DB_HOST", "PROD_DB_PORT", "PROD_API_KEY"}
	for i, e := range result {
		if e.Key != expected[i] {
			t.Errorf("entry %d: got key %q, want %q", i, e.Key, expected[i])
		}
	}
}

func TestScope_PreservesValues(t *testing.T) {
	opts := scoper.DefaultOptions()
	result := scoper.Scope(entries(), "staging", opts)
	if result[0].Value != "localhost" {
		t.Errorf("expected value %q, got %q", "localhost", result[0].Value)
	}
}

func TestScope_EmptyScopeIsNoop(t *testing.T) {
	opts := scoper.DefaultOptions()
	result := scoper.Scope(entries(), "", opts)
	if result[0].Key != "DB_HOST" {
		t.Errorf("expected key unchanged, got %q", result[0].Key)
	}
}

func TestScope_LowercaseScope(t *testing.T) {
	opts := scoper.Options{Separator: "_", UpperCase: false}
	result := scoper.Scope(entries(), "dev", opts)
	if result[0].Key != "dev_DB_HOST" {
		t.Errorf("expected %q, got %q", "dev_DB_HOST", result[0].Key)
	}
}

func TestUnscope_StripsPrefix(t *testing.T) {
	opts := scoper.DefaultOptions()
	scoped := scoper.Scope(entries(), "prod", opts)
	result := scoper.Unscope(scoped, "prod", opts)
	expected := []string{"DB_HOST", "DB_PORT", "API_KEY"}
	for i, e := range result {
		if e.Key != expected[i] {
			t.Errorf("entry %d: got key %q, want %q", i, e.Key, expected[i])
		}
	}
}

func TestUnscope_LeavesNonMatchingKeysUnchanged(t *testing.T) {
	opts := scoper.DefaultOptions()
	input := []envfile.Entry{
		{Key: "PROD_DB_HOST", Value: "localhost"},
		{Key: "STAGING_DB_HOST", Value: "remote"},
	}
	result := scoper.Unscope(input, "prod", opts)
	if result[1].Key != "STAGING_DB_HOST" {
		t.Errorf("non-matching key should be unchanged, got %q", result[1].Key)
	}
}

func TestFilterByScope_ReturnsMatchingEntries(t *testing.T) {
	opts := scoper.DefaultOptions()
	input := []envfile.Entry{
		{Key: "PROD_DB_HOST", Value: "localhost"},
		{Key: "STAGING_DB_HOST", Value: "remote"},
		{Key: "PROD_API_KEY", Value: "secret"},
	}
	result := scoper.FilterByScope(input, "prod", opts)
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result[0].Key != "PROD_DB_HOST" {
		t.Errorf("unexpected key %q", result[0].Key)
	}
}

func TestFilterByScope_EmptyScopeReturnsAll(t *testing.T) {
	opts := scoper.DefaultOptions()
	result := scoper.FilterByScope(entries(), "", opts)
	if len(result) != len(entries()) {
		t.Errorf("expected all entries returned for empty scope")
	}
}

func TestScope_DoesNotMutateOriginal(t *testing.T) {
	opts := scoper.DefaultOptions()
	orig := entries()
	scoper.Scope(orig, "prod", opts)
	if orig[0].Key != "DB_HOST" {
		t.Errorf("original entries were mutated")
	}
}
