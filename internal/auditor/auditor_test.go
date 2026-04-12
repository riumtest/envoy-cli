package auditor

import (
	"testing"

	"github.com/yourusername/envoy-cli/internal/masker"
)

func mkEntries(pairs ...string) []Entry {
	var out []Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestAudit_Clean(t *testing.T) {
	entries := mkEntries("APP_NAME", "envoy", "PORT", "8080")
	m := masker.New()
	report := Audit(entries, m)
	if report.Total != 0 {
		t.Errorf("expected 0 findings, got %d", report.Total)
	}
}

func TestAudit_EmptyValue(t *testing.T) {
	entries := mkEntries("APP_NAME", "")
	report := Audit(entries, nil)
	if report.Warnings != 1 {
		t.Errorf("expected 1 warning, got %d", report.Warnings)
	}
}

func TestAudit_EmptyKey(t *testing.T) {
	entries := []Entry{{Key: "", Value: "something"}}
	report := Audit(entries, nil)
	if report.Errors != 1 {
		t.Errorf("expected 1 error, got %d", report.Errors)
	}
}

func TestAudit_ShortSensitiveValue(t *testing.T) {
	entries := mkEntries("SECRET_KEY", "abc")
	m := masker.New()
	report := Audit(entries, m)
	if report.Warnings == 0 {
		t.Error("expected warning for short sensitive value")
	}
}

func TestAudit_TodoMarker(t *testing.T) {
	entries := mkEntries("DATABASE_URL", "TODO: fill this in")
	report := Audit(entries, nil)
	if report.Infos != 1 {
		t.Errorf("expected 1 info, got %d", report.Infos)
	}
}

func TestAudit_UnderscorePrefix(t *testing.T) {
	entries := mkEntries("_INTERNAL", "value")
	report := Audit(entries, nil)
	if report.Infos != 1 {
		t.Errorf("expected 1 info for underscore key, got %d", report.Infos)
	}
}

func TestAudit_CountsAreAccurate(t *testing.T) {
	entries := []Entry{
		{Key: "", Value: "x"},
		{Key: "EMPTY_VAL", Value: ""},
		{Key: "NOTE", Value: "FIXME later"},
	}
	report := Audit(entries, nil)
	if report.Errors != 1 {
		t.Errorf("expected 1 error, got %d", report.Errors)
	}
	if report.Warnings != 1 {
		t.Errorf("expected 1 warning, got %d", report.Warnings)
	}
	if report.Infos != 1 {
		t.Errorf("expected 1 info, got %d", report.Infos)
	}
	if report.Total != 3 {
		t.Errorf("expected total 3, got %d", report.Total)
	}
}
