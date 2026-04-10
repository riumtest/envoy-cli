package patcher_test

import (
	"testing"

	"github.com/user/envoy-cli/internal/envfile"
	"github.com/user/envoy-cli/internal/patcher"
)

func entries() []envfile.Entry {
	return []envfile.Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "SECRET_KEY", Value: "abc123"},
	}
}

func TestApply_SetExistingKey(t *testing.T) {
	res, err := patcher.Apply(entries(), []patcher.Patch{
		{Op: patcher.OpSet, Key: "APP_ENV", Value: "staging"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Entries[0].Value != "staging" {
		t.Errorf("expected staging, got %s", res.Entries[0].Value)
	}
	if len(res.Applied) != 1 {
		t.Errorf("expected 1 applied, got %d", len(res.Applied))
	}
}

func TestApply_SetNewKey(t *testing.T) {
	res, err := patcher.Apply(entries(), []patcher.Patch{
		{Op: patcher.OpSet, Key: "NEW_KEY", Value: "newval"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 4 {
		t.Errorf("expected 4 entries, got %d", len(res.Entries))
	}
}

func TestApply_DeleteExistingKey(t *testing.T) {
	res, err := patcher.Apply(entries(), []patcher.Patch{
		{Op: patcher.OpDelete, Key: "DB_HOST"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(res.Entries))
	}
	if len(res.Applied) != 1 {
		t.Errorf("expected 1 applied, got %d", len(res.Applied))
	}
}

func TestApply_DeleteMissingKey(t *testing.T) {
	res, err := patcher.Apply(entries(), []patcher.Patch{
		{Op: patcher.OpDelete, Key: "MISSING"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(res.Skipped))
	}
}

func TestApply_RenameKey(t *testing.T) {
	res, err := patcher.Apply(entries(), []patcher.Patch{
		{Op: patcher.OpRename, Key: "APP_ENV", NewKey: "APP_ENVIRONMENT"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Entries[0].Key != "APP_ENVIRONMENT" {
		t.Errorf("expected APP_ENVIRONMENT, got %s", res.Entries[0].Key)
	}
}

func TestApply_UnknownOp(t *testing.T) {
	_, err := patcher.Apply(entries(), []patcher.Patch{
		{Op: "upsert", Key: "FOO"},
	})
	if err == nil {
		t.Error("expected error for unknown op")
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	orig := entries()
	_, err := patcher.Apply(orig, []patcher.Patch{
		{Op: patcher.OpSet, Key: "APP_ENV", Value: "test"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if orig[0].Value != "production" {
		t.Errorf("original mutated: got %s", orig[0].Value)
	}
}
