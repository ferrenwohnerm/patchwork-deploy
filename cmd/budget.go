package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

func init() {
	budgetCmd := &cobra.Command{
		Use:   "budget",
		Short: "Manage apply budgets for patches per environment",
	}

	setCmd := &cobra.Command{
		Use:   "set <env> <patch> <limit>",
		Short: "Set a maximum apply count for a patch",
		Args:  cobra.ExactArgs(3),
		RunE:  runBudgetSet,
	}

	clearCmd := &cobra.Command{
		Use:   "clear <env> <patch>",
		Short: "Remove the budget limit for a patch",
		Args:  cobra.ExactArgs(2),
		RunE:  runBudgetClear,
	}

	listCmd := &cobra.Command{
		Use:   "list <env>",
		Short: "List all budgets for an environment",
		Args:  cobra.ExactArgs(1),
		RunE:  runBudgetList,
	}

	budgetCmd.AddCommand(setCmd, clearCmd, listCmd)
	rootCmd.AddCommand(budgetCmd)
}

func runBudgetSet(cmd *cobra.Command, args []string) error {
	envName, patch, limitStr := args[0], args[1], args[2]
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return fmt.Errorf("limit must be an integer: %w", err)
	}
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	if err := env.SetBudget(st, envName, patch, limit); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "budget set: %s/%s = %d\n", envName, patch, limit)
	return state.Save(st, workDir())
}

func runBudgetClear(cmd *cobra.Command, args []string) error {
	envName, patch := args[0], args[1]
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	env.ClearBudget(st, envName, patch)
	fmt.Fprintf(os.Stdout, "budget cleared: %s/%s\n", envName, patch)
	return state.Save(st, workDir())
}

func runBudgetList(cmd *cobra.Command, args []string) error {
	envName := args[0]
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	budgets := env.ListBudgets(st, envName)
	records := st.ForEnvironment(envName)
	usedMap := map[string]int{}
	for _, r := range records {
		usedMap[r.Patch]++
	}
	entries := env.CollectBudgetEntries(budgets, usedMap)
	fmt.Fprintln(os.Stdout, env.FormatBudgets(entries))
	return nil
}
