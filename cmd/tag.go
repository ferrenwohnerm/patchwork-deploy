package cmd

import (
	"fmt"
	"os"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{Use: "tag", Short: "Manage patch tags"}

var tagAddCmd = &cobra.Command{
	Use:   "add <env> <patch> <tag>",
	Short: "Add a tag to a patch in an environment",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTagAdd(args[0], args[1], args[2])
	},
}

var tagListCmd = &cobra.Command{
	Use:   "list <env>",
	Short: "List tags for an environment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTagList(args[0])
	},
}

var tagRemoveCmd = &cobra.Command{
	Use:   "remove <env> <patch>",
	Short: "Remove a tag from a patch",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTagRemove(args[0], args[1])
	},
}

func init() {
	tagCmd.AddCommand(tagAddCmd, tagListCmd, tagRemoveCmd)
	rootCmd.AddCommand(tagCmd)
}

func runTagAdd(environment, patch, tag string) error {
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	if err := env.AddTag(st, environment, patch, tag); err != nil {
		return err
	}
	return state.Save(st, workDir())
}

func runTagList(environment string) error {
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	tags, err := env.ListTags(st, environment)
	if err != nil {
		return err
	}
	if len(tags) == 0 {
		fmt.Fprintln(os.Stdout, "no tags found")
		return nil
	}
	for _, t := range tags {
		fmt.Fprintf(os.Stdout, "%s\t%s\t%s\n", t.Environment, t.Patch, t.Tag)
	}
	return nil
}

func runTagRemove(environment, patch string) error {
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	if err := env.RemoveTag(st, environment, patch); err != nil {
		return err
	}
	return state.Save(st, workDir())
}
