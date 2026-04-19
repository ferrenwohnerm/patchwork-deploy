package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

func init() {
	priorityCmd := &cobra.Command{
		Use:   "priority",
		Short: "Manage patch priorities within an environment",
	}

	setCmd := &cobra.Command{
		Use:   "set <env> <patch> <priority>",
		Short: "Set a numeric priority for a patch",
		Args:  cobra.ExactArgs(3),
		RunE:  runPrioritySet,
	}

	listCmd := &cobra.Command{
		Use:   "list <env>",
		Short: "List all patch priorities for an environment",
		Args:  cobra.ExactArgs(1),
		RunE:  runPriorityList,
	}

	clearCmd := &cobra.Command{
		Use:   "clear <env> <patch>",
		Short: "Clear the priority for a patch",
		Args:  cobra.ExactArgs(2),
		RunE:  runPriorityClear,
	}

	priorityCmd.AddCommand(setCmd, listCmd, clearCmd)
	rootCmd.AddCommand(priorityCmd)
}

func runPrioritySet(cmd *cobra.Command, args []string) error {
	envName, patch, rawN := args[0], args[1], args[2]
	n, err := strconv.Atoi(rawN)
	if err != nil {
		return fmt.Errorf("priority must be an integer: %w", err)
	}
	dir := workDir(cmd)
	st, err := state.Load(dir)
	if err != nil {
		return err
	}
	if err := env.SetPriority(st, envName, patch, n); err != nil {
		return err
	}
	return state.Save(dir, st)
}

func runPriorityList(cmd *cobra.Command, args []string) error {
	dir := workDir(cmd)
	st, err := state.Load(dir)
	if err != nil {
		return err
	}
	m := env.ListPriorities(st, args[0])
	if len(m) == 0 {
		fmt.Fprintln(os.Stdout, "no priorities set")
		return nil
	}
	for patch, p := range m {
		fmt.Fprintf(os.Stdout, "%-30s %d\n", patch, p)
	}
	return nil
}

func runPriorityClear(cmd *cobra.Command, args []string) error {
	dir := workDir(cmd)
	st, err := state.Load(dir)
	if err != nil {
		return err
	}
	env.ClearPriority(st, args[0], args[1])
	return state.Save(dir, st)
}
