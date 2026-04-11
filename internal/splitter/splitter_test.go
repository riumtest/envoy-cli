package splitter_test

import (
	"testing"

	"github.com/user/envoy-cli/internal/envfile"
	"github.com/user/envoy-cli/internal/splitter"
)

func entries() []envfile.Entry {
	return []envfile.Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "AWS_ACCESS_KEY", Value: "AKIA123"},
		{Key: "AWS_SECRET", Value: "secret"},
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DEBUG", Value: "true"},
	}
}

func TestSplit_ByPrefix(t *testing.T) {
	opts := splitter.DefaultOptions()
	opts.Prefixes = []string{"DB_", "AWS_"}
	result := splitter.Split(entries(), opts)

	if len(result["DB_"]) != 2 {
		t.Errorf("expected 2 DB_ entries, got %d", len(result["DB_"]))
	}
	if len(result["AWS_"]) != 2 {
		t.Errorf("expected 2 AWS_ entries, got %d", len(result["AWS_"]))
	}
	if len(result[""]) != 2 {
		t.Errorf("expected 2 unmatched entries, got %d", len(result[""]))
	}
}

func TestSplit_StripPrefix(t *testing.T) {
	opts := splitter.DefaultOptions()
	opts.Prefixes = []string{"DB_"}
	opts.StripPrefix = true
	result := splitter.Split(entries(), opts)

	for _, e := range result["DB_"] {
		if e.Key == "DB_HOST" || e.Key == "DB_PORT" {
			t.Errorf("prefix was not stripped: %s", e.Key)
		}
	}
	keys := map[string]bool{}
	for _, e := range result["DB_"] {
		keys[e.Key] = true
	}
	if !keys["HOST"] || !keys["PORT"] {
		t.Errorf("expected HOST and PORT after strip, got %v", keys)
	}
}

func TestSplit_ExcludeUnmatched(t *testing.T) {
	opts := splitter.DefaultOptions()
	opts.Prefixes = []string{"DB_"}
	opts.IncludeUnmatched = false
	result := splitter.Split(entries(), opts)

	if _, ok := result[""]; ok {
		t.Error("expected no unmatched group when IncludeUnmatched is false")
	}
}

func TestSplit_NoPrefixes(t *testing.T) {
	opts := splitter.DefaultOptions()
	opts.Prefixes = []string{}
	result := splitter.Split(entries(), opts)

	if len(result[""]) != len(entries()) {
		t.Errorf("expected all entries in unmatched group, got %d", len(result[""]))
	}
}

func TestSplit_EmptyEntries(t *testing.T) {
	opts := splitter.DefaultOptions()
	opts.Prefixes = []string{"DB_"}
	result := splitter.Split([]envfile.Entry{}, opts)

	if len(result) != 0 {
		t.Errorf("expected empty result for empty input, got %v", result)
	}
}
