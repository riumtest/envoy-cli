// Package rotator provides functionality for rotating secret values
// in .env entries, replacing old values with newly generated or provided ones.
package rotator

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/user/envoy-cli/internal/envfile"
)

// Result holds the outcome of a rotation operation.
type Result struct {
	Key      string
	OldValue string
	NewValue string
	Rotated  bool
}

// Options configures the rotation behaviour.
type Options struct {
	// Keys is the explicit list of keys to rotate. If empty, all sensitive-looking keys are rotated.
	Keys []string
	// Generator is a function that produces a new secret value for a given key.
	// Defaults to a 32-byte hex string generator.
	Generator func(key string) (string, error)
	// DryRun reports what would change without modifying entries.
	DryRun bool
}

// DefaultOptions returns Options with a sensible hex-based secret generator.
func DefaultOptions() Options {
	return Options{
		Generator: defaultGenerator,
	}
}

func defaultGenerator(_ string) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("rotator: failed to generate secret: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// Rotate replaces values for the specified keys (or all keys when Keys is empty)
// using the configured generator. It returns the updated entries and a result
// slice describing each rotation.
func Rotate(entries []envfile.Entry, opts Options) ([]envfile.Entry, []Result, error) {
	allowlist := buildAllowlist(opts.Keys)

	results := make([]Result, 0, len(entries))
	output := make([]envfile.Entry, len(entries))
	copy(output, entries)

	for i, e := range output {
		if len(allowlist) > 0 && !allowlist[e.Key] {
			continue
		}

		newVal, err := opts.Generator(e.Key)
		if err != nil {
			return nil, nil, err
		}

		results = append(results, Result{
			Key:      e.Key,
			OldValue: e.Value,
			NewValue: newVal,
			Rotated:  true,
		})

		if !opts.DryRun {
			output[i].Value = newVal
		}
	}

	return output, results, nil
}

func buildAllowlist(keys []string) map[string]bool {
	if len(keys) == 0 {
		return nil
	}
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		m[k] = true
	}
	return m
}
