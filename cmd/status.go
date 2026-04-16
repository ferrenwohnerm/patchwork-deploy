package cmd

import (
	"fmt"
	"os"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

var statusEnv string

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show deployment status for one or all environments",
	RunE:  runStatus,
}

func init() {
	statusCmd.Flags().StringVarP(&statusEnv, "env", "e", "", "environment name (omit for all)")
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, _ []string) error {
	cfg, err := loadConfig(cmd)
	if err != nil {
		return err
	}

	st, err := state.Load(cfg.StateFile)
	if err != nil {
		return fmt.Errorf("load state: %w", err)
	}

	if statusEnv != "" {
		s, err := env.Status(st, statusEnv)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}
		fmt.Println(s)
		if len(s.Tags) > 0 {
			fmt.Printf("  tags: %v\n", s.Tags)
		}
		return nil
	}

	all := env.StatusAll(st)
	if len(all) == 0 {
		fmt.Println("no environments found in state")
		return nil
	}
	for _, s := range all {
		fmt.Println(s)
	}
	return nil
}
