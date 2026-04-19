package cmd

import (
	"fmt"
	"os"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

func init() {
	triggerCmd := &cobra.Command{Use: "trigger", Short: "Manage patch triggers"}

	addCmd := &cobra.Command{
		Use:   "add <env> <patch> <event> <action>",
		Short: "Add a trigger for a patch event",
		Args:  cobra.ExactArgs(4),
		RunE:  runTriggerAdd,
	}

	listCmd := &cobra.Command{
		Use:   "list <env>",
		Short: "List triggers for an environment",
		Args:  cobra.ExactArgs(1),
		RunE:  runTriggerList,
	}

	removeCmd := &cobra.Command{
		Use:   "remove <env> <patch> <event>",
		Short: "Remove a trigger",
		Args:  cobra.ExactArgs(3),
		RunE:  runTriggerRemove,
	}

	triggerCmd.AddCommand(addCmd, listCmd, removeCmd)
	rootCmd.AddCommand(triggerCmd)
}

func runTriggerAdd(cmd *cobra.Command, args []string) error {
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	if err := env.SetTrigger(st, args[0], args[1], args[2], args[3]); err != nil {
		return err
	}
	return state.Save(workDir(), st)
}

func runTriggerList(cmd *cobra.Command, args []string) error {
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	triggers, err := env.ListTriggers(st, args[0])
	if err != nil {
		return err
	}
	if len(triggers) == 0 {
		fmt.Fprintln(os.Stdout, "no triggers defined")
		return nil
	}
	for _, t := range triggers {
		fmt.Fprintf(os.Stdout, "patch=%-20s event=%-20s action=%s\n", t.Patch, t.Event, t.Action)
	}
	return nil
}

func runTriggerRemove(cmd *cobra.Command, args []string) error {
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	if err := env.RemoveTrigger(st, args[0], args[1], args[2]); err != nil {
		return err
	}
	return state.Save(workDir(), st)
}
