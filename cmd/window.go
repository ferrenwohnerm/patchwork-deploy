package cmd

import (
	"fmt"
	"os"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

var windowCmd = &cobra.Command{
	Use:   "window",
	Short: "Manage deployment windows for patches",
}

var windowSetCmd = &cobra.Command{
	Use:   "set <env> <patch> <start> <end>",
	Short: "Set a deployment window (HH:MM format)",
	Args:  cobra.ExactArgs(4),
	RunE:  runWindowSet,
}

var windowClearCmd = &cobra.Command{
	Use:   "clear <env> <patch>",
	Short: "Clear the deployment window for a patch",
	Args:  cobra.ExactArgs(2),
	RunE:  runWindowClear,
}

var windowListCmd = &cobra.Command{
	Use:   "list <env>",
	Short: "List all deployment windows for an environment",
	Args:  cobra.ExactArgs(1),
	RunE:  runWindowList,
}

func init() {
	windowCmd.AddCommand(windowSetCmd, windowClearCmd, windowListCmd)
	rootCmd.AddCommand(windowCmd)
}

func runWindowSet(cmd *cobra.Command, args []string) error {
	environment, patch, start, end := args[0], args[1], args[2], args[3]
	st, path, err := loadStateForCmd(cmd)
	if err != nil {
		return err
	}
	if err := env.SetWindow(st, environment, patch, start, end); err != nil {
		return err
	}
	if err := state.Save(st, path); err != nil {
		return fmt.Errorf("saving state: %w", err)
	}
	fmt.Fprintf(os.Stdout, "window set for %s in %s: %s - %s\n", patch, environment, start, end)
	return nil
}

func runWindowClear(cmd *cobra.Command, args []string) error {
	environment, patch := args[0], args[1]
	st, path, err := loadStateForCmd(cmd)
	if err != nil {
		return err
	}
	env.ClearWindow(st, environment, patch)
	if err := state.Save(st, path); err != nil {
		return fmt.Errorf("saving state: %w", err)
	}
	fmt.Fprintf(os.Stdout, "window cleared for %s in %s\n", patch, environment)
	return nil
}

func runWindowList(cmd *cobra.Command, args []string) error {
	environment := args[0]
	st, _, err := loadStateForCmd(cmd)
	if err != nil {
		return err
	}
	wins := env.ListWindows(st, environment)
	fmt.Println(env.FormatWindows(wins))
	return nil
}
