package cmd_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/yourusername/envoy-cli/cmd"
)

func writeRedactTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "envoy-redact-*.env")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func execRedact(t *testing.T, args ...string) (string, error) {
	t.Helper()
	var buf bytes.Buffer
	root := &cobra.Command{Use: "envoy"}
	cmd.RegisterCommands(root)
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs(append([]string{"redact"}, args...))
	err := root.Execute()
	return buf.String(), err
}

func TestRedactCmd_MasksSensitiveKeys(t *testing.T) {
	f := writeRedactTempEnv(t, "DB_PASSWORD=secret\nAPP_ENV=production\n")
	out, err := execRedact(t, f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "DB_PASSWORD=***") {
		t.Errorf("expected DB_PASSWORD to be masked, got:\n%s", out)
	}
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("expected APP_ENV to remain plain, got:\n%s", out)
	}
}

func TestRedactCmd_CustomPlaceholder(t *testing.T) {
	f := writeRedactTempEnv(t, "API_KEY=mykey\n")
	out, err := execRedact(t, "--placeholder", "<hidden>", f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "API_KEY=<hidden>") {
		t.Errorf("expected custom placeholder, got:\n%s", out)
	}
}

func TestRedactCmd_JSONFormat(t *testing.T) {
	f := writeRedactTempEnv(t, "SECRET=abc\nHOST=localhost\n")
	out, err := execRedact(t, "--format", "json", f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "\"SECRET\"") {
		t.Errorf("expected JSON output, got:\n%s", out)
	}
}

func TestRedactCmd_MissingFile(t *testing.T) {
	_, err := execRedact(t, "/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
