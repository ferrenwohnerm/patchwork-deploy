package cmd

import (
	"fmt"
	"os"

	"github.com/patchwork-deploy/internal/config"
	"github.com/patchwork-deploy/internal/diff"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

var diffEnvA, diffEnvB string

func init() {
	diffCmd := &cobra.Command{
		Use:   "diff",
		Short: "Compare applied patches between two environments",
		RunE:  runDiff,
	}
	diffCmd.Flags().StringVar(&diffEnvA, "env-a", "", "first environment (required)")
	diffCmd.Flags().StringVar(&diffEnvB, "env-b", "", "second environment (required)")
	_ = diffCmd.MarkFlagRequired("env-a")
	_ = diffCmd.MarkFlagRequired("env-b")
	rootCmd.AddCommand(diffCmd)
}

func runDiff(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if _, err := cfg.GetEnvironment(diffEnvA); err != nil {
		return fmt.Errorf("unknown environment: %s", diffEnvA)
	}
	if _, err := cfg.GetEnvironment(diffEnvB); err != nil {
		return fmt.Errorf("unknown environment: %s", diffEnvB)
	}

	st, err := state.Load(cfg.StateFile)
	if err != nil {
		return fmt.Errorf("load state: %w", err)
	}

	recsA := diff.FromStateRecords(st.ForEnvironment(diffEnvA))
	recsB := diff.FromStateRecords(st.ForEnvironment(diffEnvB))

	result := diff.Compare(recsA, recsB)
	fmt.Fprint(os.Stdout, diff.Format(result, diffEnvA, diffEnvB))
	return nil
}
