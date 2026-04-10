package cmd_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/envoy-cli/cmd"
)

func writeGroupTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "group-*.env")
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

func TestGroupCmd_TextOutput(t *testing.T) {
	file := writeGroupTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\nREDIS_HOST=redis\n")
	var buf bytes.Buffer
	cmd.RootCmd().SetOut(&buf)
	cmd.RootCmd().SetArgs([]string{"group", file})
	if err := cmd.RootCmd().Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[DB]") {
		t.Errorf("expected [DB] group in output, got: %s", out)
	}
	if !strings.Contains(out, "[REDIS]") {
		t.Errorf("expected [REDIS] group in output, got: %s", out)
	}
}

func TestGroupCmd_JSONOutput(t *testing.T) {
	file := writeGroupTempEnv(t, "APP_NAME=myapp\nAPP_ENV=prod\nDB_HOST=localhost\n")
	var buf bytes.Buffer
	cmd.RootCmd().SetOut(&buf)
	cmd.RootCmd().SetArgs([]string{"group", file, "--format", "json"})
	if err := cmd.RootCmd().Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"name"`) {
		t.Errorf("expected JSON output with 'name' field, got: %s", out)
	}
}

func TestGroupCmd_MinSize(t *testing.T) {
	file := writeGroupTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\nLONELY=yes\n")
	var buf bytes.Buffer
	cmd.RootCmd().SetOut(&buf)
	cmd.RootCmd().SetArgs([]string{"group", file, "--min-size", "2"})
	if err := cmd.RootCmd().Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "[_other]") {
		t.Errorf("_other group should be filtered out by min-size=2")
	}
}

func TestGroupCmd_MissingFile(t *testing.T) {
	cmd.RootCmd().SetArgs([]string{"group", "/nonexistent/path.env"})
	if err := cmd.RootCmd().Execute(); err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
