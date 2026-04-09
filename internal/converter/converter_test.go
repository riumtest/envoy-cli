package converter_test

import (
	"strings"
	"testing"

	"github.com/user/envoy-cli/internal/converter"
	"github.com/user/envoy-cli/internal/envfile"
)

func entries() []envfile.Entry {
	return []envfile.Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_PASS", Value: "s3cr3t"},
	}
}

func TestConvert_DotEnv(t *testing.T) {
	out, err := converter.Convert(entries(), converter.FormatDotEnv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("expected dotenv line, got: %s", out)
	}
	if !strings.Contains(out, "DB_PASS=s3cr3t") {
		t.Errorf("expected dotenv line, got: %s", out)
	}
}

func TestConvert_Export(t *testing.T) {
	out, err := converter.Convert(entries(), converter.FormatExport)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "export APP_ENV=") {
		t.Errorf("expected export line, got: %s", out)
	}
}

func TestConvert_JSON(t *testing.T) {
	out, err := converter.Convert(entries(), converter.FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(out, "{") || !strings.HasSuffix(strings.TrimSpace(out), "}") {
		t.Errorf("expected JSON object, got: %s", out)
	}
	if !strings.Contains(out, `"APP_ENV"`) {
		t.Errorf("expected key in JSON, got: %s", out)
	}
}

func TestConvert_YAML(t *testing.T) {
	out, err := converter.Convert(entries(), converter.FormatYAML)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "APP_ENV:") {
		t.Errorf("expected YAML line, got: %s", out)
	}
}

func TestConvert_UnknownFormat(t *testing.T) {
	_, err := converter.Convert(entries(), "toml")
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
	if !strings.Contains(err.Error(), "toml") {
		t.Errorf("error should mention format name, got: %v", err)
	}
}

func TestConvert_EmptyEntries(t *testing.T) {
	out, err := converter.Convert([]envfile.Entry{}, converter.FormatDotEnv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Errorf("expected empty output, got: %q", out)
	}
}
