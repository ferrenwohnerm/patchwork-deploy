package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/patchwork-deploy/internal/config"
	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

var scheduleCmd = &cobra.Command{Use: "schedule", Short: "Manage scheduled patch deployments"}

var scheduleAddCmd = &cobra.Command{
	Use:   "add <env> <patch> <run-at>",
	Short: "Schedule a patch to run at a given time (RFC3339)",
	Args:  cobra.ExactArgs(3),
	RunE:  runScheduleAdd,
}

var scheduleListCmd = &cobra.Command{
	Use:   "list <env>",
	Short: "List scheduled patches for an environment",
	Args:  cobra.ExactArgs(1),
	RunE:  runScheduleList,
}

var scheduleCancelCmd = &cobra.Command{
	Use:   "cancel <env> <patch>",
	Short: "Cancel a scheduled patch",
	Args:  cobra.ExactArgs(2),
	RunE:  runScheduleCancel,
}

func init() {
	scheduleCmd.AddCommand(scheduleAddCmd, scheduleListCmd, scheduleCancelCmd)
	rootCmd.AddCommand(scheduleCmd)
}

func runScheduleAdd(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return err
	}
	runAt, err := time.Parse(time.RFC3339, args[2])
	if err != nil {
		return fmt.Errorf("invalid time format, use RFC3339: %w", err)
	}
	st, err := state.Load(workDir(cfg))
	if err != nil {
		return err
	}
	if err := env.Schedule(st, args[0], args[1], runAt); err != nil {
		return err
	}
	return state.Save(workDir(cfg), st)
}

func runScheduleList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return err
	}
	st, err := state.Load(workDir(cfg))
	if err != nil {
		return err
	}
	entries := env.ListScheduled(st, args[0])
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PATCH\tRUN AT")
	for _, e := range entries {
		fmt.Fprintf(w, "%s\t%s\n", e.Patch, e.RunAt.Format(time.RFC3339))
	}
	return w.Flush()
}

func runScheduleCancel(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return err
	}
	st, err := state.Load(workDir(cfg))
	if err != nil {
		return err
	}
	if err := env.CancelScheduled(st, args[0], args[1]); err != nil {
		return err
	}
	return state.Save(workDir(cfg), st)
}
