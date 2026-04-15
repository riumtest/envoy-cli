// Package cipher provides AES-GCM encryption and decryption for .env entry values.
//
// It is designed to allow sensitive environment variable values to be stored
// encrypted in version control or secrets stores, and decrypted at runtime
// using a symmetric key.
//
// Usage:
//
//	key := base64.StdEncoding.EncodeToString(myAES32ByteKey)
//	encrypted, err := cipher.Encrypt(entries, key, []string{"DB_PASSWORD", "API_KEY"})
//	decrypted, err := cipher.Decrypt(encrypted, key, nil)
//
// Keys must be base64-encoded AES keys of 16, 24, or 32 bytes (AES-128, AES-192, AES-256).
// Encrypted values are stored as base64-encoded AES-GCM ciphertext with a prepended nonce.
package cipher
