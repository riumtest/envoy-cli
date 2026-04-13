package counter_test

import (
	"testing"

	"envoy-cli/internal/counter"
	"envoy-cli/internal/envfile"
)

func entries(pairs ...string) []envfile.Entry {
	var out []envfile.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, envfile.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestCount_Total(t *testing.T) {
	e := entries("DB_HOST", "localhost", "DB_PORT", "5432", "APP_ENV", "prod")
	r := counter.Count(e)
	if r.Total != 3 {
		t.Errorf("expected Total=3, got %d", r.Total)
	}
}

func TestCount_EmptyAndNonEmpty(t *testing.T) {
	e := entries("KEY_A", "value", "KEY_B", "", "KEY_C", "")
	r := counter.Count(e)
	if r.Empty != 2 {
		t.Errorf("expected Empty=2, got %d", r.Empty)
	}
	if r.NonEmpty != 1 {
		t.Errorf("expected NonEmpty=1, got %d", r.NonEmpty)
	}
}

func TestCount_UniqueAndDuplicated(t *testing.T) {
	e := []envfile.Entry{
		{Key: "FOO", Value: "1"},
		{Key: "FOO", Value: "2"},
		{Key: "BAR", Value: "3"},
	}
	r := counter.Count(e)
	if r.Unique != 1 {
		t.Errorf("expected Unique=1, got %d", r.Unique)
	}
	if r.Duplicated != 1 {
		t.Errorf("expected Duplicated=1, got %d", r.Duplicated)
	}
}

func TestCount_ByPrefix(t *testing.T) {
	e := entries(
		"DB_HOST", "localhost",
		"DB_PORT", "5432",
		"APP_ENV", "prod",
		"NOPREFIX", "val",
	)
	r := counter.Count(e)
	if r.ByPrefix["DB"] != 2 {
		t.Errorf("expected DB prefix count=2, got %d", r.ByPrefix["DB"])
	}
	if r.ByPrefix["APP"] != 1 {
		t.Errorf("expected APP prefix count=1, got %d", r.ByPrefix["APP"])
	}
	if r.ByPrefix["NOPREFIX"] != 1 {
		t.Errorf("expected NOPREFIX prefix count=1, got %d", r.ByPrefix["NOPREFIX"])
	}
}

func TestTopPrefixes_Order(t *testing.T) {
	e := entries(
		"DB_HOST", "h",
		"DB_PORT", "p",
		"DB_NAME", "n",
		"APP_ENV", "e",
		"APP_KEY", "k",
		"SVC_URL", "u",
	)
	r := counter.Count(e)
	top := counter.TopPrefixes(r, 2)
	if len(top) != 2 {
		t.Fatalf("expected 2 top prefixes, got %d", len(top))
	}
	if top[0] != "DB" {
		t.Errorf("expected first prefix=DB, got %s", top[0])
	}
	if top[1] != "APP" {
		t.Errorf("expected second prefix=APP, got %s", top[1])
	}
}

func TestCount_EmptyEntries(t *testing.T) {
	r := counter.Count(nil)
	if r.Total != 0 || r.Empty != 0 || r.Unique != 0 {
		t.Errorf("expected all zeros for empty input, got %+v", r)
	}
}
