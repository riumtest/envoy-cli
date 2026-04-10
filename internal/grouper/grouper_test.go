package grouper_test

import (
	"testing"

	"github.com/envoy-cli/internal/envfile"
	"github.com/envoy-cli/internal/grouper"
)

func entries() []envfile.Entry {
	return []envfile.Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "REDIS_HOST", Value: "redis"},
		{Key: "REDIS_PORT", Value: "6379"},
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "STANDALONE", Value: "yes"},
	}
}

func TestGroupByPrefix_BasicGroups(t *testing.T) {
	groups := grouper.GroupByPrefix(entries(), grouper.DefaultOptions())
	m := grouper.ToMap(groups)

	if len(m["DB"]) != 2 {
		t.Errorf("expected 2 DB entries, got %d", len(m["DB"]))
	}
	if len(m["REDIS"]) != 2 {
		t.Errorf("expected 2 REDIS entries, got %d", len(m["REDIS"]))
	}
	if len(m["APP"]) != 1 {
		t.Errorf("expected 1 APP entry, got %d", len(m["APP"]))
	}
	if len(m["_other"]) != 1 {
		t.Errorf("expected 1 _other entry, got %d", len(m["_other"]))
	}
}

func TestGroupByPrefix_MinSize(t *testing.T) {
	opts := grouper.DefaultOptions()
	opts.MinSize = 2
	groups := grouper.GroupByPrefix(entries(), opts)
	m := grouper.ToMap(groups)

	if _, ok := m["APP"]; ok {
		t.Error("APP group should be excluded due to MinSize=2")
	}
	if _, ok := m["_other"]; ok {
		t.Error("_other group should be excluded due to MinSize=2")
	}
	if len(m["DB"]) != 2 {
		t.Errorf("expected DB group to remain, got %d entries", len(m["DB"]))
	}
}

func TestGroupByPrefix_SortedGroups(t *testing.T) {
	opts := grouper.DefaultOptions()
	opts.SortGroups = true
	groups := grouper.GroupByPrefix(entries(), opts)

	for i := 1; i < len(groups); i++ {
		if groups[i-1].Name > groups[i].Name {
			t.Errorf("groups not sorted: %s > %s", groups[i-1].Name, groups[i].Name)
		}
	}
}

func TestGroupByPrefix_EmptyEntries(t *testing.T) {
	groups := grouper.GroupByPrefix([]envfile.Entry{}, grouper.DefaultOptions())
	if len(groups) != 0 {
		t.Errorf("expected 0 groups for empty input, got %d", len(groups))
	}
}

func TestGroupByPrefix_CustomDelimiter(t *testing.T) {
	es := []envfile.Entry{
		{Key: "prod.host", Value: "h1"},
		{Key: "prod.port", Value: "80"},
		{Key: "dev.host", Value: "h2"},
	}
	opts := grouper.DefaultOptions()
	opts.Delimiter = "."
	groups := grouper.GroupByPrefix(es, opts)
	m := grouper.ToMap(groups)

	if len(m["prod"]) != 2 {
		t.Errorf("expected 2 prod entries, got %d", len(m["prod"]))
	}
	if len(m["dev"]) != 1 {
		t.Errorf("expected 1 dev entry, got %d", len(m["dev"]))
	}
}
