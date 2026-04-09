package filter_test

import (
	"testing"

	"github.com/user/envoy-cli/internal/envfile"
	"github.com/user/envoy-cli/internal/filter"
)

func entries() []envfile.Entry {
	return []envfile.Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PASSWORD", Value: ""},
		{Key: "FEATURE_FLAG", Value: "true"},
	}
}

func TestFilter_NoOptions(t *testing.T) {
	result := filter.Filter(entries(), filter.Options{})
	if len(result) != 5 {
		t.Fatalf("expected 5 entries, got %d", len(result))
	}
}

func TestFilter_ByPrefix(t *testing.T) {
	result := filter.Filter(entries(), filter.Options{Prefix: "APP_"})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	for _, e := range result {
		if e.Key != "APP_NAME" && e.Key != "APP_ENV" {
			t.Errorf("unexpected key %q", e.Key)
		}
	}
}

func TestFilter_ByKeyContains(t *testing.T) {
	result := filter.Filter(entries(), filter.Options{KeyContains: "DB"})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestFilter_ExcludeEmpty(t *testing.T) {
	result := filter.Filter(entries(), filter.Options{ExcludeEmpty: true})
	if len(result) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(result))
	}
	for _, e := range result {
		if e.Key == "DB_PASSWORD" {
			t.Error("DB_PASSWORD should have been excluded")
		}
	}
}

func TestFilter_ByAllowlist(t *testing.T) {
	result := filter.Filter(entries(), filter.Options{Keys: []string{"APP_NAME", "FEATURE_FLAG"}})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestFilter_DoesNotMutateOriginal(t *testing.T) {
	orig := entries()
	filter.Filter(orig, filter.Options{Prefix: "DB_", ExcludeEmpty: true})
	if len(orig) != 5 {
		t.Error("original slice was mutated")
	}
}

func TestFilter_CombinedOptions(t *testing.T) {
	result := filter.Filter(entries(), filter.Options{
		Prefix:       "DB_",
		ExcludeEmpty: true,
	})
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if result[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %q", result[0].Key)
	}
}
