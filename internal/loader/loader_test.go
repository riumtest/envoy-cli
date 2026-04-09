package loader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envoy-cli/internal/loader"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidFile(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nDATABASE_URL=postgres://localhost/db\n")

	ef, err := loader.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ef.Path != path {
		// path may be abs-resolved; check suffix
		if !filepath.IsAbs(ef.Path) {
			t.Errorf("expected absolute path, got %q", ef.Path)
		}
	}
	if len(ef.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(ef.Entries))
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := loader.Load("/nonexistent/path/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_EmptyFile(t *testing.T) {
	path := writeTempEnv(t, "")

	ef, err := loader.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ef.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(ef.Entries))
	}
}

func TestLoadPair_BothValid(t *testing.T) {
	base := writeTempEnv(t, "KEY=base_val\n")
	target := writeTempEnv(t, "KEY=target_val\n")

	b, tgt, err := loader.LoadPair(base, target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(b.Entries) != 1 || len(tgt.Entries) != 1 {
		t.Errorf("expected 1 entry each, got base=%d target=%d", len(b.Entries), len(tgt.Entries))
	}
}

func TestLoadPair_MissingTarget(t *testing.T) {
	base := writeTempEnv(t, "KEY=val\n")

	_, _, err := loader.LoadPair(base, "/no/such/file.env")
	if err == nil {
		t.Fatal("expected error when target is missing")
	}
}

func TestExists(t *testing.T) {
	path := writeTempEnv(t, "X=1\n")

	if !loader.Exists(path) {
		t.Errorf("expected Exists to return true for %q", path)
	}
	if loader.Exists("/no/such/file.env") {
		t.Error("expected Exists to return false for missing file")
	}
}
