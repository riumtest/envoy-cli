// Package loader provides functionality for loading .env files
// from the filesystem, with support for multiple environments and
// path resolution helpers.
package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/envoy-cli/internal/envfile"
)

// EnvFile represents a loaded environment file with its parsed entries.
type EnvFile struct {
	Path    string
	Entries []envfile.Entry
}

// Load reads and parses a .env file from the given path.
// It returns an EnvFile containing the resolved path and parsed entries,
// or an error if the file cannot be read or parsed.
func Load(path string) (*EnvFile, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("loader: resolving path %q: %w", path, err)
	}

	f, err := os.Open(abs)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("loader: file not found: %q", abs)
		}
		return nil, fmt.Errorf("loader: opening file %q: %w", abs, err)
	}
	defer f.Close()

	entries, err := envfile.Parse(f)
	if err != nil {
		return nil, fmt.Errorf("loader: parsing file %q: %w", abs, err)
	}

	return &EnvFile{
		Path:    abs,
		Entries: entries,
	}, nil
}

// LoadPair loads two .env files and returns them as a pair.
// Useful for diff operations that compare two environments.
func LoadPair(basePath, targetPath string) (*EnvFile, *EnvFile, error) {
	base, err := Load(basePath)
	if err != nil {
		return nil, nil, fmt.Errorf("loader: loading base file: %w", err)
	}

	target, err := Load(targetPath)
	if err != nil {
		return nil, nil, fmt.Errorf("loader: loading target file: %w", err)
	}

	return base, target, nil
}

// Exists reports whether the given path points to a readable file.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
