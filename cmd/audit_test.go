package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func writeAuditTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "audit-*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func execAudit(t *testing.T, args ...string) (string, error) {
	t.Helper()
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(append([]string{"audit"}, args...))
	err := rootCmd.Execute()
	return buf.String(), err
}

func TestAuditCmd_CleanFile(t *testing.T) {
	path := writeAuditTempEnv(t, "APP_NAME=envoy\nPORT=8080\n")
	out, err := execAudit(t, path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No issues") {
		t.Errorf("expected no-issues message, got: %s", out)
	}
}

func TestAuditCmd_EmptyValue(t *testing.T) {
	path := writeAuditTempEnv(t, "MISSING_VAL=\n")
	out, err := execAudit(t, path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "warning") {
		t.Errorf("expected warning in output, got: %s", out)
	}
}

func TestAuditCmd_JSONFormat(t *testing.T) {
	path := writeAuditTempEnv(t, "APP=ok\n")
	out, err := execAudit(t, "--format", "json", path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "\"total\"") {
		t.Errorf("expected JSON output with total field, got: %s", out)
	}
}

func TestAuditCmd_MissingFile(t *testing.T) {
	_, err := execAudit(t, "/nonexistent/path/.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
