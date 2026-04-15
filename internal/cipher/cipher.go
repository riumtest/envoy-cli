// Package cipher provides symmetric encryption and decryption for .env entry values.
// It supports AES-GCM encryption with a base64-encoded key, allowing sensitive
// values to be stored encrypted and decrypted at runtime.
package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"envoy-cli/internal/envfile"
)

// ErrInvalidKey is returned when the provided key is not valid for AES.
var ErrInvalidKey = errors.New("cipher: key must be 16, 24, or 32 bytes after base64 decoding")

// Encrypt encrypts the values of matching entries using AES-GCM.
// keyB64 must be a base64-encoded AES key (16, 24, or 32 bytes).
// Only entries whose keys are in the keys allowlist are encrypted.
// If keys is empty, all entries are encrypted.
func Encrypt(entries []envfile.Entry, keyB64 string, keys []string) ([]envfile.Entry, error) {
	block, err := newBlock(keyB64)
	if err != nil {
		return nil, err
	}
	allowlist := buildSet(keys)
	result := make([]envfile.Entry, len(entries))
	for i, e := range entries {
		if len(allowlist) == 0 || allowlist[e.Key] {
			enc, err := encryptValue(block, e.Value)
			if err != nil {
				return nil, fmt.Errorf("cipher: encrypt key %q: %w", e.Key, err)
			}
			e.Value = enc
		}
		result[i] = e
	}
	return result, nil
}

// Decrypt decrypts the values of matching entries using AES-GCM.
// keyB64 must be the same base64-encoded key used during encryption.
// If keys is empty, all entries are decrypted.
func Decrypt(entries []envfile.Entry, keyB64 string, keys []string) ([]envfile.Entry, error) {
	block, err := newBlock(keyB64)
	if err != nil {
		return nil, err
	}
	allowlist := buildSet(keys)
	result := make([]envfile.Entry, len(entries))
	for i, e := range entries {
		if len(allowlist) == 0 || allowlist[e.Key] {
			dec, err := decryptValue(block, e.Value)
			if err != nil {
				return nil, fmt.Errorf("cipher: decrypt key %q: %w", e.Key, err)
			}
			e.Value = dec
		}
		result[i] = e
	}
	return result, nil
}

func newBlock(keyB64 string) (cipher.Block, error) {
	raw, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidKey, err)
	}
	if l := len(raw); l != 16 && l != 24 && l != 32 {
		return nil, ErrInvalidKey
	}
	return aes.NewCipher(raw)
}

func encryptValue(block cipher.Block, plaintext string) (string, error) {
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	sealed := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

func decryptValue(block cipher.Block, ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("base64 decode: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	ns := gcm.NonceSize()
	if len(data) < ns {
		return "", errors.New("ciphertext too short")
	}
	plain, err := gcm.Open(nil, data[:ns], data[ns:], nil)
	if err != nil {
		return "", fmt.Errorf("aes-gcm open: %w", err)
	}
	return string(plain), nil
}

func buildSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
