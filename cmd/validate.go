package cmd

import (
	"fmt"
	"os"

	"github.com/patchwork-deploy/internal/config"
	"github.com/patchwork-deploy/internal/env"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate environment definitions in the config file",
	RunE:  runValidate,
}

func init() {
	validateCmd.Flags().StringP("config", "c", "patchwork.yaml", "path to config file")
	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, _ []string) error {
	cfgPath, _ := cmd.Flags().GetString("config")

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	results := env.ValidateAll(cfg)

	allValid := true
	for _, r := range results {
		if r.Valid() {
			fmt.Fprintln(cmd.OutOrStdout(), r.String())
		} else {
			fmt.Fprintln(cmd.ErrOrStderr(), r.String())
			allValid = false
		}
	}

	if !allValid {
		os.Exit(1)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "\nAll %d environment(s) valid.\n", len(results))
	return nil
}
