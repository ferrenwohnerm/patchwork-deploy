package cmd

import (
	"fmt"
	"os"

	"github.com/patchwork-deploy/internal/config"
	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

var cloneCmd = &cobra.Command{
	Use:   "clone <source> <target>",
	Short: "Clone all patch state from one environment into a new target environment",
	Args:  cobra.ExactArgs(2),
	RunE:  runClone,
}

func init() {
	cloneCmd.Flags().String("config", "patchwork.yaml", "path to patchwork config file")
	cloneCmd.Flags().String("state", ".patchwork-state.json", "path to state file")
	rootCmd.AddCommand(cloneCmd)
}

func runClone(cmd *cobra.Command, args []string) error {
	source := args[0]
	target := args[1]

	cfgPath, _ := cmd.Flags().GetString("config")
	statePath, _ := cmd.Flags().GetString("state")

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	st, err := state.Load(statePath)
	if err != nil {
		return fmt.Errorf("loading state: %w", err)
	}

	result, err := env.Clone(cfg, st, source, target)
	if err != nil {
		return fmt.Errorf("clone failed: %w", err)
	}

	if err := state.Save(statePath, st); err != nil {
		return fmt.Errorf("saving state: %w", err)
	}

	fmt.Fprintf(os.Stdout, "cloned %d patch(es) from %q to %q\n",
		result.Copied, result.SourceEnv, result.TargetEnv)
	return nil
}
