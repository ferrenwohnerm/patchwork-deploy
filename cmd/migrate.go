package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/patchwork-deploy/internal/config"
	"github.com/patchwork-deploy/internal/state"
)

var (
	migrateEnv        string
	migrateLegacyFile string
	migratePrune      bool
)

func init() {
	migrateCmd.Flags().StringVarP(&migrateEnv, "env", "e", "", "target environment (required)")
	migrateCmd.Flags().StringVarP(&migrateLegacyFile, "import", "i", "", "path to legacy applied-patches file")
	migrateCmd.Flags().BoolVar(&migratePrune, "prune", false, "remove state records for patches no longer on disk")
	_ = migrateCmd.MarkFlagRequired("env")
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Import legacy state or prune stale records",
	RunE:  runMigrate,
}

func runMigrate(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	env, err := cfg.GetEnvironment(migrateEnv)
	if err != nil {
		return err
	}

	s, err := state.Load(cfg.StateFile)
	if err != nil {
		return fmt.Errorf("loading state: %w", err)
	}

	if migrateLegacyFile != "" {
		n, err := state.Import(s, migrateEnv, migrateLegacyFile)
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "Imported %d record(s) into environment %q\n", n, migrateEnv)
	}

	if migratePrune {
		n, err := state.Prune(s, migrateEnv, env.PatchDir)
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "Pruned %d stale record(s) from environment %q\n", n, migrateEnv)
	}

	if err := s.Save(cfg.StateFile); err != nil {
		return fmt.Errorf("saving state: %w", err)
	}
	return nil
}
