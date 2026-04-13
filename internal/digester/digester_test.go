package digester_test

import (
	"testing"

	"github.com/user/envoy-cli/internal/digester"
	"github.com/user/envoy-cli/internal/envfile"
)

func entries(pairs ...string) []envfile.Entry {
	out := make([]envfile.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, envfile.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestDigest_SHA256_Deterministic(t *testing.T) {
	e := entries("HOST", "localhost", "PORT", "5432")
	r1, err := digester.Digest(e, digester.AlgoSHA256)
	if err != nil {
		t.Fatal(err)
	}
	r2, err := digester.Digest(e, digester.AlgoSHA256)
	if err != nil {
		t.Fatal(err)
	}
	if r1.Digest != r2.Digest {
		t.Errorf("expected identical digests, got %q and %q", r1.Digest, r2.Digest)
	}
	if r1.KeyCount != 2 {
		t.Errorf("expected KeyCount 2, got %d", r1.KeyCount)
	}
}

func TestDigest_OrderIndependent(t *testing.T) {
	a := entries("HOST", "localhost", "PORT", "5432")
	b := entries("PORT", "5432", "HOST", "localhost")
	ra, _ := digester.Digest(a, digester.AlgoSHA256)
	rb, _ := digester.Digest(b, digester.AlgoSHA256)
	if ra.Digest != rb.Digest {
		t.Errorf("order should not affect digest: %q != %q", ra.Digest, rb.Digest)
	}
}

func TestDigest_MD5(t *testing.T) {
	e := entries("KEY", "value")
	r, err := digester.Digest(e, digester.AlgoMD5)
	if err != nil {
		t.Fatal(err)
	}
	if r.Algorithm != digester.AlgoMD5 {
		t.Errorf("expected md5 algorithm, got %q", r.Algorithm)
	}
	if len(r.Digest) != 32 {
		t.Errorf("expected 32-char md5 hex digest, got len %d", len(r.Digest))
	}
}

func TestDigest_UnsupportedAlgorithm(t *testing.T) {
	e := entries("K", "v")
	_, err := digester.Digest(e, "blake3")
	if err == nil {
		t.Error("expected error for unsupported algorithm")
	}
}

func TestDigest_EmptyKeySkipped(t *testing.T) {
	e := []envfile.Entry{{Key: "", Value: "ghost"}, {Key: "REAL", Value: "yes"}}
	r, err := digester.Digest(e, digester.AlgoSHA256)
	if err != nil {
		t.Fatal(err)
	}
	if r.KeyCount != 1 {
		t.Errorf("expected KeyCount 1 (empty key skipped), got %d", r.KeyCount)
	}
}

func TestEqual_SameEntries(t *testing.T) {
	a := entries("A", "1", "B", "2")
	b := entries("B", "2", "A", "1")
	ok, err := digester.Equal(a, b, digester.AlgoSHA256)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("expected entries to be equal")
	}
}

func TestEqual_DifferentEntries(t *testing.T) {
	a := entries("A", "1")
	b := entries("A", "2")
	ok, err := digester.Equal(a, b, digester.AlgoSHA256)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Error("expected entries to differ")
	}
}
