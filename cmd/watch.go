package cmd

import (
	"fmt"
	"os"

	"github.com/patchwork-deploy/internal/config"
	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:   "watch <environment>",
	Short: "Check an environment for drift against its last snapshot",
	Args:  cobra.ExactArgs(1),
	RunE:  runWatch,
}

func init() {
	rootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, args []string) error {
	cfgPath, _ := cmd.Flags().GetString("config")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	envName := args[0]
	if _, err := cfg.GetEnvironment(envName); err != nil {
		return fmt.Errorf("unknown environment: %s", envName)
	}

	st, err := state.Load(workDir(cfg))
	if err != nil {
		return fmt.Errorf("load state: %w", err)
	}

	res, err := env.Watch(st, envName)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Environment : %s\n", res.Environment)
	fmt.Fprintf(os.Stdout, "Patches     : %d\n", res.PatchCount)
	fmt.Fprintf(os.Stdout, "Last Applied: %s\n", res.LastApplied)
	fmt.Fprintf(os.Stdout, "Checked At  : %s\n", res.LastChecked.Format("2006-01-02T15:04:05Z"))

	if res.Drifted {
		fmt.Fprintf(os.Stdout, "Drift       : YES — %s\n", res.DriftDetail)
		os.Exit(1)
	} else {
		fmt.Fprintf(os.Stdout, "Drift       : none\n")
	}
	return nil
}
