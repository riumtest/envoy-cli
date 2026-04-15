package merger_test

import (
	"testing"

	"github.com/user/envoy-cli/internal/envfile"
	"github.com/user/envoy-cli/internal/merger"
)

func entries(pairs ...string) []envfile.Entry {
	var out []envfile.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, envfile.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestMerge_NoConflicts(t *testing.T) {
	files := map[string][]envfile.Entry{
		"base.env": entries("APP_NAME", "myapp", "DEBUG", "false"),
		"prod.env": entries("PORT", "8080"),
	}
	order := []string{"base.env", "prod.env"}

	result, err := merger.Merge(files, order, merger.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(result.Entries))
	}
	if len(result.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(result.Conflicts))
	}
}

func TestMerge_StrategyFirst(t *testing.T) {
	files := map[string][]envfile.Entry{
		"base.env": entries("DB_HOST", "localhost"),
		"prod.env": entries("DB_HOST", "prod-db.internal"),
	}
	order := []string{"base.env", "prod.env"}

	result, err := merger.Merge(files, order, merger.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Entries[0].Value != "localhost" {
		t.Errorf("StrategyFirst: expected 'localhost', got %q", result.Entries[0].Value)
	}
	if len(result.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(result.Conflicts))
	}
}

func TestMerge_StrategyLast(t *testing.T) {
	files := map[string][]envfile.Entry{
		"base.env": entries("DB_HOST", "localhost"),
		"prod.env": entries("DB_HOST", "prod-db.internal"),
	}
	order := []string{"base.env", "prod.env"}

	result, err := merger.Merge(files, order, merger.StrategyLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Entries[0].Value != "prod-db.internal" {
		t.Errorf("StrategyLast: expected 'prod-db.internal', got %q", result.Entries[0].Value)
	}
}

func TestMerge_StrategyError(t *testing.T) {
	files := map[string][]envfile.Entry{
		"a.env": entries("KEY", "val1"),
		"b.env": entries("KEY", "val2"),
	}
	order := []string{"a.env", "b.env"}

	_, err := merger.Merge(files, order, merger.StrategyError)
	if err == nil {
		t.Fatal("expected error for duplicate key, got nil")
	}
}

func TestMerge_MissingFileInOrder(t *testing.T) {
	files := map[string][]envfile.Entry{
		"base.env": entries("FOO", "bar"),
	}
	order := []string{"base.env", "missing.env"}

	result, err := merger.Merge(files, order, merger.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(result.Entries))
	}
}

func TestMerge_EmptyFiles(t *testing.T) {
	files := map[string][]envfile.Entry{
		"a.env": entries(),
		"b.env": entries(),
	}
	order := []string{"a.env", "b.env"}

	result, err := merger.Merge(files, order, merger.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(result.Entries))
	}
	if len(result.Conflicts) != 0 {
		t.Errorf("expected 0 conflicts, got %d", len(result.Conflicts))
	}
}
