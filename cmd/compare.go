package cmd

import (
	"fmt"
	"os"

	"github.com/patchwork-deploy/internal/config"
	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

var compareCmd = &cobra.Command{
	Use:   "compare <source> <target>",
	Short: "Compare applied patches between two environments",
	Args:  cobra.ExactArgs(2),
	RunE:  runCompare,
}

func init() {
	compareCmd.Flags().String("config", "patchwork.yaml", "path to config file")
	compareCmd.Flags().String("state", ".patchwork-state.json", "path to state file")
	rootCmd.AddCommand(compareCmd)
}

func runCompare(cmd *cobra.Command, args []string) error {
	cfgPath, _ := cmd.Flags().GetString("config")
	statePath, _ := cmd.Flags().GetString("state")
	source := args[0]
	target := args[1]

	_, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	st, err := state.Load(statePath)
	if err != nil {
		return fmt.Errorf("load state: %w", err)
	}

	res, err := env.CompareEnvironments(st, source, target)
	if err != nil {
		return err
	}

	w := os.Stdout
	fmt.Fprintf(w, "Comparing %s → %s\n", res.Source, res.Target)
	fmt.Fprintf(w, "  In both       : %d\n", len(res.InBoth))
	for _, p := range res.InBoth {
		fmt.Fprintf(w, "    = %s\n", p)
	}
	fmt.Fprintf(w, "  Only in %-8s: %d\n", source, len(res.OnlyInSource))
	for _, p := range res.OnlyInSource {
		fmt.Fprintf(w, "    < %s\n", p)
	}
	fmt.Fprintf(w, "  Only in %-8s: %d\n", target, len(res.OnlyInTarget))
	for _, p := range res.OnlyInTarget {
		fmt.Fprintf(w, "    > %s\n", p)
	}
	return nil
}
