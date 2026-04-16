package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/envoy-cli/internal/loader"
	"github.com/user/envoy-cli/internal/pitcher"
)

var pitchCmd = &cobra.Command{
	Use:   "pitch <src> <dst>",
	Short: "Promote env entries from a source file into a destination file",
	Args:  cobra.ExactArgs(2),
	RunE:  runPitch,
}

func init() {
	pitchCmd.Flags().StringSliceP("keys", "k", nil, "Keys to promote (default: all)")
	pitchCmd.Flags().BoolP("no-overwrite", "n", false, "Skip keys already present in destination")
	pitchCmd.Flags().StringP("prefix", "p", "", "Prefix to prepend to promoted keys in destination")
	pitchCmd.Flags().StringP("format", "f", "text", "Output format: text or json")
	pitchCmd.Flags().StringP("output", "o", "", "Write result to file instead of stdout")
	rootCmd.AddCommand(pitchCmd)
}

func runPitch(cmd *cobra.Command, args []string) error {
	srcPath, dstPath := args[0], args[1]
	src, err := loader.Load(srcPath)
	if err != nil {
		return fmt.Errorf("loading source: %w", err)
	}
	dst, err := loader.Load(dstPath)
	if err != nil {
		return fmt.Errorf("loading destination: %w", err)
	}

	keys, _ := cmd.Flags().GetStringSlice("keys")
	noOverwrite, _ := cmd.Flags().GetBool("no-overwrite")
	prefix, _ := cmd.Flags().GetString("prefix")
	format, _ := cmd.Flags().GetString("format")
	output, _ := cmd.Flags().GetString("output")

	opts := pitcher.DefaultOptions()
	opts.Keys = keys
	opts.Overwrite = !noOverwrite
	opts.Prefix = prefix

	res, err := pitcher.Pitch(src, dst, opts)
	if err != nil {
		return err
	}

	w := os.Stdout
	if output != "" {
		f, err := os.Create(output)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	if format == "json" {
		return json.NewEncoder(w).Encode(map[string]interface{}{
			"promoted": res.Promoted,
			"skipped":  res.Skipped,
			"entries":  res.Entries,
		})
	}

	fmt.Fprintf(w, "Promoted: %s\n", strings.Join(res.Promoted, ", "))
	if len(res.Skipped) > 0 {
		fmt.Fprintf(w, "Skipped:  %s\n", strings.Join(res.Skipped, ", "))
	}
	for _, e := range res.Entries {
		fmt.Fprintf(w, "%s=%s\n", e.Key, e.Value)
	}
	return nil
}
