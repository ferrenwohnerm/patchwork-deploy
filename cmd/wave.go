package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

func init() {
	waveCmd := &cobra.Command{
		Use:   "wave",
		Short: "Manage deployment wave assignments for patches",
	}

	setCmd := &cobra.Command{
		Use:   "set <env> <patch> <wave>",
		Short: "Assign a wave number to a patch",
		Args:  cobra.ExactArgs(3),
		RunE:  runWaveSet,
	}

	clearCmd := &cobra.Command{
		Use:   "clear <env> <patch>",
		Short: "Remove wave assignment from a patch",
		Args:  cobra.ExactArgs(2),
		RunE:  runWaveClear,
	}

	listCmd := &cobra.Command{
		Use:   "list <env>",
		Short: "List all wave assignments for an environment",
		Args:  cobra.ExactArgs(1),
		RunE:  runWaveList,
	}

	waveCmd.AddCommand(setCmd, clearCmd, listCmd)
	rootCmd.AddCommand(waveCmd)
}

func runWaveSet(cmd *cobra.Command, args []string) error {
	envName, patch, waveStr := args[0], args[1], args[2]
	waveNum, err := strconv.Atoi(waveStr)
	if err != nil {
		return fmt.Errorf("wave must be an integer: %w", err)
	}
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	if err := env.SetWave(st, envName, patch, waveNum); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "wave %d assigned to %s in %s\n", waveNum, patch, envName)
	return state.Save(st, workDir())
}

func runWaveClear(cmd *cobra.Command, args []string) error {
	envName, patch := args[0], args[1]
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	env.ClearWave(st, envName, patch)
	fmt.Fprintf(cmd.OutOrStdout(), "wave cleared for %s in %s\n", patch, envName)
	return state.Save(st, workDir())
}

func runWaveList(cmd *cobra.Command, args []string) error {
	envName := args[0]
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	entries := env.ListWaves(st, envName)
	if len(entries) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no wave assignments")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PATCH\tWAVE")
	for _, e := range entries {
		fmt.Fprintf(w, "%s\t%d\n", e.Patch, e.Wave)
	}
	return w.Flush()
}
