package differ_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourusername/envoy-cli/internal/differ"
	"github.com/yourusername/envoy-cli/internal/envfile"
	"github.com/yourusername/envoy-cli/internal/masker"
)

func TestTextFormatter_NoChanges(t *testing.T) {
	base := []envfile.Entry{{Key: "FOO", Value: "bar"}}
	target := []envfile.Entry{{Key: "FOO", Value: "bar"}}
	result := differ.Compare(base, target)

	f := differ.NewFormatter("text", nil)
	var buf bytes.Buffer
	if err := f.Format(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected no-differences message, got: %q", buf.String())
	}
}

func TestTextFormatter_WithChanges(t *testing.T) {
	base := []envfile.Entry{{Key: "FOO", Value: "old"}, {Key: "GONE", Value: "bye"}}
	target := []envfile.Entry{{Key: "FOO", Value: "new"}, {Key: "ADDED", Value: "hi"}}
	result := differ.Compare(base, target)

	f := differ.NewFormatter("text", nil)
	var buf bytes.Buffer
	_ = f.Format(&buf, result)
	out := buf.String()

	if !strings.Contains(out, "~ FOO") {
		t.Errorf("expected changed FOO, got: %q", out)
	}
	if !strings.Contains(out, "- GONE") {
		t.Errorf("expected removed GONE, got: %q", out)
	}
	if !strings.Contains(out, "+ ADDED") {
		t.Errorf("expected added ADDED, got: %q", out)
	}
}

func TestTextFormatter_WithMasking(t *testing.T) {
	base := []envfile.Entry{{Key: "API_SECRET", Value: "old-secret"}}
	target := []envfile.Entry{{Key: "API_SECRET", Value: "new-secret"}}
	result := differ.Compare(base, target)

	m := masker.New()
	f := differ.NewFormatter("text", m)
	var buf bytes.Buffer
	_ = f.Format(&buf, result)
	out := buf.String()

	if strings.Contains(out, "old-secret") || strings.Contains(out, "new-secret") {
		t.Errorf("expected secrets to be masked, got: %q", out)
	}
	if !strings.Contains(out, "***") {
		t.Errorf("expected mask placeholder, got: %q", out)
	}
}

func TestMaskValue(t *testing.T) {
	base := []envfile.Entry{{Key: "PASSWORD", Value: "hunter2"}}
	target := []envfile.Entry{{Key: "PASSWORD", Value: "hunter3"}}
	result := differ.Compare(base, target)

	m := masker.New()
	f := differ.NewFormatter("text", m)
	var buf bytes.Buffer
	_ = f.Format(&buf, result)

	if strings.Contains(buf.String(), "hunter") {
		t.Errorf("password should be masked")
	}
}

func TestJSONFormatter(t *testing.T) {
	base := []envfile.Entry{{Key: "FOO", Value: "old"}}
	target := []envfile.Entry{{Key: "FOO", Value: "new"}, {Key: "BAR", Value: "added"}}
	result := differ.Compare(base, target)

	f := differ.NewFormatter("json", nil)
	var buf bytes.Buffer
	if err := f.Format(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out []map[string]string
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 changes in JSON, got %d", len(out))
	}
}
