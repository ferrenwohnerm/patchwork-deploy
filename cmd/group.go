package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/patchwork-deploy/internal/config"
	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

func init() {
	groupCmd := &cobra.Command{Use: "group", Short: "Manage environment groups"}

	addCmd := &cobra.Command{
		Use:   "add <group> <env>",
		Short: "Add an environment to a group",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGroupAdd(args[0], args[1])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list [group]",
		Short: "List groups or members of a group",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGroupList(args)
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <group> <env>",
		Short: "Remove an environment from a group",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGroupRemove(args[0], args[1])
		},
	}

	groupCmd.AddCommand(addCmd, listCmd, removeCmd)
	rootCmd.AddCommand(groupCmd)
}

func runGroupAdd(group, envName string) error {
	cfg, err := config.Load("patchwork.yaml")
	if err != nil {
		return err
	}
	st, err := state.Load(cfg.StateFile)
	if err != nil {
		return err
	}
	if err := env.AddToGroup(st, group, envName); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "added %q to group %q\n", envName, group)
	return state.Save(cfg.StateFile, st)
}

func runGroupList(args []string) error {
	cfg, err := config.Load("patchwork.yaml")
	if err != nil {
		return err
	}
	st, err := state.Load(cfg.StateFile)
	if err != nil {
		return err
	}
	if len(args) == 0 {
		groups := env.ListAllGroups(st)
		if len(groups) == 0 {
			fmt.Println("no groups defined")
			return nil
		}
		fmt.Println(strings.Join(groups, "\n"))
		return nil
	}
	members, err := env.ListGroup(st, args[0])
	if err != nil {
		return err
	}
	fmt.Println(strings.Join(members, "\n"))
	return nil
}

func runGroupRemove(group, envName string) error {
	cfg, err := config.Load("patchwork.yaml")
	if err != nil {
		return err
	}
	st, err := state.Load(cfg.StateFile)
	if err != nil {
		return err
	}
	if err := env.RemoveFromGroup(st, group, envName); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "removed %q from group %q\n", envName, group)
	return state.Save(cfg.StateFile, st)
}
