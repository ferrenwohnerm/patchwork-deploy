package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
)

func init() {
	setCmd := &cobra.Command{
		Use:   "set <env> <patch> <limit>",
		Short: "Set an apply ceiling for a patch in an environment",
		Args:  cobra.ExactArgs(3),
		RunE:  runCeilingSet,
	}
	clearCmd := &cobra.Command{
		Use:   "clear <env> <patch>",
		Short: "Remove the ceiling for a patch",
		Args:  cobra.ExactArgs(2),
		RunE:  runCeilingClear,
	}
	listCmd := &cobra.Command{
		Use:   "list <env>",
		Short: "List all ceilings for an environment",
		Args:  cobra.ExactArgs(1),
		RunE:  runCeilingList,
	}
	ceilingCmd := &cobra.Command{
		Use:   "ceiling",
		Short: "Manage patch apply ceilings per environment",
	}
	ceilingCmd.AddCommand(setCmd, clearCmd, listCmd)
	rootCmd.AddCommand(ceilingCmd)
}

func runCeilingSet(cmd *cobra.Command, args []string) error {
	envName, patch, limitStr := args[0], args[1], args[2]
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return fmt.Errorf("limit must be an integer: %w", err)
	}
	dir := workDir(cmd)
	st, err := state.Load(dir)
	if err != nil {
		return err
	}
	if err := env.SetCeiling(st, envName, patch, limit); err != nil {
		return err
	}
	return state.Save(dir, st)
}

func runCeilingClear(cmd *cobra.Command, args []string) error {
	envName, patch := args[0], args[1]
	dir := workDir(cmd)
	st, err := state.Load(dir)
	if err != nil {
		return err
	}
	if err := env.ClearCeiling(st, envName, patch); err != nil {
		return err
	}
	return state.Save(dir, st)
}

func runCeilingList(cmd *cobra.Command, args []string) error {
	envName := args[0]
	dir := workDir(cmd)
	st, err := state.Load(dir)
	if err != nil {
		return err
	}
	ceilings, err := env.ListCeilings(st, envName)
	if err != nil {
		return err
	}
	var results []env.CeilingResult
	for patch, limit := range ceilings {
		current := env.CheckCeiling(st, envName, patch)
		results = append(results, env.CeilingResult{
			Patch:    patch,
			Limit:    limit,
			Current:  current,
			Exceeded: current >= limit,
		})
	}
	fmt.Fprintln(os.Stdout, env.FormatCeilings(results))
	return nil
}
