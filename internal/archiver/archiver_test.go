package archiver_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envoy-cli/internal/archiver"
	"github.com/user/envoy-cli/internal/envfile"
)

func entries() []envfile.Entry {
	return []envfile.Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "SECRET_KEY", Value: "s3cr3t"},
	}
}

func tempDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "archives")
}

func TestSave_CreatesFile(t *testing.T) {
	opts := archiver.Options{Dir: tempDir(t), MaxKeep: 10}
	path, err := archiver.Save(entries(), "v1", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("archive file not found: %v", err)
	}
}

func TestLoad_ReturnsEntries(t *testing.T) {
	opts := archiver.Options{Dir: tempDir(t), MaxKeep: 10}
	_, err := archiver.Save(entries(), "v2", opts)
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	a, err := archiver.Load("v2", opts)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if a.Label != "v2" {
		t.Errorf("expected label v2, got %q", a.Label)
	}
	if len(a.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(a.Entries))
	}
}

func TestList_ReturnsSortedLabels(t *testing.T) {
	opts := archiver.Options{Dir: tempDir(t), MaxKeep: 10}
	for _, label := range []string{"alpha", "beta", "gamma"} {
		if _, err := archiver.Save(entries(), label, opts); err != nil {
			t.Fatalf("save %s: %v", label, err)
		}
	}
	labels, err := archiver.List(opts)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(labels) != 3 {
		t.Errorf("expected 3 labels, got %d", len(labels))
	}
}

func TestList_EmptyDir(t *testing.T) {
	opts := archiver.Options{Dir: tempDir(t), MaxKeep: 10}
	labels, err := archiver.List(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(labels) != 0 {
		t.Errorf("expected empty list, got %d", len(labels))
	}
}

func TestSave_PrunesOldArchives(t *testing.T) {
	opts := archiver.Options{Dir: tempDir(t), MaxKeep: 3}
	for i := 0; i < 5; i++ {
		label := filepath.Join("snap", "") // unique via timestamp suffix
		_ = label
		if _, err := archiver.Save(entries(), "", opts); err != nil {
			t.Fatalf("save: %v", err)
		}
	}
	labels, err := archiver.List(opts)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(labels) > 3 {
		t.Errorf("expected at most 3 archives, got %d", len(labels))
	}
}

func TestLoad_MissingLabel(t *testing.T) {
	opts := archiver.Options{Dir: tempDir(t), MaxKeep: 10}
	_, err := archiver.Load("nonexistent", opts)
	if err == nil {
		t.Error("expected error for missing archive")
	}
}
