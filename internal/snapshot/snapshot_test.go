package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"envoy-cli/internal/envfile"
	"envoy-cli/internal/snapshot"
)

func testEntries() []envfile.Entry {
	return []envfile.Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "SECRET_KEY", Value: "s3cr3t"},
	}
}

func TestSave_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "snap.json")

	err := snapshot.Save(testEntries(), ".env", "test-label", dest)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if _, err := os.Stat(dest); os.IsNotExist(err) {
		t.Fatal("expected snapshot file to exist")
	}
}

func TestLoad_ReturnsSnapshot(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "snap.json")

	entries := testEntries()
	if err := snapshot.Save(entries, ".env.prod", "v1", dest); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	snap, err := snapshot.Load(dest)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if snap.Label != "v1" {
		t.Errorf("expected label 'v1', got %q", snap.Label)
	}
	if snap.Source != ".env.prod" {
		t.Errorf("expected source '.env.prod', got %q", snap.Source)
	}
	if len(snap.Entries) != len(entries) {
		t.Errorf("expected %d entries, got %d", len(entries), len(snap.Entries))
	}
	if snap.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestSave_DefaultLabel(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "snap.json")

	before := time.Now().Unix()
	if err := snapshot.Save(testEntries(), ".env", "", dest); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	snap, err := snapshot.Load(dest)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if snap.Label == "" {
		t.Error("expected auto-generated label, got empty string")
	}
	_ = before
}

func TestToMap_ReturnsKeyValuePairs(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "snap.json")

	if err := snapshot.Save(testEntries(), ".env", "map-test", dest); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	snap, _ := snapshot.Load(dest)
	m := snap.ToMap()

	if m["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", m["APP_ENV"])
	}
	if m["SECRET_KEY"] != "s3cr3t" {
		t.Errorf("expected SECRET_KEY=s3cr3t, got %q", m["SECRET_KEY"])
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
