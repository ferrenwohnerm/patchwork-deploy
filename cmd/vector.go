package cmd

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

func init() {
	vectorCmd := &cobra.Command{
		Use:   "vector",
		Short: "Manage routing vectors for patches",
	}

	setCmd := &cobra.Command{
		Use:   "set <env> <patch> <vector>",
		Short: "Assign a routing vector to a patch",
		Args:  cobra.ExactArgs(3),
		RunE:  runVectorSet,
	}

	clearCmd := &cobra.Command{
		Use:   "clear <env> <patch>",
		Short: "Remove the routing vector from a patch",
		Args:  cobra.ExactArgs(2),
		RunE:  runVectorClear,
	}

	listCmd := &cobra.Command{
		Use:   "list <env>",
		Short: "List all routing vectors for an environment",
		Args:  cobra.ExactArgs(1),
		RunE:  runVectorList,
	}

	vectorCmd.AddCommand(setCmd, clearCmd, listCmd)
	rootCmd.AddCommand(vectorCmd)
}

func runVectorSet(cmd *cobra.Command, args []string) error {
	envName, patch, vector := args[0], args[1], args[2]
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	if err := env.SetVector(st, envName, patch, vector); err != nil {
		return err
	}
	if err := state.Save(st, workDir()); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "vector %q set for %s/%s\n", vector, envName, patch)
	return nil
}

func runVectorClear(cmd *cobra.Command, args []string) error {
	envName, patch := args[0], args[1]
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	if err := env.ClearVector(st, envName, patch); err != nil {
		return err
	}
	return state.Save(st, workDir())
}

func runVectorList(cmd *cobra.Command, args []string) error {
	envName := args[0]
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	vecs := env.ListVectors(st, envName)
	if len(vecs) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no vectors set")
		return nil
	}
	patches := make([]string, 0, len(vecs))
	for p := range vecs {
		patches = append(patches, p)
	}
	sort.Strings(patches)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PATCH\tVECTOR")
	for _, p := range patches {
		fmt.Fprintf(w, "%s\t%s\n", p, vecs[p])
	}
	return w.Flush()
}
