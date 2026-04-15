package cipher_test

import (
	"encoding/base64"
	"strings"
	"testing"

	"envoy-cli/internal/cipher"
	"envoy-cli/internal/envfile"
)

func entries() []envfile.Entry {
	return []envfile.Entry{
		{Key: "DB_PASSWORD", Value: "s3cr3t"},
		{Key: "API_KEY", Value: "abc123"},
		{Key: "APP_NAME", Value: "myapp"},
	}
}

// validKey returns a base64-encoded 32-byte AES key.
func validKey() string {
	raw := make([]byte, 32)
	for i := range raw {
		raw[i] = byte(i + 1)
	}
	return base64.StdEncoding.EncodeToString(raw)
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	key := validKey()
	orig := entries()

	enc, err := cipher.Encrypt(orig, key, nil)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	for i, e := range enc {
		if e.Value == orig[i].Value {
			t.Errorf("entry %q was not encrypted", e.Key)
		}
	}

	dec, err := cipher.Decrypt(enc, key, nil)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	for i, e := range dec {
		if e.Value != orig[i].Value {
			t.Errorf("key %q: got %q, want %q", e.Key, e.Value, orig[i].Value)
		}
	}
}

func TestEncrypt_SpecificKeys(t *testing.T) {
	key := validKey()
	orig := entries()

	enc, err := cipher.Encrypt(orig, key, []string{"DB_PASSWORD"})
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if enc[0].Value == orig[0].Value {
		t.Error("DB_PASSWORD should be encrypted")
	}
	if enc[1].Value != orig[1].Value {
		t.Error("API_KEY should be unchanged")
	}
	if enc[2].Value != orig[2].Value {
		t.Error("APP_NAME should be unchanged")
	}
}

func TestDecrypt_SpecificKeys(t *testing.T) {
	key := validKey()
	orig := entries()

	all, _ := cipher.Encrypt(orig, key, nil)
	dec, err := cipher.Decrypt(all, key, []string{"DB_PASSWORD"})
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if dec[0].Value != "s3cr3t" {
		t.Errorf("DB_PASSWORD: got %q, want %q", dec[0].Value, "s3cr3t")
	}
	// API_KEY still encrypted
	if dec[1].Value == "abc123" {
		t.Error("API_KEY should still be encrypted")
	}
}

func TestEncrypt_InvalidKey(t *testing.T) {
	_, err := cipher.Encrypt(entries(), "not-valid-base64!!!", nil)
	if err == nil {
		t.Fatal("expected error for invalid key")
	}
	if !strings.Contains(err.Error(), "cipher") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestEncrypt_WrongKeyLength(t *testing.T) {
	// 10 bytes — not a valid AES key length
	shortKey := base64.StdEncoding.EncodeToString([]byte("tooshort!!"))
	_, err := cipher.Encrypt(entries(), shortKey, nil)
	if err == nil {
		t.Fatal("expected error for wrong key length")
	}
}

func TestDecrypt_TamperedCiphertext(t *testing.T) {
	key := validKey()
	enc, _ := cipher.Encrypt(entries(), key, nil)
	enc[0].Value = base64.StdEncoding.EncodeToString([]byte("garbage"))
	_, err := cipher.Decrypt(enc, key, nil)
	if err == nil {
		t.Fatal("expected error decrypting tampered ciphertext")
	}
}
