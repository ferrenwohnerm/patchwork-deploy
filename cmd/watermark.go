package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/patchwork-deploy/internal/config"
	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

func init() {
	watermarkCmd := &cobra.Command{
		Use:   "watermark",
		Short: "Manage soft patch-count warning thresholds per environment",
	}

	setCmd := &cobra.Command{
		Use:   "set <env> <limit>",
		Short: "Set a watermark limit for an environment",
		Args:  cobra.ExactArgs(2),
		RunE:  runWatermarkSet,
	}

	clearCmd := &cobra.Command{
		Use:   "clear <env>",
		Short: "Remove the watermark for an environment",
		Args:  cobra.ExactArgs(1),
		RunE:  runWatermarkClear,
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all configured watermarks",
		Args:  cobra.NoArgs,
		RunE:  runWatermarkList,
	}

	watermarkCmd.AddCommand(setCmd, clearCmd, listCmd)
	rootCmd.AddCommand(watermarkCmd)
}

func runWatermarkSet(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return err
	}
	st, err := state.Load(workDir(cfg))
	if err != nil {
		return err
	}
	limit, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("limit must be an integer: %w", err)
	}
	if err := env.SetWatermark(st, args[0], limit); err != nil {
		return err
	}
	return state.Save(workDir(cfg), st)
}

func runWatermarkClear(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return err
	}
	st, err := state.Load(workDir(cfg))
	if err != nil {
		return err
	}
	if err := env.ClearWatermark(st, args[0]); err != nil {
		return err
	}
	return state.Save(workDir(cfg), st)
}

func runWatermarkList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return err
	}
	st, err := state.Load(workDir(cfg))
	if err != nil {
		return err
	}
	var envNames []string
	for _, e := range cfg.Environments {
		envNames = append(envNames, e.Name)
	}
	entries := env.CollectWatermarks(st, envNames)
	fmt.Fprintln(os.Stdout, env.FormatWatermarks(entries))
	return nil
}
