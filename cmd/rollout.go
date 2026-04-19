package cmd

import (
	"fmt"
	"os"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

var rolloutCmd = &cobra.Command{Use: "rollout", Short: "Manage rollout strategies for patches"}

func init() {
	var stateFile string

	setCmd := &cobra.Command{
		Use:   "set <env> <patch> <strategy>",
		Short: "Set a rollout strategy (canary, blue-green, immediate)",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRolloutSet(stateFile, args[0], args[1], env.RolloutStrategy(args[2]))
		},
	}
	setCmd.Flags().StringVar(&stateFile, "state", "patchwork.state.json", "state file path")

	listCmd := &cobra.Command{
		Use:   "list <env>",
		Short: "List rollout strategies for an environment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRolloutList(stateFile, args[0])
		},
	}
	listCmd.Flags().StringVar(&stateFile, "state", "patchwork.state.json", "state file path")

	clearCmd := &cobra.Command{
		Use:   "clear <env> <patch>",
		Short: "Clear rollout strategy for a patch",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRolloutClear(stateFile, args[0], args[1])
		},
	}
	clearCmd.Flags().StringVar(&stateFile, "state", "patchwork.state.json", "state file path")

	rolloutCmd.AddCommand(setCmd, listCmd, clearCmd)
	rootCmd.AddCommand(rolloutCmd)
}

func runRolloutSet(stateFile, environment, patch string, strategy env.RolloutStrategy) error {
	if !env.IsValidStrategy(strategy) {
		return fmt.Errorf("invalid rollout strategy %q: must be one of canary, blue-green, immediate", strategy)
	}
	st, err := state.Load(stateFile)
	if err != nil {
		return err
	}
	if err := env.SetRollout(st, environment, patch, strategy); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "rollout strategy %q set for %s in %s\n", strategy, patch, environment)
	return state.Save(stateFile, st)
}

func runRolloutList(stateFile, environment string) error {
	st, err := state.Load(stateFile)
	if err != nil {
		return err
	}
	list := env.ListRollouts(st, environment)
	fmt.Fprintln(os.Stdout, env.FormatRollouts(environment, list))
	return nil
}

func runRolloutClear(stateFile, environment, patch string) error {
	st, err := state.Load(stateFile)
	if err != nil {
		return err
	}
	env.ClearRollout(st, environment, patch)
	fmt.Fprintf(os.Stdout, "rollout strategy cleared for %s in %s\n", patch, environment)
	return state.Save(stateFile, st)
}
