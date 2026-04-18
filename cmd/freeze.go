package cmd

import (
	"fmt"
	"os"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

var freezeCmd = &cobra.Command{
	Use:   "freeze <environment>",
	Short: "Freeze an environment to prevent patch application",
	Args:  cobra.ExactArgs(1),
	RunE:  runFreeze,
}

var unfreezeCmd = &cobra.Command{
	Use:   "unfreeze <environment>",
	Short: "Unfreeze a previously frozen environment",
	Args:  cobra.ExactArgs(1),
	RunE:  runUnfreeze,
}

func init() {
	rootCmd.AddCommand(freezeCmd)
	rootCmd.AddCommand(unfreezeCmd)
}

func runFreeze(cmd *cobra.Command, args []string) error {
	dir := workDir(cmd)
	st, err := state.Load(dir)
	if err != nil {
		return fmt.Errorf("load state: %w", err)
	}
	res, err := env.Freeze(st, args[0])
	if err != nil {
		return err
	}
	if err := state.Save(dir, st); err != nil {
		return fmt.Errorf("save state: %w", err)
	}
	if res.Frozen {
		fmt.Fprintf(os.Stdout, "environment %q is now frozen\n", res.Environment)
	} else {
		fmt.Fprintf(os.Stdout, "environment %q was already frozen\n", res.Environment)
	}
	return nil
}

func runUnfreeze(cmd *cobra.Command, args []string) error {
	dir := workDir(cmd)
	st, err := state.Load(dir)
	if err != nil {
		return fmt.Errorf("load state: %w", err)
	}
	res, err := env.Unfreeze(st, args[0])
	if err != nil {
		return err
	}
	if err := state.Save(dir, st); err != nil {
		return fmt.Errorf("save state: %w", err)
	}
	if !res.Frozen {
		fmt.Fprintf(os.Stdout, "environment %q has been unfrozen\n", res.Environment)
	} else {
		fmt.Fprintf(os.Stdout, "environment %q was already unfrozen\n", res.Environment)
	}
	return nil
}
