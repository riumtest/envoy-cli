// Package digester computes deterministic hashes (digests) for .env entry
// sets, enabling change detection and integrity verification across environments.
package digester

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/user/envoy-cli/internal/envfile"
)

// Algorithm selects the hashing algorithm used for digest computation.
type Algorithm string

const (
	AlgoSHA256 Algorithm = "sha256"
	AlgoMD5    Algorithm = "md5"
)

// Result holds the computed digest alongside metadata.
type Result struct {
	Algorithm Algorithm `json:"algorithm"`
	Digest    string    `json:"digest"`
	KeyCount  int       `json:"key_count"`
}

// Digest computes a deterministic hash over the provided entries.
// Keys are sorted before hashing so order does not affect the result.
// Only non-empty keys are included in the digest.
func Digest(entries []envfile.Entry, algo Algorithm) (Result, error) {
	pairs := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.Key == "" {
			continue
		}
		pairs = append(pairs, fmt.Sprintf("%s=%s", e.Key, e.Value))
	}
	sort.Strings(pairs)
	raw := strings.Join(pairs, "\n")

	var digest string
	switch algo {
	case AlgoMD5:
		sum := md5.Sum([]byte(raw)) //nolint:gosec
		digest = hex.EncodeToString(sum[:])
	case AlgoSHA256, "":
		sum := sha256.Sum256([]byte(raw))
		digest = hex.EncodeToString(sum[:])
	default:
		return Result{}, fmt.Errorf("unsupported algorithm: %q", algo)
	}

	if algo == "" {
		algo = AlgoSHA256
	}

	return Result{
		Algorithm: algo,
		Digest:    digest,
		KeyCount:  len(pairs),
	}, nil
}

// Equal returns true when two entry sets produce the same digest.
func Equal(a, b []envfile.Entry, algo Algorithm) (bool, error) {
	ra, err := Digest(a, algo)
	if err != nil {
		return false, err
	}
	rb, err := Digest(b, algo)
	if err != nil {
		return false, err
	}
	return ra.Digest == rb.Digest, nil
}
