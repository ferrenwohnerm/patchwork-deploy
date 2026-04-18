package cmd

import (
	"fmt"
	"strings"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

func init() {
	depsCmd := &cobra.Command{Use: "deps", Short: "Manage patch dependencies"}

	addCmd := &cobra.Command{
		Use:   "add <env> <patch> <depends-on>",
		Short: "Record that a patch depends on another",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDepsAdd(args[0], args[1], args[2])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list <env> <patch>",
		Short: "List dependencies of a patch",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDepsList(args[0], args[1])
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <env> <patch> <depends-on>",
		Short: "Remove a dependency edge",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDepsRemove(args[0], args[1], args[2])
		},
	}

	depsCmd.AddCommand(addCmd, listCmd, removeCmd)
	rootCmd.AddCommand(depsCmd)
}

func runDepsAdd(environment, patch, dependsOn string) error {
	if patch == dependsOn {
		return fmt.Errorf("patch %q cannot depend on itself", patch)
	}
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	if err := env.SetDep(st, environment, patch, dependsOn); err != nil {
		return err
	}
	if err := state.Save(st, workDir()); err != nil {
		return err
	}
	fmt.Printf("dependency recorded: %s -> %s (env: %s)\n", patch, dependsOn, environment)
	return nil
}

func runDepsList(environment, patch string) error {
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	deps, err := env.GetDeps(st, environment, patch)
	if err != nil {
		return err
	}
	if len(deps) == 0 {
		fmt.Printf("%s has no recorded dependencies\n", patch)
		return nil
	}
	fmt.Printf("dependencies of %s:\n  %s\n", patch, strings.Join(deps, "\n  "))
	return nil
}

func runDepsRemove(environment, patch, dependsOn string) error {
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	if err := env.RemoveDep(st, environment, patch, dependsOn); err != nil {
		return err
	}
	return state.Save(st, workDir())
}
