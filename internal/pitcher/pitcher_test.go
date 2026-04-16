package pitcher_test

import (
	"testing"

	"github.com/user/envoy-cli/internal/envfile"
	"github.com/user/envoy-cli/internal/pitcher"
)

func entries(kv ...string) []envfile.Entry {
	var out []envfile.Entry
	for i := 0; i+1 < len(kv); i += 2 {
		out = append(out, envfile.Entry{Key: kv[i], Value: kv[i+1]})
	}
	return out
}

func TestPitch_AllKeys(t *testing.T) {
	src := entries("FOO", "bar", "BAZ", "qux")
	dst := entries("EXISTING", "val")
	opts := pitcher.DefaultOptions()
	res, err := pitcher.Pitch(src, dst, opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Promoted) != 2 {
		t.Fatalf("expected 2 promoted, got %d", len(res.Promoted))
	}
	if len(res.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(res.Entries))
	}
}

func TestPitch_NoOverwrite(t *testing.T) {
	src := entries("FOO", "new")
	dst := entries("FOO", "old")
	opts := pitcher.DefaultOptions()
	opts.Overwrite = false
	res, err := pitcher.Pitch(src, dst, opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Skipped) != 1 {
		t.Fatalf("expected 1 skipped, got %d", len(res.Skipped))
	}
	for _, e := range res.Entries {
		if e.Key == "FOO" && e.Value != "old" {
			t.Errorf("expected FOO=old, got %s", e.Value)
		}
	}
}

func TestPitch_WithPrefix(t *testing.T) {
	src := entries("DB_HOST", "localhost")
	dst := entries("APP_NAME", "myapp")
	opts := pitcher.DefaultOptions()
	opts.Prefix = "prod"
	res, err := pitcher.Pitch(src, dst, opts)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, e := range res.Entries {
		if e.Key == "PROD_DB_HOST" {
			found = true
		}
	}
	if !found {
		t.Error("expected PROD_DB_HOST in entries")
	}
}

func TestPitch_FilteredKeys(t *testing.T) {
	src := entries("FOO", "1", "BAR", "2", "BAZ", "3")
	dst := entries("EXISTING", "x")
	opts := pitcher.DefaultOptions()
	opts.Keys = []string{"FOO", "BAZ"}
	res, err := pitcher.Pitch(src, dst, opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Promoted) != 2 {
		t.Fatalf("expected 2 promoted, got %d", len(res.Promoted))
	}
}

func TestPitch_MissingFilteredKey(t *testing.T) {
	src := entries("FOO", "1")
	dst := entries()
	opts := pitcher.DefaultOptions()
	opts.Keys = []string{"MISSING"}
	_, err := pitcher.Pitch(src, dst, opts)
	if err == nil {
		t.Error("expected error for missing key")
	}
}

func TestPitch_OverwriteUpdatesValue(t *testing.T) {
	src := entries("FOO", "updated")
	dst := entries("FOO", "original")
	opts := pitcher.DefaultOptions()
	res, err := pitcher.Pitch(src, dst, opts)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range res.Entries {
		if e.Key == "FOO" && e.Value != "updated" {
			t.Errorf("expected FOO=updated, got %s", e.Value)
		}
	}
}
