package cmd

import (
	"fmt"
	"os"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

func init() {
	badgeCmd := &cobra.Command{Use: "badge", Short: "Manage environment badges"}

	setCmd := &cobra.Command{
		Use:   "set <env> <key> <value>",
		Short: "Attach a badge to an environment",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBadgeSet(args[0], args[1], args[2])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list <env>",
		Short: "List all badges for an environment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBadgeList(args[0])
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <env> <key>",
		Short: "Remove a badge from an environment",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBadgeRemove(args[0], args[1])
		},
	}

	badgeCmd.AddCommand(setCmd, listCmd, removeCmd)
	rootCmd.AddCommand(badgeCmd)
}

func runBadgeSet(envName, key, value string) error {
	st, path, err := loadStateFromWorkDir()
	if err != nil {
		return err
	}
	if err := env.SetBadge(st, envName, key, value); err != nil {
		return err
	}
	return state.Save(st, path)
}

func runBadgeList(envName string) error {
	st, _, err := loadStateFromWorkDir()
	if err != nil {
		return err
	}
	badges := env.ListBadges(st, envName)
	if len(badges) == 0 {
		fmt.Fprintf(os.Stdout, "no badges set for %s\n", envName)
		return nil
	}
	for k, v := range badges {
		fmt.Fprintf(os.Stdout, "%s = %s\n", k, v)
	}
	return nil
}

func runBadgeRemove(envName, key string) error {
	st, path, err := loadStateFromWorkDir()
	if err != nil {
		return err
	}
	if err := env.RemoveBadge(st, envName, key); err != nil {
		return err
	}
	return state.Save(st, path)
}
