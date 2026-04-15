package cmd

import (
	"fmt"
	"os"

	"github.com/patchwork-deploy/internal/config"
	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

var promoteCmd = &cobra.Command{
	Use:   "promote <from-env> <to-env>",
	Short: "Promote applied patches from one environment to another",
	Args:  cobra.ExactArgs(2),
	RunE:  runPromote,
}

func init() {
	promoteCmd.Flags().String("config", "patchwork.yaml", "path to patchwork config file")
	promoteCmd.Flags().String("state", ".patchwork-state.json", "path to state file")
	promoteCmd.Flags().Bool("dry-run", false, "preview promotion without persisting changes")
	rootCmd.AddCommand(promoteCmd)
}

func runPromote(cmd *cobra.Command, args []string) error {
	fromEnv := args[0]
	toEnv := args[1]

	cfgPath, _ := cmd.Flags().GetString("config")
	statePath, _ := cmd.Flags().GetString("state")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	_, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	st, err := state.Load(statePath)
	if err != nil {
		return fmt.Errorf("loading state: %w", err)
	}

	result, err := env.Promote(st, fromEnv, toEnv)
	if err != nil {
		return err
	}

	for _, p := range result.Promoted {
		fmt.Fprintf(os.Stdout, "  promoted: %s\n", p)
	}
	for _, p := range result.Skipped {
		fmt.Fprintf(os.Stdout, "  skipped:  %s (already applied)\n", p)
	}
	fmt.Fprintf(os.Stdout, "\npromoted %d patch(es) from %s → %s\n",
		len(result.Promoted), fromEnv, toEnv)

	if dryRun {
		fmt.Fprintln(os.Stdout, "(dry-run: state not saved)")
		return nil
	}

	if err := state.Save(st, statePath); err != nil {
		return fmt.Errorf("saving state: %w", err)
	}
	return nil
}
