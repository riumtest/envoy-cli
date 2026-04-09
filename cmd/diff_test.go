package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
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

func TestDiffCmd_NoChanges(t *testing.T) {
	base := writeTempEnv(t, "APP_NAME=envoy\nDEBUG=false\n")
	target := writeTempEnv(t, "APP_NAME=envoy\nDEBUG=false\n")

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"diff", base, target})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDiffCmd_WithChanges(t *testing.T) {
	base := writeTempEnv(t, "APP_NAME=envoy\nDEBUG=false\n")
	target := writeTempEnv(t, "APP_NAME=envoy\nDEBUG=true\nNEW_KEY=hello\n")

	rootCmd.SetArgs([]string{"diff", base, target})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDiffCmd_JSONFormat(t *testing.T) {
	base := writeTempEnv(t, "APP=foo\n")
	target := writeTempEnv(t, "APP=bar\n")

	rootCmd.SetArgs([]string{"diff", "--format", "json", base, target})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDiffCmd_MissingFile(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nonexistent.env")
	existing := writeTempEnv(t, "KEY=value\n")

	rootCmd.SetArgs([]string{"diff", missing, existing})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
	if !strings.Contains(err.Error(), "base file") {
		t.Errorf("expected error to mention 'base file', got: %v", err)
	}
}
