package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writePitchTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestPitchCmd_AllKeys(t *testing.T) {
	src := writePitchTempEnv(t, "FOO=bar\nBAZ=qux\n")
	dst := writePitchTempEnv(t, "EXISTING=val\n")
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"pitch", src, dst})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "FOO") {
		t.Errorf("expected FOO in output, got: %s", out)
	}
}

func TestPitchCmd_NoOverwrite(t *testing.T) {
	src := writePitchTempEnv(t, "FOO=new\n")
	dst := writePitchTempEnv(t, "FOO=old\n")
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"pitch", "--no-overwrite", src, dst})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "Skipped") {
		t.Errorf("expected Skipped in output, got: %s", out)
	}
}

func TestPitchCmd_JSONFormat(t *testing.T) {
	src := writePitchTempEnv(t, "DB_HOST=localhost\n")
	dst := writePitchTempEnv(t, "APP=myapp\n")
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"pitch", "--format", "json", src, dst})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "promoted") {
		t.Errorf("expected JSON with promoted key, got: %s", buf.String())
	}
}

func TestPitchCmd_MissingFile(t *testing.T) {
	rootCmd.SetArgs([]string{"pitch", "/no/such/src.env", "/no/such/dst.env"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for missing file")
	}
}
