package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/envoy-cli/internal/converter"
	"github.com/user/envoy-cli/internal/loader"
)

var convertFormat string
var convertOutput string

var convertCmd = &cobra.Command{
	Use:   "convert <file>",
	Short: "Convert a .env file to another format",
	Long: `Convert reads a .env file and outputs its contents in the specified format.

Supported formats: dotenv, export, json, yaml`,
	Args: cobra.ExactArgs(1),
	RunE: runConvert,
}

func init() {
	convertCmd.Flags().StringVarP(&convertFormat, "format", "f", "dotenv", "output format: dotenv, export, json, yaml")
	convertCmd.Flags().StringVarP(&convertOutput, "output", "o", "", "write output to file instead of stdout")
	rootCmd.AddCommand(convertCmd)
}

func runConvert(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	entries, err := loader.Load(filePath)
	if err != nil {
		return fmt.Errorf("failed to load %q: %w", filePath, err)
	}

	fmt := converter.Format(strings.ToLower(convertFormat))
	result, err := converter.Convert(entries, fmt)
	if err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}

	if convertOutput != "" {
		if err := os.WriteFile(convertOutput, []byte(result), 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		cmd.Printf("Written to %s\n", convertOutput)
		return nil
	}

	cmd.Print(result)
	return nil
}
