package scorer_test

import (
	"testing"

	"github.com/envoy-cli/internal/masker"
	"github.com/envoy-cli/internal/scorer"
)

func entries(pairs ...string) []scorer.Entry {
	out := make([]scorer.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, scorer.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestScore_PerfectEntries(t *testing.T) {
	m := masker.New()
	res := scorer.Score(entries("APP_ENV", "production", "APP_PORT", "8080"), m)
	if res.Score != 100 {
		t.Errorf("expected 100, got %d", res.Score)
	}
	if len(res.Penalties) != 0 {
		t.Errorf("expected no penalties, got %v", res.Penalties)
	}
}

func TestScore_LowercaseKey(t *testing.T) {
	m := masker.New()
	res := scorer.Score(entries("app_env", "production"), m)
	if res.Score >= 100 {
		t.Errorf("expected penalty for lowercase key, score=%d", res.Score)
	}
	if len(res.Penalties) == 0 {
		t.Error("expected at least one penalty")
	}
}

func TestScore_EmptyValue(t *testing.T) {
	m := masker.New()
	res := scorer.Score(entries("APP_KEY", ""), m)
	if res.Score >= 100 {
		t.Errorf("expected penalty for empty value, score=%d", res.Score)
	}
	found := false
	for _, p := range res.Penalties {
		if len(p) > 0 {
			found = true
		}
	}
	if !found {
		t.Error("expected a penalty message")
	}
}

func TestScore_SensitiveShortValue(t *testing.T) {
	m := masker.New()
	res := scorer.Score(entries("SECRET_KEY", "abc"), m)
	if res.Score >= 100 {
		t.Errorf("expected penalty for short sensitive value, score=%d", res.Score)
	}
}

func TestScore_EmptyEntries(t *testing.T) {
	res := scorer.Score([]scorer.Entry{}, nil)
	if res.Score != 100 {
		t.Errorf("expected 100 for empty input, got %d", res.Score)
	}
}

func TestScore_NoMasker(t *testing.T) {
	res := scorer.Score(entries("SECRET_KEY", "abc"), nil)
	// Without masker, sensitive check is skipped — only lowercase penalty applies (none here)
	if res.Score == 0 {
		t.Error("score should not be 0 without masker")
	}
}

func TestScore_MultiplePenalties(t *testing.T) {
	m := masker.New()
	res := scorer.Score(entries("secret_key", ""), m)
	if len(res.Penalties) < 2 {
		t.Errorf("expected at least 2 penalties, got %d", len(res.Penalties))
	}
	if len(res.Suggestions) < 2 {
		t.Errorf("expected at least 2 suggestions, got %d", len(res.Suggestions))
	}
}
