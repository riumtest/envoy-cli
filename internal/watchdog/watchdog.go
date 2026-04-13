// Package watchdog provides functionality to watch an .env file for changes
// and report drift against a saved baseline snapshot.
package watchdog

import (
	"fmt"
	"time"

	"github.com/envoy-cli/envoy/internal/differ"
	"github.com/envoy-cli/envoy/internal/envfile"
	"github.com/envoy-cli/envoy/internal/loader"
	"github.com/envoy-cli/envoy/internal/snapshot"
)

// Options configures the watchdog behaviour.
type Options struct {
	// PollInterval is how often the file is checked for changes.
	PollInterval time.Duration
	// MaxDrift is the maximum number of changed keys before an alert fires.
	MaxDrift int
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		PollInterval: 5 * time.Second,
		MaxDrift:     0,
	}
}

// Alert is emitted when drift is detected.
type Alert struct {
	File    string
	Changes []differ.Change
	At      time.Time
}

// Watch polls the given file against the snapshot at snapshotPath and sends
// alerts on the returned channel. Call the returned stop func to terminate.
func Watch(filePath, snapshotPath string, opts Options, alerts chan<- Alert) (stop func(), err error) {
	snap, err := snapshot.Load(snapshotPath)
	if err != nil {
		return nil, fmt.Errorf("watchdog: load snapshot: %w", err)
	}

	baseline := snap.ToMap()

	quit := make(chan struct{})
	go func() {
		ticker := time.NewTicker(opts.PollInterval)
		defer ticker.Stop()
		for {
			select {
			case <-quit:
				return
			case <-ticker.C:
				current, loadErr := loader.Load(filePath)
				if loadErr != nil {
					continue
				}
				var baseEntries []envfile.Entry
				for k, v := range baseline {
					baseEntries = append(baseEntries, envfile.Entry{Key: k, Value: v})
				}
				result := differ.Compare(baseEntries, current)
				var drifted []differ.Change
				for _, c := range result.Changes {
					if c.Type != differ.Unchanged {
						drifted = append(drifted, c)
					}
				}
				if len(drifted) > opts.MaxDrift {
					alerts <- Alert{File: filePath, Changes: drifted, At: time.Now()}
				}
			}
		}
	}()

	return func() { close(quit) }, nil
}
