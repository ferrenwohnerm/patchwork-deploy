package cmd

import (
	"fmt"
	"os"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

var chainCmd = &cobra.Command{Use: "chain", Short: "Manage patch chains within an environment"}

var chainSetCmd = &cobra.Command{
	Use:   "set <env> <patch> <successor>",
	Short: "Set a successor patch that must follow the given patch",
	Args:  cobra.ExactArgs(3),
	RunE:  runChainSet,
}

var chainRemoveCmd = &cobra.Command{
	Use:   "remove <env> <patch>",
	Short: "Remove a chain entry for a patch",
	Args:  cobra.ExactArgs(2),
	RunE:  runChainRemove,
}

var chainListCmd = &cobra.Command{
	Use:   "list <env>",
	Short: "List all chain entries for an environment",
	Args:  cobra.ExactArgs(1),
	RunE:  runChainList,
}

func init() {
	chainCmd.AddCommand(chainSetCmd, chainRemoveCmd, chainListCmd)
	rootCmd.AddCommand(chainCmd)
}

func runChainSet(cmd *cobra.Command, args []string) error {
	envName, patch, successor := args[0], args[1], args[2]
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	if err := env.SetChain(st, envName, patch, successor); err != nil {
		return err
	}
	if err := state.Save(st, workDir()); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "chain set: %s -> %s in %s\n", patch, successor, envName)
	return nil
}

func runChainRemove(cmd *cobra.Command, args []string) error {
	envName, patch := args[0], args[1]
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	if err := env.RemoveChain(st, envName, patch); err != nil {
		return err
	}
	return state.Save(st, workDir())
}

func runChainList(cmd *cobra.Command, args []string) error {
	envName := args[0]
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	chains := env.ListChains(st, envName)
	fmt.Println(env.FormatChains(chains))
	return nil
}
