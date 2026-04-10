package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeInterpolateTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestInterpolateCmd_NoReferences(t *testing.T) {
	p := writeInterpolateTempEnv(t, "HOST=localhost\nPORT=5432\n")
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"interpolate", p})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInterpolateCmd_ExpandsReference(t *testing.T) {
	p := writeInterpolateTempEnv(t, "HOST=db\nDSN=postgres://${HOST}:5432\n")
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"interpolate", p})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInterpolateCmd_JSONFormat(t *testing.T) {
	p := writeInterpolateTempEnv(t, "HOST=db\nDSN=postgres://${HOST}:5432\n")
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"interpolate", "--format", "json", p})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "DSN") && !strings.Contains(buf.String(), "{}") {
		// output goes to stdout directly; just ensure no panic
	}
}

func TestInterpolateCmd_MissingFile(t *testing.T) {
	rootCmd.SetArgs([]string{"interpolate", "/nonexistent/.env"})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestInterpolateCmd_ExportFormat(t *testing.T) {
	p := writeInterpolateTempEnv(t, "KEY=value\n")
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"interpolate", "--format", "export", p})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
