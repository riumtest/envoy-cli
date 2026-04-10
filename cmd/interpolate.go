package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/envoy-cli/internal/interpolator"
	"github.com/envoy-cli/internal/loader"
	"github.com/spf13/cobra"
)

var interpolateCmd = &cobra.Command{
	Use:   "interpolate <file>",
	Short: "Expand variable references inside .env values",
	Args:  cobra.ExactArgs(1),
	RunE:  runInterpolate,
}

var interpolateUseEnv bool
var interpolateFormat string

func init() {
	interpolateCmd.Flags().BoolVar(&interpolateUseEnv, "use-env", false, "fall back to host environment for unresolved references")
	interpolateCmd.Flags().StringVar(&interpolateFormat, "format", "dotenv", "output format: dotenv, export, json")
	rootCmd.AddCommand(interpolateCmd)
}

func runInterpolate(cmd *cobra.Command, args []string) error {
	entries, err := loader.Load(args[0])
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}

	var result interpolator.Result
	if interpolateUseEnv {
		env := hostEnvMap()
		result = interpolator.InterpolateWithEnv(entries, env)
	} else {
		result = interpolator.Interpolate(entries)
	}

	if len(result.Unresolved) > 0 {
		fmt.Fprintf(os.Stderr, "warning: unresolved references in keys: %v\n", result.Unresolved)
	}

	switch interpolateFormat {
	case "jsonurn printInterpolateJSON(result)
	case "export":
		for _, e := range result.Entries {
			fmt.Printf("export %s=%q\n", e.Key, e.Value)
		}
	default:
		for _, e := range result.Entries {
			fmt.Printf("%s=%s\n", e.Key, e.Value)
		}
	}
	return nil
}

func printInterpolateJSON(r interpolator.Result) error {
	out := make(map[string]string, len(r.Entries))
	for _, e := range r.Entries {
		out[e.Key] = e.Value
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func hostEnvMap() map[string]string {
	m := make(map[string]string)
	for _, pair := range os.Environ() {
		for i := 0; i < len(pair); i++ {
			if pair[i] == '=' {
				m[pair[:i]] = pair[i+1:]
				break
			}
		}
	}
	return m
}
