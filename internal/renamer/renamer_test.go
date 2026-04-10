package renamer

import (
	"testing"

	"github.com/envoy-cli/internal/envfile"
)

func entries() []envfile.Entry {
	return []envfile.Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "APP_ENV", Value: "production"},
	}
}

func TestRename_SingleRule(t *testing.T) {
	out, res, err := Rename(entries(), []Rule{{From: "DB_HOST", To: "DATABASE_HOST"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Renamed) != 1 {
		t.Fatalf("expected 1 renamed, got %d", len(res.Renamed))
	}
	if out[0].Key != "DATABASE_HOST" {
		t.Errorf("expected DATABASE_HOST, got %s", out[0].Key)
	}
}

func TestRename_MissingFromKey(t *testing.T) {
	_, res, err := Rename(entries(), []Rule{{From: "MISSING_KEY", To: "NEW_KEY"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 {
		t.Fatalf("expected 1 skipped, got %d", len(res.Skipped))
	}
}

func TestRename_ConflictWithExistingKey(t *testing.T) {
	_, res, err := Rename(entries(), []Rule{{From: "DB_HOST", To: "DB_PORT"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflict) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(res.Conflict))
	}
}

func TestRename_MultipleRules(t *testing.T) {
	rules := []Rule{
		{From: "DB_HOST", To: "DATABASE_HOST"},
		{From: "APP_ENV", To: "ENVIRONMENT"},
	}
	out, res, err := Rename(entries(), rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Renamed) != 2 {
		t.Fatalf("expected 2 renamed, got %d", len(res.Renamed))
	}
	keys := map[string]bool{}
	for _, e := range out {
		keys[e.Key] = true
	}
	if !keys["DATABASE_HOST"] || !keys["ENVIRONMENT"] {
		t.Errorf("expected renamed keys in output, got %v", keys)
	}
}

func TestRename_EmptyRules(t *testing.T) {
	out, res, err := Rename(entries(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Renamed) != 0 {
		t.Errorf("expected 0 renamed")
	}
	if len(out) != len(entries()) {
		t.Errorf("expected entries unchanged")
	}
}

func TestRename_InvalidRule(t *testing.T) {
	_, _, err := Rename(entries(), []Rule{{From: "", To: "SOMETHING"}})
	if err == nil {
		t.Error("expected error for empty From, got nil")
	}
}
