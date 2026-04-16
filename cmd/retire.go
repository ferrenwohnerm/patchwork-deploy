package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

var retireCmd = &cobra.Command{
	Use:   "retire",
	Short: "Manage retired environments",
}

var retireEnvCmd = &cobra.Command{
	Use:   "env <environment>",
	Short: "Retire an environment, preserving its history",
	Args:  cobra.ExactArgs(1),
	RunE:  runRetire,
}

var listRetiredCmd = &cobra.Command{
	Use:   "list",
	Short: "List all retired environments",
	RunE:  runListRetired,
}

func init() {
	retireCmd.AddCommand(retireEnvCmd)
	retireCmd.AddCommand(listRetiredCmd)
	RootCmd.AddCommand(retireCmd)
}

func runRetire(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	environment := args[0]
	st, err := state.Load(cfg.StateFile)
	if err != nil {
		return fmt.Errorf("load state: %w", err)
	}
	if err := env.Retire(st, environment); err != nil {
		return err
	}
	if err := state.Save(cfg.StateFile, st); err != nil {
		return fmt.Errorf("save state: %w", err)
	}
	fmt.Fprintf(os.Stdout, "environment %q retired\n", environment)
	return nil
}

func runListRetired(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	st, err := state.Load(cfg.StateFile)
	if err != nil {
		return fmt.Errorf("load state: %w", err)
	}
	names := env.ListRetired(st)
	if len(names) == 0 {
		fmt.Fprintln(os.Stdout, "no retired environments")
		return nil
	}
	fmt.Fprintln(os.Stdout, strings.Join(names, "\n"))
	return nil
}
