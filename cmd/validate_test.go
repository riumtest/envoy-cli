package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeValidateTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env: %v", err)
	}
	return path
}

func TestValidateCmd_ValidFile(t *testing.T) {
	path := writeValidateTempEnv(t, "APP_ENV=production\nDATABASE_URL=postgres://localhost/db\n")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"validate", path})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "valid") {
		t.Errorf("expected valid message, got: %s", buf.String())
	}
}

func TestValidateCmd_DuplicateKeys(t *testing.T) {
	path := writeValidateTempEnv(t, "APP_ENV=production\nAPP_ENV=staging\n")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"validate", "--warn-only", path})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Issues found") {
		t.Errorf("expected issues output, got: %s", out)
	}
}

func TestValidateCmd_MissingFile(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"validate", "/nonexistent/.env"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestValidateCmd_EmptyValues(t *testing.T) {
	path := writeValidateTempEnv(t, "APP_ENV=\nDATABASE_URL=postgres://localhost/db\n")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"validate", "--warn-only", path})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "WARN") && !strings.Contains(out, "Issues") {
		t.Errorf("expected warning about empty value, got: %s", out)
	}
}
