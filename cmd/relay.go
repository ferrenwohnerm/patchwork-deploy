package cmd

import (
	"fmt"
	"os"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

var relayCmd = &cobra.Command{
	Use:   "relay",
	Short: "Manage patch relay forwarding between environments",
}

func init() {
	relayCmd.AddCommand(relaySetCmd, relayRemoveCmd, relayListCmd)
	rootCmd.AddCommand(relayCmd)
}

var relaySetCmd = &cobra.Command{
	Use:   "set <env> <patch> <target-env>",
	Short: "Forward patch events from env to target-env",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRelaySet(args[0], args[1], args[2])
	},
}

var relayRemoveCmd = &cobra.Command{
	Use:   "remove <env> <patch>",
	Short: "Remove relay configuration for a patch",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRelayRemove(args[0], args[1])
	},
}

var relayListCmd = &cobra.Command{
	Use:   "list <env>",
	Short: "List all relay mappings for an environment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRelayList(args[0])
	},
}

func runRelaySet(envName, patch, target string) error {
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	if err := env.SetRelay(st, envName, patch, target); err != nil {
		return err
	}
	return state.Save(st, workDir())
}

func runRelayRemove(envName, patch string) error {
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	if err := env.RemoveRelay(st, envName, patch); err != nil {
		return err
	}
	return state.Save(st, workDir())
}

func runRelayList(envName string) error {
	st, err := state.Load(workDir())
	if err != nil {
		return err
	}
	relays := env.ListRelays(st, envName)
	fmt.Fprint(os.Stdout, env.FormatRelays(envName, relays))
	return nil
}
