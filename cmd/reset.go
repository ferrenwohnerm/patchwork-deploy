package cmd

import (
	"fmt"
	"os"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

var resetDryRun bool

func init() {
	resetCmd := &cobra.Command{
		Use:   "reset <environment>",
		Short: "Remove all applied patch records for an environment",
		Args:  cobra.ExactArgs(1),
		RunE:  runReset,
	}
	resetCmd.Flags().BoolVar(&resetDryRun, "dry-run", false, "Preview changes without modifying state")
	rootCmd.AddCommand(resetCmd)
}

func runReset(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	envName := args[0]
	wd := workDir()

	st, err := state.Load(wd + "/" + cfg.StateFile)
	if err != nil {
		return fmt.Errorf("loading state: %w", err)
	}

	result, err := env.Reset(st, envName, resetDryRun)
	if err != nil {
		return err
	}

	if resetDryRun {
		fmt.Fprintf(os.Stdout, "[dry-run] would remove %d record(s) for environment %q\n", result.Removed, result.Environment)
		return nil
	}

	if err := state.Save(st, wd+"/"+cfg.StateFile); err != nil {
		return fmt.Errorf("saving state: %w", err)
	}

	fmt.Fprintf(os.Stdout, "reset %d record(s) for environment %q\n", result.Removed, result.Environment)
	return nil
}
