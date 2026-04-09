package cmd

import (
	"fmt"
	"os"

	"envoy-cli/internal/loader"
	"envoy-cli/internal/sorter"

	"github.com/spf13/cobra"
)

var (
	sortOrder   string
	sortByValue bool
)

var sortCmd = &cobra.Command{
	Use:   "sort <file>",
	Short: "Sort entries in a .env file by key or value",
	Args:  cobra.ExactArgs(1),
	RunE:  runSort,
}

func init() {
	sortCmd.Flags().StringVarP(&sortOrder, "order", "o", "asc", "Sort order: asc or desc")
	sortCmd.Flags().BoolVarP(&sortByValue, "by-value", "v", false, "Sort by value instead of key")
	rootCmd.AddCommand(sortCmd)
}

func runSort(cmd *cobra.Command, args []string) error {
	entries, err := loader.Load(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return err
	}

	order := sorter.Ascending
	if sortOrder == "desc" {
		order = sorter.Descending
	}

	opts := sorter.Options{
		Order:   order,
		ByValue: sortByValue,
	}

	sorted := sorter.Sort(entries, opts)

	for _, e := range sorted {
		if e.Value == "" {
			fmt.Fprintf(cmd.OutOrStdout(), "%s=\n", e.Key)
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", e.Key, e.Value)
		}
	}

	return nil
}
