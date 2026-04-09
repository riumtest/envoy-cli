package exporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envoy-cli/internal/envfile"
	"github.com/user/envoy-cli/internal/exporter"
	"github.com/user/envoy-cli/internal/masker"
)

var testEntries = []envfile.Entry{
	{Key: "APP_NAME", Value: "myapp"},
	{Key: "DB_PASSWORD", Value: "s3cr3t"},
	{Key: "PORT", Value: "8080"},
}

func TestWrite_DotEnv(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Write(&buf, testEntries, exporter.Options{Format: exporter.FormatDotEnv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "APP_NAME=myapp") {
		t.Errorf("expected APP_NAME=myapp in output, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_PASSWORD=s3cr3t") {
		t.Errorf("expected plain password in dotenv output")
	}
}

func TestWrite_DotEnv_Masked(t *testing.T) {
	m := masker.New()
	var buf bytes.Buffer
	err := exporter.Write(&buf, testEntries, exporter.Options{
		Format:        exporter.FormatDotEnv,
		MaskSensitive: true,
		Masker:        m,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "s3cr3t") {
		t.Errorf("expected password to be masked, got:\n%s", out)
	}
	if !strings.Contains(out, "APP_NAME=myapp") {
		t.Errorf("non-sensitive value should not be masked")
	}
}

func TestWrite_Export(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Write(&buf, testEntries, exporter.Options{Format: exporter.FormatExport})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "export APP_NAME='myapp'") {
		t.Errorf("expected export statement, got:\n%s", out)
	}
}

func TestWrite_JSON(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Write(&buf, testEntries, exporter.Options{Format: exporter.FormatJSON})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"key": "APP_NAME"`) {
		t.Errorf("expected JSON key APP_NAME, got:\n%s", out)
	}
	if !strings.Contains(out, `"value": "myapp"`) {
		t.Errorf("expected JSON value myapp, got:\n%s", out)
	}
}
