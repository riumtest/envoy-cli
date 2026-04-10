package cloner_test

import (
	"testing"

	"envoy-cli/internal/cloner"
	"envoy-cli/internal/envfile"
)

func entries(pairs ...string) []envfile.Entry {
	var out []envfile.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, envfile.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestClone_AllKeys(t *testing.T) {
	dst := entries("HOST", "localhost")
	src := entries("PORT", "8080", "DEBUG", "true")

	res, err := cloner.Clone(dst, src, cloner.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Cloned != 2 {
		t.Errorf("expected 2 cloned, got %d", res.Cloned)
	}
	if len(res.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(res.Entries))
	}
}

func TestClone_Overwrite(t *testing.T) {
	dst := entries("HOST", "localhost", "PORT", "3000")
	src := entries("PORT", "8080")

	res, err := cloner.Clone(dst, src, cloner.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Cloned != 1 {
		t.Errorf("expected 1 cloned, got %d", res.Cloned)
	}
	for _, e := range res.Entries {
		if e.Key == "PORT" && e.Value != "8080" {
			t.Errorf("expected PORT=8080, got %s", e.Value)
		}
	}
}

func TestClone_NoOverwrite_Conflict(t *testing.T) {
	dst := entries("PORT", "3000")
	src := entries("PORT", "8080")

	opts := cloner.Options{Overwrite: false}
	res, err := cloner.Clone(dst, src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Conflict != 1 {
		t.Errorf("expected 1 conflict, got %d", res.Conflict)
	}
	if res.Entries[0].Value != "3000" {
		t.Errorf("expected original value preserved")
	}
}

func TestClone_FilteredKeys(t *testing.T) {
	dst := entries()
	src := entries("HOST", "localhost", "PORT", "8080", "DEBUG", "true")

	opts := cloner.Options{Keys: []string{"PORT"}, Overwrite: true}
	res, err := cloner.Clone(dst, src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Cloned != 1 {
		t.Errorf("expected 1 cloned, got %d", res.Cloned)
	}
	if res.Skipped != 2 {
		t.Errorf("expected 2 skipped, got %d", res.Skipped)
	}
	if len(res.Entries) != 1 || res.Entries[0].Key != "PORT" {
		t.Errorf("expected only PORT entry")
	}
}

func TestClone_EmptySrc(t *testing.T) {
	dst := entries("HOST", "localhost")
	res, err := cloner.Clone(dst, nil, cloner.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Cloned != 0 {
		t.Errorf("expected 0 cloned")
	}
	if len(res.Entries) != 1 {
		t.Errorf("expected dst unchanged")
	}
}
