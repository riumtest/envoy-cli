package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envoy-cli/internal/loader"
	"github.com/user/envoy-cli/internal/rotator"
)

var rotateCmd = &cobra.Command{
	Use:   "rotate <file>",
	Short: "Rotate secret values in a .env file",
	Args:  cobra.ExactArgs(1),
	RunE:  runRotate,
}

func init() {
	rotateCmd.Flags().StringSliceP("keys", "k", nil, "Comma-separated list of keys to rotate (default: all)")
	rotateCmd.Flags().Bool("dry-run", false, "Preview rotations without writing changes")
	rotateCmd.Flags().StringP("output", "o", "text", "Output format: text or json")
	rotateCmd.Flags().StringP("write", "w", "", "Write updated entries back to this file (defaults to input file)")
	rootCmd.AddCommand(rotateCmd)
}

func runRotate(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	keys, _ := cmd.Flags().GetStringSlice("keys")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	format, _ := cmd.Flags().GetString("output")
	writeTo, _ := cmd.Flags().GetString("write")

	entries, err := loader.Load(filePath)
	if err != nil {
		return fmt.Errorf("rotate: %w", err)
	}

	opts := rotator.DefaultOptions()
	opts.Keys = keys
	opts.DryRun = dryRun

	updated, results, err := rotator.Rotate(entries, opts)
	if err != nil {
		return fmt.Errorf("rotate: %w", err)
	}

	if format == "json" {
		return json.NewEncoder(os.Stdout).Encode(results)
	}

	if len(results) == 0 {
		fmt.Println("No keys rotated.")
		return nil
	}

	for _, r := range results {
		if dryRun {
			fmt.Printf("[dry-run] %s: %s → %s\n", r.Key, r.OldValue, r.NewValue)
		} else {
			fmt.Printf("rotated  %s\n", r.Key)
		}
	}

	if dryRun {
		return nil
	}

	dest := writeTo
	if dest == "" {
		dest = filePath
	}

	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("rotate: cannot write file: %w", err)
	}
	defer f.Close()

	for _, e := range updated {
		fmt.Fprintf(f, "%s=%s\n", e.Key, e.Value)
	}

	fmt.Printf("\nWrote %d entr(ies) to %s\n", len(updated), dest)
	return nil
}
