package cmd

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/patchwork-deploy/internal/config"
	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

func init() {
	spikeCmd := &cobra.Command{
		Use:   "spike",
		Short: "Manage spike (concurrency) limits for patches",
	}

	setCmd := &cobra.Command{
		Use:   "set <env> <patch> <limit>",
		Short: "Set a spike concurrency limit for a patch",
		Args:  cobra.ExactArgs(3),
		RunE:  runSpikeSet,
	}

	clearCmd := &cobra.Command{
		Use:   "clear <env> <patch>",
		Short: "Clear the spike limit for a patch",
		Args:  cobra.ExactArgs(2),
		RunE:  runSpikeClear,
	}

	listCmd := &cobra.Command{
		Use:   "list <env>",
		Short: "List all spike limits for an environment",
		Args:  cobra.ExactArgs(1),
		RunE:  runSpikeList,
	}

	spikeCmd.AddCommand(setCmd, clearCmd, listCmd)
	rootCmd.AddCommand(spikeCmd)
}

func runSpikeSet(cmd *cobra.Command, args []string) error {
	envName, patch, limitStr := args[0], args[1], args[2]
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return fmt.Errorf("invalid limit %q: must be an integer", limitStr)
	}
	cfg, err := config.Load(workDir())
	if err != nil {
		return err
	}
	st, err := state.Load(cfg.StateFile)
	if err != nil {
		return err
	}
	if err := env.SetSpike(st, envName, patch, limit); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "spike limit %d set for %s in %s\n", limit, patch, envName)
	return state.Save(cfg.StateFile, st)
}

func runSpikeClear(cmd *cobra.Command, args []string) error {
	envName, patch := args[0], args[1]
	cfg, err := config.Load(workDir())
	if err != nil {
		return err
	}
	st, err := state.Load(cfg.StateFile)
	if err != nil {
		return err
	}
	env.ClearSpike(st, envName, patch)
	fmt.Fprintf(cmd.OutOrStdout(), "spike limit cleared for %s in %s\n", patch, envName)
	return state.Save(cfg.StateFile, st)
}

func runSpikeList(cmd *cobra.Command, args []string) error {
	envName := args[0]
	cfg, err := config.Load(workDir())
	if err != nil {
		return err
	}
	st, err := state.Load(cfg.StateFile)
	if err != nil {
		return err
	}
	spikes := env.ListSpikes(st, envName)
	if len(spikes) == 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "no spike limits set for %s\n", envName)
		return nil
	}
	keys := make([]string, 0, len(spikes))
	for k := range spikes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(cmd.OutOrStdout(), "  %s: %d\n", k, spikes[k])
	}
	return nil
}
