package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

func init() {
	drainCmd := &cobra.Command{
		Use:   "drain <environment>",
		Short: "Drain an environment by clearing its scheduled patches",
		Args:  cobra.ExactArgs(1),
		RunE:  runDrain,
	}
	drainCmd.Flags().Bool("dry-run", false, "Preview what would be removed without making changes")
	drainCmd.Flags().Bool("undo", false, "Remove the drain sentinel from the environment")
	drainCmd.Flags().String("state", "patchwork.state.json", "Path to state file")
	RootCmd.AddCommand(drainCmd)
}

func runDrain(cmd *cobra.Command, args []string) error {
	environment := args[0]
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	undo, _ := cmd.Flags().GetBool("undo")
	statePath, _ := cmd.Flags().GetString("state")

	st, err := state.Load(statePath)
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	if undo {
		if err := env.Undrain(st, environment); err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "drain sentinel removed for %q\n", environment)
		return state.Save(statePath, st)
	}

	result, err := env.Drain(st, environment, dryRun)
	if err != nil {
		return err
	}

	if len(result.Removed) == 0 {
		fmt.Fprintf(os.Stdout, "no scheduled patches to drain for %q\n", environment)
		return nil
	}

	prefix := ""
	if dryRun {
		prefix = "[dry-run] "
	}
	fmt.Fprintf(os.Stdout, "%sdrained %d scheduled patch(es) from %q:\n", prefix, len(result.Removed), environment)
	for _, p := range result.Removed {
		fmt.Fprintf(os.Stdout, "  - %s\n", strings.TrimSpace(p))
	}

	if dryRun {
		return nil
	}
	return state.Save(statePath, st)
}
