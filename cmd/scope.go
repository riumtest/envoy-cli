package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/envoy-cli/internal/loader"
	"github.com/yourusername/envoy-cli/internal/scoper"
)

var scopeCmd = &cobra.Command{
	Use:   "scope [file]",
	Short: "Add or remove a scope prefix from .env keys",
	Args:  cobra.ExactArgs(1),
	RunE:  runScope,
}

func init() {
	scopeCmd.Flags().StringP("scope", "s", "", "Scope name to apply or strip (required)")
	scopeCmd.Flags().BoolP("unscope", "u", false, "Strip the scope prefix instead of adding it")
	scopeCmd.Flags().BoolP("filter", "f", false, "Only output entries matching the scope prefix")
	scopeCmd.Flags().StringP("format", "o", "dotenv", "Output format: dotenv or json")
	_ = scopeCmd.MarkFlagRequired("scope")
	rootCmd.AddCommand(scopeCmd)
}

func runScope(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	scope, _ := cmd.Flags().GetString("scope")
	unscope, _ := cmd.Flags().GetBool("unscope")
	filterOnly, _ := cmd.Flags().GetBool("filter")
	format, _ := cmd.Flags().GetString("format")

	entries, err := loader.Load(filePath)
	if err != nil {
		return fmt.Errorf("failed to load file: %w", err)
	}

	opts := scoper.DefaultOptions()

	if filterOnly {
		entries = scoper.FilterByScope(entries, scope, opts)
	} else if unscope {
		entries = scoper.Unscope(entries, scope, opts)
	} else {
		entries = scoper.Scope(entries, scope, opts)
	}

	switch format {
	case "json":
		m := make(map[string]string, len(entries))
		for _, e := range entries {
			m[e.Key] = e.Value
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(m)
	default:
		for _, e := range entries {
			fmt.Fprintf(os.Stdout, "%s=%s\n", e.Key, e.Value)
		}
	}
	return nil
}
