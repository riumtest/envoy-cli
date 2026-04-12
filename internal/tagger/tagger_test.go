package tagger_test

import (
	"testing"

	"github.com/envoy-cli/envoy/internal/envfile"
	"github.com/envoy-cli/envoy/internal/tagger"
)

func entries() []envfile.Entry {
	return []envfile.Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PASS", Value: "secret"},
		{Key: "AWS_KEY", Value: "AKIA..."},
		{Key: "AWS_SECRET", Value: "abc123"},
		{Key: "APP_ENV", Value: "production"},
	}
}

func TestTag_NoRules(t *testing.T) {
	result := tagger.Tag(entries(), tagger.DefaultOptions())
	if len(result) != 5 {
		t.Fatalf("expected 5 entries, got %d", len(result))
	}
	for _, r := range result {
		if len(r.Tags) != 0 {
			t.Errorf("expected no tags, got %v for key %s", r.Tags, r.Key)
		}
	}
}

func TestTag_PrefixRule(t *testing.T) {
	opts := tagger.DefaultOptions()
	opts.Rules = map[string]string{
		"database": "DB_",
		"aws":      "AWS_",
	}
	result := tagger.Tag(entries(), opts)

	tagMap := map[string][]string{}
	for _, r := range result {
		tagMap[r.Key] = r.Tags
	}

	if !contains(tagMap["DB_HOST"], "database") {
		t.Errorf("DB_HOST should have 'database' tag")
	}
	if !contains(tagMap["AWS_KEY"], "aws") {
		t.Errorf("AWS_KEY should have 'aws' tag")
	}
	if contains(tagMap["APP_ENV"], "database") || contains(tagMap["APP_ENV"], "aws") {
		t.Errorf("APP_ENV should have no tags")
	}
}

func TestTag_CaseInsensitive(t *testing.T) {
	opts := tagger.DefaultOptions()
	opts.Rules = map[string]string{"database": "db_"}
	opts.CaseSensitive = false

	result := tagger.Tag(entries(), opts)
	for _, r := range result {
		if r.Key == "DB_HOST" && !contains(r.Tags, "database") {
			t.Errorf("DB_HOST should match case-insensitive prefix 'db_'")
		}
	}
}

func TestFilterByTag(t *testing.T) {
	opts := tagger.DefaultOptions()
	opts.Rules = map[string]string{"aws": "AWS_"}
	tagged := tagger.Tag(entries(), opts)

	filtered := tagger.FilterByTag(tagged, "aws")
	if len(filtered) != 2 {
		t.Fatalf("expected 2 aws entries, got %d", len(filtered))
	}
}

func TestToEntries(t *testing.T) {
	tagged := tagger.Tag(entries(), tagger.DefaultOptions())
	out := tagger.ToEntries(tagged)
	if len(out) != len(entries()) {
		t.Errorf("ToEntries length mismatch")
	}
	if out[0].Key != "DB_HOST" {
		t.Errorf("unexpected key: %s", out[0].Key)
	}
}

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
