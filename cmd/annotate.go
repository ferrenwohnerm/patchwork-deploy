package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

func init() {
	annotateCmd := &cobra.Command{Use: "annotate", Short: "Manage patch annotations"}

	setCmd := &cobra.Command{
		Use:   "set <env> <patch> <text>",
		Short: "Attach an annotation to a patch",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAnnotateSet(args[0], args[1], args[2])
		},
	}

	getCmd := &cobra.Command{
		Use:   "get <env> <patch>",
		Short: "Print the annotation for a patch",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAnnotateGet(args[0], args[1])
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <env> <patch>",
		Short: "Remove the annotation from a patch",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAnnotateRemove(args[0], args[1])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list <env>",
		Short: "List all annotations for an environment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAnnotateList(args[0])
		},
	}

	annotateCmd.AddCommand(setCmd, getCmd, removeCmd, listCmd)
	rootCmd.AddCommand(annotateCmd)
}

func runAnnotateSet(environment, patch, text string) error {
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	if err := env.SetAnnotation(st, environment, patch, text); err != nil {
		return err
	}
	return state.Save(st, workDir())
}

func runAnnotateGet(environment, patch string) error {
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	v, ok := env.GetAnnotation(st, environment, patch)
	if !ok {
		fmt.Fprintf(os.Stderr, "no annotation set for %s in %s\n", patch, environment)
		return nil
	}
	fmt.Println(v)
	return nil
}

func runAnnotateRemove(environment, patch string) error {
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	env.RemoveAnnotation(st, environment, patch)
	return state.Save(st, workDir())
}

func runAnnotateList(environment string) error {
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	list := env.ListAnnotations(st, environment)
	if len(list) == 0 {
		fmt.Println("no annotations")
		return nil
	}
	patches := make([]string, 0, len(list))
	for p := range list {
		patches = append(patches, p)
	}
	sort.Strings(patches)
	for _, p := range patches {
		fmt.Printf("%-30s  %s\n", p, list[p])
	}
	return nil
}
