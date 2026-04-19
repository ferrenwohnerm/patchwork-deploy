package cmd

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

func init() {
	setCmd := &cobra.Command{
		Use:   "set <env> <patch> <weight>",
		Short: "Set a numeric weight for a patch in an environment",
		Args:  cobra.ExactArgs(3),
		RunE:  runWeightSet,
	}
	listCmd := &cobra.Command{
		Use:   "list <env>",
		Short: "List all patch weights for an environment",
		Args:  cobra.ExactArgs(1),
		RunE:  runWeightList,
	}
	clearCmd := &cobra.Command{
		Use:   "clear <env> <patch>",
		Short: "Clear the weight for a patch in an environment",
		Args:  cobra.ExactArgs(2),
		RunE:  runWeightClear,
	}
	weightCmd := &cobra.Command{
		Use:   "weight",
		Short: "Manage patch weights within environments",
	}
	weightCmd.AddCommand(setCmd, listCmd, clearCmd)
	rootCmd.AddCommand(weightCmd)
}

func runWeightSet(cmd *cobra.Command, args []string) error {
	envName, patch, rawWeight := args[0], args[1], args[2]
	w, err := strconv.Atoi(rawWeight)
	if err != nil {
		return fmt.Errorf("weight must be an integer: %v", err)
	}
	dir := workDir(cmd)
	st, err := state.Load(dir)
	if err != nil {
		return err
	}
	if err := env.SetWeight(st, envName, patch, w); err != nil {
		return err
	}
	if err := state.Save(st, dir); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "weight %d set for patch %q in %q\n", w, patch, envName)
	return nil
}

func runWeightList(cmd *cobra.Command, args []string) error {
	envName := args[0]
	dir := workDir(cmd)
	st, err := state.Load(dir)
	if err != nil {
		return err
	}
	weights := env.ListWeights(st, envName)
	if len(weights) == 0 {
		fmt.Fprintln(os.Stdout, "no weights set")
		return nil
	}
	patches := make([]string, 0, len(weights))
	for p := range weights {
		patches = append(patches, p)
	}
	sort.Strings(patches)
	for _, p := range patches {
		fmt.Fprintf(os.Stdout, "%-30s %d\n", p, weights[p])
	}
	return nil
}

func runWeightClear(cmd *cobra.Command, args []string) error {
	envName, patch := args[0], args[1]
	dir := workDir(cmd)
	st, err := state.Load(dir)
	if err != nil {
		return err
	}
	env.ClearWeight(st, envName, patch)
	if err := state.Save(st, dir); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "weight cleared for patch %q in %q\n", patch, envName)
	return nil
}
