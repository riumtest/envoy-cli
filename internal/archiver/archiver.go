// Package archiver provides functionality for archiving and restoring
// collections of env entries as timestamped snapshots in a directory.
package archiver

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/user/envoy-cli/internal/envfile"
)

// Archive represents a saved archive entry.
type Archive struct {
	Label     string           `json:"label"`
	CreatedAt time.Time        `json:"created_at"`
	Entries   []envfile.Entry  `json:"entries"`
}

// DefaultOptions returns sensible defaults for archiving.
func DefaultOptions() Options {
	return Options{
		Dir:      ".envoy-archives",
		MaxKeep:  10,
	}
}

// Options configures archive behaviour.
type Options struct {
	Dir     string
	MaxKeep int
}

// Save writes entries to the archive directory with an optional label.
func Save(entries []envfile.Entry, label string, opts Options) (string, error) {
	if err := os.MkdirAll(opts.Dir, 0o755); err != nil {
		return "", fmt.Errorf("archiver: create dir: %w", err)
	}
	if label == "" {
		label = time.Now().UTC().Format("20060102-150405")
	}
	a := Archive{
		Label:     label,
		CreatedAt: time.Now().UTC(),
		Entries:   entries,
	}
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return "", fmt.Errorf("archiver: marshal: %w", err)
	}
	name := fmt.Sprintf("%s.json", label)
	path := filepath.Join(opts.Dir, name)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", fmt.Errorf("archiver: write: %w", err)
	}
	_ = pruneOld(opts)
	return path, nil
}

// Load reads an archived snapshot by label from the archive directory.
func Load(label string, opts Options) (*Archive, error) {
	path := filepath.Join(opts.Dir, label+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("archiver: read %q: %w", path, err)
	}
	var a Archive
	if err := json.Unmarshal(data, &a); err != nil {
		return nil, fmt.Errorf("archiver: unmarshal: %w", err)
	}
	return &a, nil
}

// List returns archive labels sorted by creation time (newest first).
func List(opts Options) ([]string, error) {
	entries, err := os.ReadDir(opts.Dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("archiver: list: %w", err)
	}
	var labels []string
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".json" {
			labels = append(labels, e.Name()[:len(e.Name())-5])
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(labels)))
	return labels, nil
}

// pruneOld removes oldest archives beyond MaxKeep.
func pruneOld(opts Options) error {
	labels, err := List(opts)
	if err != nil || len(labels) <= opts.MaxKeep {
		return err
	}
	for _, l := range labels[opts.MaxKeep:] {
		_ = os.Remove(filepath.Join(opts.Dir, l+".json"))
	}
	return nil
}
