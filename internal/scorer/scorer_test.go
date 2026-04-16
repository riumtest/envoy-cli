package scorer_test

import (
	"testing"

	"github.com/envoy-cli/internal/envfile"
	"github.com/envoy-cli/internal/scorer"
)

func entries(kvs ...string) []envfile.Entry {
	var out []envfile.Entry
	for i := 0; i+1 < len(kvs); i += 2 {
		out = append(out, envfile.Entry{Key: kvs[i], Value: kvs[i+1]})
	}
	return out
}

func TestScore_PerfectEntries(t *testing.T) {
	res := scorer.Score(entries("FOO", "somevalue"))
	if res[0].Score != scorer.MaxScore {
		t.Fatalf("expected %d got %d notes=%v", scorer.MaxScore, res[0].Score, res[0].Notes)
	}
}

func TestScore_LowercaseKey(t *testing.T) {
	res := scorer.Score(entries("foo", "val"))
	if res[0].Score >= scorer.MaxScore {
		t.Fatal("expected penalty for lowercase key")
	}
}

func TestScore_EmptyValue(t *testing.T) {
	res := scorer.Score(entries("FOO", ""))
	if res[0].Score >= scorer.MaxScore {
		t.Fatal("expected penalty for empty value")
	}
}

func TestScore_SensitiveShortValue(t *testing.T) {
	res := scorer.Score(entries("SECRET", "abc"))
	if res[0].Score >= scorer.MaxScore {
		t.Fatal("expected penalty for short sensitive value")
	}
}

func TestScore_EmptyKey(t *testing.T) {
	res := scorer.Score(entries("", "val"))
	if res[0].Score != 0 {
		t.Fatalf("expected 0 score for empty key, got %d", res[0].Score)
	}
}

func TestScore_MultipleEntries(t *testing.T) {
	res := scorer.Score(entries("GOOD", "value", "bad", ""))
	if len(res) != 2 {
		t.Fatalf("expected 2 results")
	}
	if res[0].Score <= res[1].Score {
		t.Fatalf("expected GOOD > bad score")
	}
}
