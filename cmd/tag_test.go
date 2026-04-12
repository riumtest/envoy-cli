package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func writeTagTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func execTag(args ...string) (string, error) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(append([]string{"tag"}, args...))
	err := rootCmd.Execute()
	return buf.String(), err
}

func TestTagCmd_TextOutput(t *testing.T) {
	file := writeTagTempEnv(t, "DB_HOST=localhost\nAWS_KEY=abc\nAPP_ENV=prod\n")
	out, err := execTag("--rule", "db=DB_", "--rule", "aws=AWS_", file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "db") {
		t.Errorf("expected 'db' tag in output, got: %s", out)
	}
}

func TestTagCmd_FilterByTag(t *testing.T) {
	file := writeTagTempEnv(t, "DB_HOST=localhost\nAWS_KEY=abc\nAPP_ENV=prod\n")
	out, err := execTag("--rule", "db=DB_", "--filter", "db", file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "AWS_KEY") {
		t.Errorf("AWS_KEY should be filtered out")
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("DB_HOST should appear in output")
	}
}

func TestTagCmd_JSONFormat(t *testing.T) {
	file := writeTagTempEnv(t, "DB_HOST=localhost\n")
	out, err := execTag("--rule", "db=DB_", "--format", "json", file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "{") {
		t.Errorf("expected JSON output, got: %s", out)
	}
}

func TestTagCmd_MissingFile(t *testing.T) {
	_, err := execTag("/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestTagCmd_InvalidRule(t *testing.T) {
	file := writeTagTempEnv(t, "DB_HOST=localhost\n")
	_, err := execTag("--rule", "badformat", file)
	if err == nil {
		t.Error("expected error for invalid rule format")
	}
}
