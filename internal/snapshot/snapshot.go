// Package snapshot provides functionality for saving and comparing
// point-in-time captures of .env file state.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"envoy-cli/internal/envfile"
)

// Snapshot represents a saved state of an env file at a point in time.
type Snapshot struct {
	Label     string            `json:"label"`
	CreatedAt time.Time         `json:"created_at"`
	Source    string            `json:"source"`
	Entries   []envfile.Entry   `json:"entries"`
}

// Save writes a snapshot of the given entries to the specified file path.
func Save(entries []envfile.Entry, source, label, dest string) error {
	if label == "" {
		label = fmt.Sprintf("snapshot-%d", time.Now().Unix())
	}

	snap := Snapshot{
		Label:     label,
		CreatedAt: time.Now().UTC(),
		Source:    source,
		Entries:   entries,
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return fmt.Errorf("snapshot: failed to create directory: %w", err)
	}

	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("snapshot: failed to create file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(snap); err != nil {
		return fmt.Errorf("snapshot: failed to encode: %w", err)
	}

	return nil
}

// Load reads a snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: failed to open file: %w", err)
	}
	defer f.Close()

	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return nil, fmt.Errorf("snapshot: failed to decode: %w", err)
	}

	return &snap, nil
}

// ToMap converts snapshot entries into a key→value map.
func (s *Snapshot) ToMap() map[string]string {
	m := make(map[string]string, len(s.Entries))
	for _, e := range s.Entries {
		m[e.Key] = e.Value
	}
	return m
}
