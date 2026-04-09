package masker_test

import (
	"testing"

	"envoy-cli/internal/masker"
)

func TestIsSensitive_DefaultPatterns(t *testing.T) {
	m := masker.New()

	sensitive := []string{
		"DB_PASSWORD",
		"API_KEY",
		"SECRET_KEY",
		"AUTH_TOKEN",
		"ACCESS_KEY_ID",
		"PRIVATE_KEY",
		"STRIPE_TOKEN",
		"AWS_SECRET",
		"CREDENTIALS_FILE",
	}

	for _, key := range sensitive {
		if !m.IsSensitive(key) {
			t.Errorf("expected key %q to be sensitive", key)
		}
	}
}

func TestIsSensitive_NotSensitive(t *testing.T) {
	m := masker.New()

	notSensitive := []string{
		"APP_ENV",
		"PORT",
		"LOG_LEVEL",
		"DATABASE_HOST",
		"REDIS_URL",
	}

	for _, key := range notSensitive {
		if m.IsSensitive(key) {
			t.Errorf("expected key %q to NOT be sensitive", key)
		}
	}
}

func TestIsSensitive_CaseInsensitive(t *testing.T) {
	m := masker.New()

	keys := []string{"db_password", "Api_Key", "secret_key"}
	for _, key := range keys {
		if !m.IsSensitive(key) {
			t.Errorf("expected key %q (lowercase) to be sensitive", key)
		}
	}
}

func TestMask_SensitiveValue(t *testing.T) {
	m := masker.New()
	got := m.Mask("DB_PASSWORD", "supersecret")
	if got != "****" {
		t.Errorf("expected masked value, got %q", got)
	}
}

func TestMask_NonSensitiveValue(t *testing.T) {
	m := masker.New()
	got := m.Mask("APP_ENV", "production")
	if got != "production" {
		t.Errorf("expected original value, got %q", got)
	}
}

func TestNewWithPatterns_CustomMask(t *testing.T) {
	m := masker.NewWithPatterns([]string{"CUSTOM_SECRET"}, "[REDACTED]")

	if !m.IsSensitive("MY_CUSTOM_SECRET") {
		t.Error("expected custom pattern to match")
	}

	got := m.Mask("MY_CUSTOM_SECRET", "value")
	if got != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", got)
	}
}
