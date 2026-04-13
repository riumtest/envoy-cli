package watchdog_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/envoy-cli/envoy/internal/envfile"
	"github.com/envoy-cli/envoy/internal/snapshot"
	"github.com/envoy-cli/envoy/internal/watchdog"
)

func writeEnvFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "watchdog-*.env")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func writeSnapshot(t *testing.T, entries []envfile.Entry) string {
	t.Helper()
	f, err := os.CreateTemp("", "watchdog-*.snap.json")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	if err := snapshot.Save(entries, f.Name(), "test"); err != nil {
		t.Fatal(err)
	}
	return f.Name()
}

func TestWatch_NoDrift(t *testing.T) {
	entries := []envfile.Entry{{Key: "FOO", Value: "bar"}}
	envPath := writeEnvFile(t, "FOO=bar\n")
	snapPath := writeSnapshot(t, entries)

	alerts := make(chan watchdog.Alert, 4)
	opts := watchdog.DefaultOptions()
	opts.PollInterval = 50 * time.Millisecond

	stop, err := watchdog.Watch(envPath, snapPath, opts, alerts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	time.Sleep(160 * time.Millisecond)
	stop()

	if len(alerts) != 0 {
		t.Fatalf("expected no alerts, got %d", len(alerts))
	}
}

func TestWatch_DetectsDrift(t *testing.T) {
	entries := []envfile.Entry{{Key: "FOO", Value: "original"}}
	envPath := writeEnvFile(t, "FOO=changed\n")
	snapPath := writeSnapshot(t, entries)

	alerts := make(chan watchdog.Alert, 4)
	opts := watchdog.DefaultOptions()
	opts.PollInterval = 50 * time.Millisecond

	stop, err := watchdog.Watch(envPath, snapPath, opts, alerts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	time.Sleep(160 * time.Millisecond)
	stop()

	if len(alerts) == 0 {
		t.Fatal("expected at least one alert")
	}
	alert := <-alerts
	if alert.File != envPath {
		t.Errorf("expected file %s, got %s", envPath, alert.File)
	}
	if len(alert.Changes) == 0 {
		t.Error("expected changes in alert")
	}
}

func TestWatch_InvalidSnapshot(t *testing.T) {
	f, _ := os.CreateTemp("", "bad-snap-*.json")
	f.WriteString("not json")
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })

	_, err := watchdog.Watch("any.env", f.Name(), watchdog.DefaultOptions(), make(chan watchdog.Alert))
	if err == nil {
		t.Fatal("expected error for invalid snapshot")
	}
}

func TestAlert_JSONSerializable(t *testing.T) {
	alert := watchdog.Alert{
		File: "test.env",
		At:   time.Now(),
	}
	b, err := json.Marshal(alert)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if len(b) == 0 {
		t.Error("expected non-empty JSON")
	}
}
