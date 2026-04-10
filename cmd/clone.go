package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"envoy-cli/internal/cloner"
	"envoy-cli/internal/loader"

	"github.com/spf13/cobra"
)

var (
	cloneKeys      []string
	cloneOverwrite bool
	cloneFormat    string
)

func init() {
	cloneCmd := &cobra.Command{
		Use:   "clone <src> <dst>",
		Short: "Clone entries from a source .env into a destination .env",
		Args:  cobra.ExactArgs(2),
		RunE:  runClone,
	}

	cloneCmd.Flags().StringSliceVarP(&cloneKeys, "keys", "k", nil, "Comma-separated list of keys to clone (default: all)")
	cloneCmd.Flags().BoolVar(&cloneOverwrite, "overwrite", true, "Overwrite existing keys in destination")
	cloneCmd.Flags().StringVarP(&cloneFormat, "format", "f", "text", "Output format: text or json")

	rootCmd.AddCommand(cloneCmd)
}

func runClone(cmd *cobra.Command, args []string) error {
	srcPath, dstPath := args[0], args[1]

	src, err := loader.Load(srcPath)
	if err != nil {
		return fmt.Errorf("loading source: %w", err)
	}

	var dst []interface{ GetKey() string }
	_ = dst

	dstEntries, err := loader.Load(dstPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("loading destination: %w", err)
	}

	opts := cloner.Options{
		Keys:      cloneKeys,
		Overwrite: cloneOverwrite,
	}

	res, err := cloner.Clone(dstEntries, src, opts)
	if err != nil {
		return fmt.Errorf("cloning: %w", err)
	}

	if cloneFormat == "json" {
		out := map[string]interface{}{
			"cloned":   res.Cloned,
			"skipped":  res.Skipped,
			"conflict": res.Conflict,
			"total":    len(res.Entries),
		}
		return json.NewEncoder(os.Stdout).Encode(out)
	}

	lines := make([]string, 0, len(res.Entries))
	for _, e := range res.Entries {
		lines = append(lines, e.Key+"="+e.Value)
	}
	fmt.Println(strings.Join(lines, "\n"))
	fmt.Fprintf(os.Stderr, "cloned=%d skipped=%d conflict=%d\n", res.Cloned, res.Skipped, res.Conflict)
	return nil
}
