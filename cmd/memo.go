package cmd

import (
	"fmt"
	"os"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

var memoCmd = &cobra.Command{
	Use:   "memo",
	Short: "Manage freeform memos attached to environments",
}

var memoSetCmd = &cobra.Command{
	Use:   "set <environment> <text>",
	Short: "Set a memo on an environment",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMemoSet(args[0], args[1])
	},
}

var memoGetCmd = &cobra.Command{
	Use:   "get <environment>",
	Short: "Get the memo for an environment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMemoGet(args[0])
	},
}

var memoClearCmd = &cobra.Command{
	Use:   "clear <environment>",
	Short: "Clear the memo for an environment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMemoClear(args[0])
	},
}

func init() {
	memoCmd.AddCommand(memoSetCmd, memoGetCmd, memoClearCmd)
	rootCmd.AddCommand(memoCmd)
}

func runMemoSet(envName, text string) error {
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	if err := env.SetMemo(st, envName, text); err != nil {
		return err
	}
	return state.Save(workDir(), st)
}

func runMemoGet(envName string) error {
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	val, err := env.GetMemo(st, envName)
	if err != nil {
		return err
	}
	if val == "" {
		fmt.Fprintln(os.Stdout, "(no memo set)")
	} else {
		fmt.Fprintln(os.Stdout, val)
	}
	return nil
}

func runMemoClear(envName string) error {
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	if err := env.ClearMemo(st, envName); err != nil {
		return err
	}
	return state.Save(workDir(), st)
}
