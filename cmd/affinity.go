package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

func init() {
	affinityCmd := &cobra.Command{
		Use:   "affinity",
		Short: "Manage patch affinity groups within an environment",
	}

	setCmd := &cobra.Command{
		Use:   "set <env> <patch> <group>",
		Short: "Assign an affinity group to a patch",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAffinitySet(args[0], args[1], args[2])
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <env> <patch>",
		Short: "Remove the affinity group from a patch",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAffinityRemove(args[0], args[1])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list <env>",
		Short: "List all affinity assignments for an environment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAffinityList(args[0])
		},
	}

	affinityCmd.AddCommand(setCmd, removeCmd, listCmd)
	rootCmd.AddCommand(affinityCmd)
}

func runAffinitySet(envName, patch, group string) error {
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	if err := env.SetAffinity(st, envName, patch, group); err != nil {
		return err
	}
	fmt.Printf("affinity %q assigned to %s/%s\n", group, envName, patch)
	return state.Save(workDir(), st)
}

func runAffinityRemove(envName, patch string) error {
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	if err := env.RemoveAffinity(st, envName, patch); err != nil {
		return err
	}
	fmt.Printf("affinity cleared for %s/%s\n", envName, patch)
	return state.Save(workDir(), st)
}

func runAffinityList(envName string) error {
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	affinities := env.ListAffinities(st, envName)
	if len(affinities) == 0 {
		fmt.Printf("no affinities set for environment %q\n", envName)
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PATCH\tGROUP")
	for patch, group := range affinities {
		fmt.Fprintf(w, "%s\t%s\n", patch, group)
	}
	return w.Flush()
}
