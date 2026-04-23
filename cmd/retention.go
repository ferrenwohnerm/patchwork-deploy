package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/your-org/patchwork-deploy/internal/env"
	"github.com/your-org/patchwork-deploy/internal/state"
)

func init() {
	retentionCmd := &cobra.Command{
		Use:   "retention",
		Short: "Manage patch record retention limits per environment",
	}

	setCmd := &cobra.Command{
		Use:   "set <env> <limit>",
		Short: "Set the maximum number of patch records to retain",
		Args:  cobra.ExactArgs(2),
		RunE:  runRetentionSet,
	}

	clearCmd := &cobra.Command{
		Use:   "clear <env>",
		Short: "Remove the retention limit for an environment",
		Args:  cobra.ExactArgs(1),
		RunE:  runRetentionClear,
	}

	getCmd := &cobra.Command{
		Use:   "get <env>",
		Short: "Show the retention limit for an environment",
		Args:  cobra.ExactArgs(1),
		RunE:  runRetentionGet,
	}

	retentionCmd.AddCommand(setCmd, clearCmd, getCmd)
	rootCmd.AddCommand(retentionCmd)
}

func runRetentionSet(cmd *cobra.Command, args []string) error {
	envName := args[0]
	limit, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("limit must be an integer: %w", err)
	}
	dir := workDir(cmd)
	st, err := state.LoadFromDir(dir)
	if err != nil {
		return err
	}
	if err := env.SetRetention(st, envName, limit); err != nil {
		return err
	}
	if err := state.SaveToDir(dir, st); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "retention limit for %q set to %d\n", envName, limit)
	return nil
}

func runRetentionClear(cmd *cobra.Command, args []string) error {
	envName := args[0]
	dir := workDir(cmd)
	st, err := state.LoadFromDir(dir)
	if err != nil {
		return err
	}
	if err := env.ClearRetention(st, envName); err != nil {
		return err
	}
	if err := state.SaveToDir(dir, st); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "retention limit cleared for %q\n", envName)
	return nil
}

func runRetentionGet(cmd *cobra.Command, args []string) error {
	envName := args[0]
	dir := workDir(cmd)
	st, err := state.LoadFromDir(dir)
	if err != nil {
		return err
	}
	n, ok := env.GetRetention(st, envName)
	if !ok {
		fmt.Fprintf(os.Stdout, "no retention limit set for %q\n", envName)
		return nil
	}
	fmt.Fprintf(os.Stdout, "retention limit for %q: %d\n", envName, n)
	return nil
}
