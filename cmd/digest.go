package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envoy-cli/internal/digester"
	"github.com/user/envoy-cli/internal/loader"
)

var digestCmd = &cobra.Command{
	Use:   "digest <file>",
	Short: "Compute a deterministic hash digest of a .env file",
	Args:  cobra.ExactArgs(1),
	RunE:  runDigest,
}

func init() {
	digestCmd.Flags().StringP("algo", "a", "sha256", "Hash algorithm: sha256 or md5")
	digestCmd.Flags().StringP("format", "f", "text", "Output format: text or json")
	rootCmd.AddCommand(digestCmd)
}

func runDigest(cmd *cobra.Command, args []string) error {
	path := args[0]
	algoStr, _ := cmd.Flags().GetString("algo")
	format, _ := cmd.Flags().GetString("format")

	entries, err := loader.Load(path)
	if err != nil {
		return fmt.Errorf("loading %q: %w", path, err)
	}

	result, err := digester.Digest(entries, digester.Algorithm(algoStr))
	if err != nil {
		return fmt.Errorf("computing digest: %w", err)
	}

	switch format {
	case "json":
		type output struct {
			File      string           `json:"file"`
			Algorithm digester.Algorithm `json:"algorithm"`
			Digest    string           `json:"digest"`
			KeyCount  int              `json:"key_count"`
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(output{
			File:      path,
			Algorithm: result.Algorithm,
			Digest:    result.Digest,
			KeyCount:  result.KeyCount,
		})
	default:
		fmt.Printf("file:      %s\n", path)
		fmt.Printf("algorithm: %s\n", result.Algorithm)
		fmt.Printf("digest:    %s\n", result.Digest)
		fmt.Printf("keys:      %d\n", result.KeyCount)
	}
	return nil
}
