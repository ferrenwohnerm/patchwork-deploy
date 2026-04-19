package cmd

import (
	"fmt"
	"os"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

var suppressCmd = &cobra.Command{
	Use:   "suppress",
	Short: "Manage patch suppression within an environment",
}

var suppressAddCmd = &cobra.Command{
	Use:   "add <env> <patch> [reason]",
	Short: "Suppress a patch so it is skipped during runs",
	Args:  cobra.RangeArgs(2, 3),
	RunE:  runSuppressAdd,
}

var suppressRemoveCmd = &cobra.Command{
	Use:   "remove <env> <patch>",
	Short: "Remove suppression from a patch",
	Args:  cobra.ExactArgs(2),
	RunE:  runSuppressRemove,
}

var suppressListCmd = &cobra.Command{
	Use:   "list <env>",
	Short: "List all suppressed patches in an environment",
	Args:  cobra.ExactArgs(1),
	RunE:  runSuppressList,
}

func init() {
	suppressCmd.AddCommand(suppressAddCmd, suppressRemoveCmd, suppressListCmd)
	rootCmd.AddCommand(suppressCmd)
}

func runSuppressAdd(cmd *cobra.Command, args []string) error {
	envName, patch := args[0], args[1]
	reason := ""
	if len(args) == 3 {
		reason = args[2]
	}
	st, err := state.Load(workDir(cmd))
	if err != nil {
		return err
	}
	if err := env.Suppress(st, envName, patch, reason); err != nil {
		return err
	}
	return state.Save(workDir(cmd), st)
}

func runSuppressRemove(cmd *cobra.Command, args []string) error {
	st, err := state.Load(workDir(cmd))
	if err != nil {
		return err
	}
	if err := env.Unsuppress(st, args[0], args[1]); err != nil {
		return err
	}
	return state.Save(workDir(cmd), st)
}

func runSuppressList(cmd *cobra.Command, args []string) error {
	st, err := state.Load(workDir(cmd))
	if err != nil {
		return err
	}
	list := env.ListSuppressed(st, args[0])
	if len(list) == 0 {
		fmt.Fprintln(os.Stdout, "no suppressed patches")
		return nil
	}
	for patch, reason := range list {
		if reason == "" {
			fmt.Fprintf(os.Stdout, "  %s\n", patch)
		} else {
			fmt.Fprintf(os.Stdout, "  %s  (%s)\n", patch, reason)
		}
	}
	return nil
}
