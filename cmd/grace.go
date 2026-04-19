package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

var graceCmd = &cobra.Command{
	Use:   "grace",
	Short: "Manage grace periods for applied patches",
}

var graceSetCmd = &cobra.Command{
	Use:   "set <env> <patch> <duration>",
	Short: "Set a grace period for a patch (e.g. 30m, 2h)",
	Args:  cobra.ExactArgs(3),
	RunE:  runGraceSet,
}

var graceClearCmd = &cobra.Command{
	Use:   "clear <env> <patch>",
	Short: "Clear the grace period for a patch",
	Args:  cobra.ExactArgs(2),
	RunE:  runGraceClear,
}

var graceCheckCmd = &cobra.Command{
	Use:   "check <env> <patch>",
	Short: "Check whether a patch is within its grace period",
	Args:  cobra.ExactArgs(2),
	RunE:  runGraceCheck,
}

func init() {
	graceCmd.AddCommand(graceSetCmd, graceClearCmd, graceCheckCmd)
	rootCmd.AddCommand(graceCmd)
}

func runGraceSet(cmd *cobra.Command, args []string) error {
	envName, patch, rawDur := args[0], args[1], args[2]
	dur, err := time.ParseDuration(rawDur)
	if err != nil {
		return fmt.Errorf("invalid duration %q: %w", rawDur, err)
	}
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	if err := env.SetGracePeriod(st, envName, patch, dur); err != nil {
		return err
	}
	return state.Save(workDir(), st)
}

func runGraceClear(cmd *cobra.Command, args []string) error {
	envName, patch := args[0], args[1]
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	env.ClearGracePeriod(st, envName, patch)
	return state.Save(workDir(), st)
}

func runGraceCheck(cmd *cobra.Command, args []string) error {
	envName, patch := args[0], args[1]
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	in, err := env.InGracePeriod(st, envName, patch)
	if err != nil {
		return err
	}
	if in {
		fmt.Fprintf(os.Stdout, "patch %q in %q is within its grace period\n", patch, envName)
	} else {
		fmt.Fprintf(os.Stdout, "patch %q in %q is NOT within a grace period\n", patch, envName)
	}
	return nil
}
