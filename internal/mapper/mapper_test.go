package mapper_test

import (
	"testing"

	"github.com/envoy-cli/envoy-cli/internal/envfile"
	"github.com/envoy-cli/envoy-cli/internal/mapper"
)

func entries() []envfile.Entry {
	return []envfile.Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "EMPTY_KEY", Value: ""},
	}
}

func TestToKeyValue(t *testing.T) {
	kv := mapper.ToKeyValue(entries())
	if kv["APP_ENV"] != "production" {
		t.Errorf("expected production, got %s", kv["APP_ENV"])
	}
	if kv["DB_PORT"] != "5432" {
		t.Errorf("expected 5432, got %s", kv["DB_PORT"])
	}
	if len(kv) != 4 {
		t.Errorf("expected 4 entries, got %d", len(kv))
	}
}

func TestToValueKey(t *testing.T) {
	vk := mapper.ToValueKey(entries())
	if vk["production"] != "APP_ENV" {
		t.Errorf("expected APP_ENV, got %s", vk["production"])
	}
	if vk["localhost"] != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", vk["localhost"])
	}
	// empty values should be excluded
	if _, ok := vk[""]; ok {
		t.Error("empty value should not appear in inverse map")
	}
}

func TestToValueKey_LastWins(t *testing.T) {
	dupes := []envfile.Entry{
		{Key: "FIRST", Value: "shared"},
		{Key: "SECOND", Value: "shared"},
	}
	vk := mapper.ToValueKey(dupes)
	if vk["shared"] != "SECOND" {
		t.Errorf("expected SECOND to win, got %s", vk["shared"])
	}
}

func TestToKeyIndex(t *testing.T) {
	ki := mapper.ToKeyIndex(entries())
	if ki["APP_ENV"] != 0 {
		t.Errorf("expected index 0, got %d", ki["APP_ENV"])
	}
	if ki["EMPTY_KEY"] != 3 {
		t.Errorf("expected index 3, got %d", ki["EMPTY_KEY"])
	}
}

func TestMap_AllFields(t *testing.T) {
	res := mapper.Map(entries())
	if res.KeyValue["DB_HOST"] != "localhost" {
		t.Errorf("KeyValue mismatch")
	}
	if res.ValueKey["5432"] != "DB_PORT" {
		t.Errorf("ValueKey mismatch")
	}
	if res.KeyIndex["DB_PORT"] != 2 {
		t.Errorf("KeyIndex mismatch, got %d", res.KeyIndex["DB_PORT"])
	}
}

func TestMap_EmptyInput(t *testing.T) {
	res := mapper.Map(nil)
	if len(res.KeyValue) != 0 || len(res.ValueKey) != 0 || len(res.KeyIndex) != 0 {
		t.Error("expected all maps to be empty for nil input")
	}
}
