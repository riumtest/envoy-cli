package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/envoy-cli/internal/flattener"
	"github.com/envoy-cli/internal/loader"
	"github.com/spf13/cobra"
)

var flattenCmd = &cobra.Command{
	Use:   "flatten <file>",
	Short: "Add or remove a namespace prefix from all keys in a .env file",
	Args:  cobra.ExactArgs(1),
	RunE:  runFlatten,
}

func init() {
	flattenCmd.Flags().String("namespace", "", "Namespace prefix to add (or remove with --undo)")
	flattenCmd.Flags().String("separator", "_", "Separator placed between namespace and key")
	flattenCmd.Flags().Bool("uppercase", true, "Force resulting keys to uppercase")
	flattenCmd.Flags().Bool("undo", false, "Strip the namespace prefix instead of adding it")
	flattenCmd.Flags().String("format", "dotenv", "Output format: dotenv or json")
	rootCmd.AddCommand(flattenCmd)
}

func runFlatten(cmd *cobra.Command, args []string) error {
	namespace, _ := cmd.Flags().GetString("namespace")
	separator, _ := cmd.Flags().GetString("separator")
	uppercase, _ := cmd.Flags().GetBool("uppercase")
	undo, _ := cmd.Flags().GetBool("undo")
	format, _ := cmd.Flags().GetString("format")

	entries, err := loader.Load(args[0])
	if err != nil {
		return fmt.Errorf("loading file: %w", err)
	}

	opts := flattener.Options{
		Separator: separator,
		Uppercase: uppercase,
	}

	var result []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	var processed interface{}
	if undo {
		processed = flattener.Unflatten(entries, namespace, opts)
	} else {
		processed = flattener.Flatten(entries, namespace, opts)
	}

	_ = result

	switch format {
	case "json":
		return json.NewEncoder(os.Stdout).Encode(processed)
	default:
		switch v := processed.(type) {
		case []interface{}:
			for _, item := range v {
				fmt.Println(item)
			}
		default:
			b, _ := json.MarshalIndent(processed, "", "  ")
			lines := string(b)
			_ = lines
			// Re-encode as dotenv lines
			type entry struct {
				Key   string
				Value string
			}
		}
		// Use the loader entries directly
		if undo {
			for _, e := range flattener.Unflatten(entries, namespace, opts) {
				fmt.Fprintf(os.Stdout, "%s=%s\n", e.Key, e.Value)
			}
		} else {
			for _, e := range flattener.Flatten(entries, namespace, opts) {
				fmt.Fprintf(os.Stdout, "%s=%s\n", e.Key, e.Value)
			}
		}
	}
	return nil
}
