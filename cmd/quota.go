package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

var quotaCmd = &cobra.Command{
	Use:   "quota",
	Short: "Manage patch quotas for environments",
}

var quotaSetCmd = &cobra.Command{
	Use:   "set <env> <limit>",
	Short: "Set the maximum number of applied patches for an environment",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runQuotaSet(args[0], args[1])
	},
}

var quotaCheckCmd = &cobra.Command{
	Use:   "check <env>",
	Short: "Check quota usage for an environment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runQuotaCheck(args[0])
	},
}

var quotaRemoveCmd = &cobra.Command{
	Use:   "remove <env>",
	Short: "Remove quota limit for an environment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runQuotaRemove(args[0])
	},
}

func init() {
	quotaCmd.AddCommand(quotaSetCmd, quotaCheckCmd, quotaRemoveCmd)
	rootCmd.AddCommand(quotaCmd)
}

func runQuotaSet(envName, limitStr string) error {
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return fmt.Errorf("invalid limit %q: %w", limitStr, err)
	}
	dir := workDir()
	st, err := state.Load(dir)
	if err != nil {
		return err
	}
	if err := env.SetQuota(st, envName, limit); err != nil {
		return err
	}
	return state.Save(dir, st)
}

func runQuotaCheck(envName string) error {
	dir := workDir()
	st, err := state.Load(dir)
	if err != nil {
		return err
	}
	res, err := env.CheckQuota(st, envName)
	if err != nil {
		return err
	}
	if res.Limit == 0 {
		fmt.Fprintf(os.Stdout, "env=%s applied=%d limit=none\n", res.Environment, res.Applied)
	} else {
		exceeded := ""
		if res.Exceeded {
			exceeded = " [EXCEEDED]"
		}
		fmt.Fprintf(os.Stdout, "env=%s applied=%d limit=%d%s\n", res.Environment, res.Applied, res.Limit, exceeded)
	}
	return nil
}

func runQuotaRemove(envName string) error {
	dir := workDir()
	st, err := state.Load(dir)
	if err != nil {
		return err
	}
	_ = env.RemoveQuota(st, envName)
	return state.Save(dir, st)
}
