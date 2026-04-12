package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envoy-cli/internal/freezer"
	"github.com/user/envoy-cli/internal/loader"
)

var freezeCmd = &cobra.Command{
	Use:   "freeze <baseline> <current>",
	Short: "Detect mutations in a .env file relative to a frozen baseline",
	Args:  cobra.ExactArgs(2),
	RunE:  runFreeze,
}

var freezeJSON bool
var freezeFail bool

func init() {
	rootCmd.AddCommand(freezeCmd)
	freezeCmd.Flags().BoolVar(&freezeJSON, "json", false, "Output violations as JSON")
	freezeCmd.Flags().BoolVar(&freezeFail, "fail", false, "Exit with non-zero status if violations are found")
}

func runFreeze(cmd *cobra.Command, args []string) error {
	baseline, err := loader.Load(args[0])
	if err != nil {
		return fmt.Errorf("loading baseline: %w", err)
	}
	current, err := loader.Load(args[1])
	if err != nil {
		return fmt.Errorf("loading current: %w", err)
	}

	frozen := freezer.Freeze(baseline)
	violations, checkErr := freezer.Check(frozen, current, freezeFail)

	if freezeJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(violations); err != nil {
			return err
		}
	} else {
		if len(violations) == 0 {
			fmt.Println("✓ No freeze violations detected.")
		} else {
			for _, v := range violations {
				switch v.Kind {
				case "mutated":
					fmt.Printf("~ %s: %q → %q\n", v.Key, v.Frozen, v.Current)
				case "deleted":
					fmt.Printf("- %s (was %q)\n", v.Key, v.Frozen)
				case "added":
					fmt.Printf("+ %s = %q\n", v.Key, v.Current)
				}
			}
		}
	}

	return checkErr
}
