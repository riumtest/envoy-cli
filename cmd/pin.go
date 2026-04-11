package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/envoy-cli/envoy/internal/loader"
	"github.com/envoy-cli/envoy/internal/pinner"
	"github.com/spf13/cobra"
)

var pinCmd = &cobra.Command{
	Use:   "pin <pinned.env> <current.env>",
	Short: "Compare a pinned .env baseline against a current file and report drift",
	Args:  cobra.ExactArgs(2),
	RunE:  runPin,
}

var pinFormat string

func init() {
	pinCmd.Flags().StringVarP(&pinFormat, "format", "f", "text", "Output format: text or json")
	rootCmd.AddCommand(pinCmd)
}

func runPin(cmd *cobra.Command, args []string) error {
	pinnedEntries, err := loader.Load(args[0])
	if err != nil {
		return fmt.Errorf("loading pinned file: %w", err)
	}

	currentEntries, err := loader.Load(args[1])
	if err != nil {
		return fmt.Errorf("loading current file: %w", err)
	}

	result := pinner.Pin(pinnedEntries, currentEntries)

	if pinFormat == "json" {
		return printPinJSON(result)
	}
	return printPinText(result)
}

func printPinText(r pinner.Result) error {
	if len(r.Drifted) == 0 && len(r.Missing) == 0 && len(r.New) == 0 {
		fmt.Println("✔ No drift detected.")
		return nil
	}
	for _, d := range r.Drifted {
		fmt.Fprintf(os.Stdout, "~ %s: %q → %q\n", d.Key, d.Pinned, d.Current)
	}
	for _, k := range r.Missing {
		fmt.Fprintf(os.Stdout, "- %s (missing from current)\n", k)
	}
	for _, k := range r.New {
		fmt.Fprintf(os.Stdout, "+ %s (new in current)\n", k)
	}
	fmt.Println(r.Summary())
	return nil
}

func printPinJSON(r pinner.Result) error {
	out := map[string]interface{}{
		"drifted": r.Drifted,
		"missing": r.Missing,
		"new":     r.New,
		"summary": r.Summary(),
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
