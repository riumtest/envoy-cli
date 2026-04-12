package extractor_test

import (
	"testing"

	"github.com/user/envoy-cli/internal/envfile"
	"github.com/user/envoy-cli/internal/extractor"
)

func entries() []envfile.Entry {
	return []envfile.Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PASSWORD", Value: "s3cr3t"},
		{Key: "FEATURE_FLAG", Value: "true"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestExtract_NoOptions(t *testing.T) {
	res, err := extractor.Extract(entries(), extractor.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Extracted != 6 {
		t.Errorf("expected 6 extracted, got %d", res.Extracted)
	}
	if res.Skipped != 0 {
		t.Errorf("expected 0 skipped, got %d", res.Skipped)
	}
}

func TestExtract_ByKeyPattern(t *testing.T) {
	res, err := extractor.Extract(entries(), extractor.Options{KeyPattern: "^APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Extracted != 2 {
		t.Errorf("expected 2 extracted, got %d", res.Extracted)
	}
	for _, e := range res.Entries {
		if e.Key != "APP_NAME" && e.Key != "APP_ENV" {
			t.Errorf("unexpected key %q", e.Key)
		}
	}
}

func TestExtract_ByValuePattern(t *testing.T) {
	res, err := extractor.Extract(entries(), extractor.Options{ValuePattern: "^[0-9]+$"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Extracted != 1 {
		t.Errorf("expected 1 extracted, got %d", res.Extracted)
	}
	if res.Entries[0].Key != "PORT" {
		t.Errorf("expected PORT, got %q", res.Entries[0].Key)
	}
}

func TestExtract_ByExplicitKeys(t *testing.T) {
	res, err := extractor.Extract(entries(), extractor.Options{Keys: []string{"DB_HOST", "PORT"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Extracted != 2 {
		t.Errorf("expected 2 extracted, got %d", res.Extracted)
	}
}

func TestExtract_CaseInsensitiveByDefault(t *testing.T) {
	res, err := extractor.Extract(entries(), extractor.Options{KeyPattern: "^app_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Extracted != 2 {
		t.Errorf("expected 2 extracted with case-insensitive match, got %d", res.Extracted)
	}
}

func TestExtract_CaseSensitive(t *testing.T) {
	res, err := extractor.Extract(entries(), extractor.Options{KeyPattern: "^app_", CaseSensitive: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Extracted != 0 {
		t.Errorf("expected 0 extracted with case-sensitive match, got %d", res.Extracted)
	}
}

func TestExtract_InvalidKeyPattern(t *testing.T) {
	_, err := extractor.Extract(entries(), extractor.Options{KeyPattern: "[invalid"})
	if err == nil {
		t.Error("expected error for invalid key pattern")
	}
}

func TestExtract_SkippedCount(t *testing.T) {
	res, err := extractor.Extract(entries(), extractor.Options{KeyPattern: "^DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Skipped != 4 {
		t.Errorf("expected 4 skipped, got %d", res.Skipped)
	}
}
