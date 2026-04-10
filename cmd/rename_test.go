package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func writeRenameTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
	return f.Name()
}

func TestRenameCmd_SingleRule(t *testing.T) {
	file := writeRenameTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\n")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"rename", "--rule", "DB_HOST=DATABASE_HOST", file})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DATABASE_HOST") {
		t.Errorf("expected DATABASE_HOST in output, got: %s", out)
	}
	if strings.Contains(out, "DB_HOST=") {
		t.Errorf("old key DB_HOST should not appear in output")
	}
}

func TestRenameCmd_MissingFile(t *testing.T) {
	rootCmd.SetArgs([]string{"rename", "--rule", "A=B", "/nonexistent/.env"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for missing file")
	}
}

func TestRenameCmd_InvalidRuleFormat(t *testing.T) {
	file := writeRenameTempEnv(t, "KEY=val\n")
	rootCmd.SetArgs([]string{"rename", "--rule", "NOEQUALS", file})
	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for invalid rule format")
	}
}

func TestRenameCmd_JSONFormat(t *testing.T) {
	file := writeRenameTempEnv(t, "APP_ENV=staging\n")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"rename", "--rule", "APP_ENV=ENVIRONMENT", "--format", "json", file})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "ENVIRONMENT") {
		t.Errorf("expected ENVIRONMENT in JSON output, got: %s", out)
	}
}
