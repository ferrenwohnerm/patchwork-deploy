package cmd

import (
	"fmt"
	"os"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

var baselineCmd = &cobra.Command{
	Use:   "baseline",
	Short: "Manage environment baselines",
}

var baselineSetCmd = &cobra.Command{
	Use:   "set <env> <patch>",
	Short: "Set a baseline patch for an environment",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runBaselineSet(args[0], args[1])
	},
}

var baselineGetCmd = &cobra.Command{
	Use:   "get <env>",
	Short: "Show the current baseline for an environment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runBaselineGet(args[0])
	},
}

var baselineClearCmd = &cobra.Command{
	Use:   "clear <env>",
	Short: "Clear the baseline for an environment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runBaselineClear(args[0])
	},
}

func init() {
	baselineCmd.AddCommand(baselineSetCmd, baselineGetCmd, baselineClearCmd)
	rootCmd.AddCommand(baselineCmd)
}

func runBaselineSet(environment, patch string) error {
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	if err := env.SetBaseline(st, environment, patch); err != nil {
		return err
	}
	if err := st.Save(workDir()); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "baseline set to %q for environment %q\n", patch, environment)
	return nil
}

func runBaselineGet(environment string) error {
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	b, ok := env.GetBaseline(st, environment)
	if !ok {
		fmt.Fprintf(os.Stdout, "no baseline set for environment %q\n", environment)
		return nil
	}
	fmt.Fprintf(os.Stdout, "baseline: %s (set at %s)\n", b.Patch, b.CreatedAt.Format("2006-01-02 15:04:05"))
	return nil
}

func runBaselineClear(environment string) error {
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	if err := env.ClearBaseline(st, environment); err != nil {
		return err
	}
	if err := st.Save(workDir()); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "baseline cleared for environment %q\n", environment)
	return nil
}
