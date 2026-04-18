package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"patchwork-deploy/internal/config"
	"patchwork-deploy/internal/state"
)

var (
	pruneEnv       string
	pruneOlderThan string
	pruneDryRun    bool
)

func init() {
	pruneCmd := &cobra.Command{
		Use:   "prune",
		Short: "Remove old state records by age or environment",
		RunE:  runPrune,
	}
	pruneCmd.Flags().StringVarP(&pruneEnv, "env", "e", "", "Restrict pruning to a specific environment")
	pruneCmd.Flags().StringVar(&pruneOlderThan, "older-than", "", "Remove records older than this duration (e.g. 720h)")
	pruneCmd.Flags().BoolVar(&pruneDryRun, "dry-run", false, "Show what would be removed without modifying state")
	_ = pruneCmd.MarkFlagRequired("older-than")
	rootCmd.AddCommand(pruneCmd)
}

func runPrune(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load("patchwork.yaml")
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	dur, err := time.ParseDuration(pruneOlderThan)
	if err != nil {
		return fmt.Errorf("invalid --older-than value %q: %w", pruneOlderThan, err)
	}

	st, err := state.Load(cfg.StateFile)
	if err != nil {
		return fmt.Errorf("loading state: %w", err)
	}

	opts := state.PruneOptions{
		OlderThan:   time.Now().Add(-dur),
		Environment: pruneEnv,
		DryRun:      pruneDryRun,
	}

	result := state.PruneByAge(st, opts)

	if pruneDryRun {
		printPrunePreview(result)
		return nil
	}

	if err := state.Save(st, cfg.StateFile); err != nil {
		return fmt.Errorf("saving state: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Pruned %d record(s), %d retained.\n", len(result.Removed), len(result.Retained))
	return nil
}

// printPrunePreview prints the records that would be removed in a dry-run.
func printPrunePreview(result state.PruneResult) {
	fmt.Fprintf(os.Stdout, "[dry-run] would remove %d record(s)\n", len(result.Removed))
	for _, r := range result.Removed {
		fmt.Fprintf(os.Stdout, "  - [%s] %s (applied %s)\n", r.Environment, r.Patch, r.AppliedAt.Format(time.RFC3339))
	}
}
