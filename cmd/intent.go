package cmd

import (
	"fmt"
	"os"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

func init() {
	intentCmd := &cobra.Command{Use: "intent", Short: "Manage deployment intents for patches"}

	setCmd := &cobra.Command{
		Use:   "set <env> <patch> <intent>",
		Short: "Set a deployment intent for a patch",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runIntentSet(args[0], args[1], args[2])
		},
	}

	getCmd := &cobra.Command{
		Use:   "get <env> <patch>",
		Short: "Get the deployment intent for a patch",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runIntentGet(args[0], args[1])
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <env> <patch>",
		Short: "Remove the deployment intent for a patch",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runIntentRemove(args[0], args[1])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list <env>",
		Short: "List all deployment intents for an environment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runIntentList(args[0])
		},
	}

	intentCmd.AddCommand(setCmd, getCmd, removeCmd, listCmd)
	rootCmd.AddCommand(intentCmd)
}

func runIntentSet(envName, patch, intent string) error {
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	if err := env.SetIntent(st, envName, patch, intent); err != nil {
		return err
	}
	return state.Save(st, workDir())
}

func runIntentGet(envName, patch string) error {
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	v, ok := env.GetIntent(st, envName, patch)
	if !ok {
		fmt.Fprintf(os.Stderr, "no intent set for %s in %s\n", patch, envName)
		return nil
	}
	fmt.Println(v)
	return nil
}

func runIntentRemove(envName, patch string) error {
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	if err := env.RemoveIntent(st, envName, patch); err != nil {
		return err
	}
	return state.Save(st, workDir())
}

func runIntentList(envName string) error {
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	intents, err := env.ListIntents(st, envName)
	if err != nil {
		return err
	}
	if len(intents) == 0 {
		fmt.Println("no intents set")
		return nil
	}
	for patch, intent := range intents {
		fmt.Printf("%-30s %s\n", patch, intent)
	}
	return nil
}
