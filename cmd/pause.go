package cmd

import (
	"fmt"
	"os"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

var pauseCmd = &cobra.Command{
	Use:   "pause <environment>",
	Short: "Pause patch application for an environment",
	Args:  cobra.ExactArgs(1),
	RunE:  runPause,
}

var unpauseCmd = &cobra.Command{
	Use:   "unpause <environment>",
	Short: "Resume patch application for an environment",
	Args:  cobra.ExactArgs(1),
	RunE:  runUnpause,
}

func init() {
	rootCmd.AddCommand(pauseCmd)
	rootCmd.AddCommand(unpauseCmd)
	pauseCmd.Flags().String("state", "patchwork.state.json", "path to state file")
	unpauseCmd.Flags().String("state", "patchwork.state.json", "path to state file")
}

func runPause(cmd *cobra.Command, args []string) error {
	statePath, _ := cmd.Flags().GetString("state")
	envName := args[0]

	st, err := state.Load(statePath)
	if err != nil {
		return fmt.Errorf("loading state: %w", err)
	}
	if err := env.Pause(st, envName); err != nil {
		return err
	}
	if err := st.Save(statePath); err != nil {
		return fmt.Errorf("saving state: %w", err)
	}
	fmt.Fprintf(os.Stdout, "environment %q paused\n", envName)
	return nil
}

func runUnpause(cmd *cobra.Command, args []string) error {
	statePath, _ := cmd.Flags().GetString("state")
	envName := args[0]

	st, err := state.Load(statePath)
	if err != nil {
		return fmt.Errorf("loading state: %w", err)
	}
	if err := env.Unpause(st, envName); err != nil {
		return err
	}
	if err := st.Save(statePath); err != nil {
		return fmt.Errorf("saving state: %w", err)
	}
	fmt.Fprintf(os.Stdout, "environment %q unpaused\n", envName)
	return nil
}
