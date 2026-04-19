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
	capCmd := &cobra.Command{Use: "capability", Short: "Manage environment capabilities"}

	addCmd := &cobra.Command{
		Use:   "add <env> <capability>",
		Short: "Add a capability to an environment",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCapabilityAdd(args[0], args[1])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list <env>",
		Short: "List capabilities for an environment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCapabilityList(args[0])
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <env> <capability>",
		Short: "Remove a capability from an environment",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCapabilityRemove(args[0], args[1])
		},
	}

	capCmd.AddCommand(addCmd, listCmd, removeCmd)
	rootCmd.AddCommand(capCmd)
}

func runCapabilityAdd(envName, cap string) error {
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	if err := env.AddCapability(st, envName, cap); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "capability %q added to %q\n", cap, envName)
	return state.Save(workDir(), st)
}

func runCapabilityList(envName string) error {
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	caps := env.ListCapabilities(st, envName)
	if len(caps) == 0 {
		fmt.Fprintln(os.Stdout, "no capabilities set")
		return nil
	}
	fmt.Fprintln(os.Stdout, strings.Join(caps, "\n"))
	return nil
}

func runCapabilityRemove(envName, cap string) error {
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	if err := env.RemoveCapability(st, envName, cap); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "capability %q removed from %q\n", cap, envName)
	return state.Save(workDir(), st)
}
