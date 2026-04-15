package cmd

import (
	"encoding/base64"
	"crypto/rand"
	"fmt"
	"os"
	"strings"

	"envoy-cli/internal/cipher"
	"envoy-cli/internal/loader"
	"envoy-cli/internal/exporter"
	"github.com/spf13/cobra"
)

var (
	cipherKey     string
	cipherKeys    []string
	cipherDecrypt bool
	cipherOutput  string
	cipherGenKey  bool
)

func init() {
	cipherCmd := &cobra.Command{
		Use:   "cipher [file]",
		Short: "Encrypt or decrypt .env values using AES-GCM",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runCipher,
	}
	cipherCmd.Flags().StringVarP(&cipherKey, "key", "k", "", "Base64-encoded AES key (16/24/32 bytes)")
	cipherCmd.Flags().StringSliceVar(&cipherKeys, "keys", nil, "Specific keys to encrypt/decrypt (default: all)")
	cipherCmd.Flags().BoolVarP(&cipherDecrypt, "decrypt", "d", false, "Decrypt values instead of encrypting")
	cipherCmd.Flags().StringVarP(&cipherOutput, "format", "f", "dotenv", "Output format: dotenv, export, json")
	cipherCmd.Flags().BoolVar(&cipherGenKey, "gen-key", false, "Generate a new random AES-256 key and exit")
	rootCmd.AddCommand(cipherCmd)
}

func runCipher(cmd *cobra.Command, args []string) error {
	if cipherGenKey {
		raw := make([]byte, 32)
		if _, err := rand.Read(raw); err != nil {
			return fmt.Errorf("generate key: %w", err)
		}
		fmt.Println(base64.StdEncoding.EncodeToString(raw))
		return nil
	}

	if len(args) == 0 {
		return fmt.Errorf("cipher: file argument required")
	}
	if cipherKey == "" {
		return fmt.Errorf("cipher: --key is required")
	}

	entries, err := loader.Load(args[0])
	if err != nil {
		return fmt.Errorf("load %q: %w", args[0], err)
	}

	var result []interface{ GetKey() string }
	if cipherDecrypt {
		result2, err := cipher.Decrypt(entries, cipherKey, cipherKeys)
		if err != nil {
			return err
		}
		entries = result2
		_ = result
	} else {
		result2, err := cipher.Encrypt(entries, cipherKey, cipherKeys)
		if err != nil {
			return err
		}
		entries = result2
		_ = result
	}

	fmt := strings.ToLower(cipherOutput)
	return exporter.Write(os.Stdout, entries, exporter.Options{
		Format: fmt,
	})
}
