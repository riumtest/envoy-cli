package aliaser_test

import (
	"testing"

	"github.com/envoy-cli/internal/aliaser"
	"github.com/envoy-cli/internal/envfile"
)

func entries() []envfile.Entry {
	return []envfile.Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "API_KEY", Value: "secret"},
	}
}

func TestAlias_SingleRule(t *testing.T) {
	rules := []aliaser.Rule{{From: "DB_HOST", To: "DATABASE_HOST"}}
	res, err := aliaser.Alias(entries(), rules, aliaser.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Aliased) != 1 {
		t.Errorf("expected 1 aliased, got %d", len(res.Aliased))
	}
	found := false
	for _, e := range res.Entries {
		if e.Key == "DATABASE_HOST" && e.Value == "localhost" {
			found = true
		}
	}
	if !found {
		t.Error("alias DATABASE_HOST not found in result entries")
	}
}

func TestAlias_MissingSourceReturnsError(t *testing.T) {
	rules := []aliaser.Rule{{From: "MISSING", To: "ALIAS"}}
	_, err := aliaser.Alias(entries(), rules, aliaser.DefaultOptions())
	if err == nil {
		t.Error("expected error for missing source key")
	}
}

func TestAlias_SkipMissing(t *testing.T) {
	opts := aliaser.DefaultOptions()
	opts.SkipMissing = true
	rules := []aliaser.Rule{{From: "MISSING", To: "ALIAS"}}
	res, err := aliaser.Alias(entries(), rules, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(res.Skipped))
	}
	if len(res.Entries) != len(entries()) {
		t.Error("entries should be unchanged when skipping missing")
	}
}

func TestAlias_ConflictNoOverwrite(t *testing.T) {
	rules := []aliaser.Rule{{From: "DB_HOST", To: "DB_PORT"}}
	res, err := aliaser.Alias(entries(), rules, aliaser.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(res.Conflicts))
	}
	// original value must be preserved
	for _, e := range res.Entries {
		if e.Key == "DB_PORT" && e.Value != "5432" {
			t.Errorf("DB_PORT value should be unchanged, got %q", e.Value)
		}
	}
}

func TestAlias_ConflictWithOverwrite(t *testing.T) {
	opts := aliaser.DefaultOptions()
	opts.Overwrite = true
	rules := []aliaser.Rule{{From: "DB_HOST", To: "DB_PORT"}}
	res, err := aliaser.Alias(entries(), rules, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range res.Entries {
		if e.Key == "DB_PORT" && e.Value != "localhost" {
			t.Errorf("expected DB_PORT overwritten to 'localhost', got %q", e.Value)
		}
	}
}

func TestAlias_MultipleRules(t *testing.T) {
	rules := []aliaser.Rule{
		{From: "DB_HOST", To: "HOST"},
		{From: "API_KEY", To: "SECRET_KEY"},
	}
	res, err := aliaser.Alias(entries(), rules, aliaser.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Aliased) != 2 {
		t.Errorf("expected 2 aliased rules, got %d", len(res.Aliased))
	}
	if len(res.Entries) != len(entries())+2 {
		t.Errorf("expected %d entries, got %d", len(entries())+2, len(res.Entries))
	}
}
