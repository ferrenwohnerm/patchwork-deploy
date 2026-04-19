package cmd

import (
	"fmt"
	"os"

	"github.com/patchwork-deploy/internal/config"
	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

var sealReason string

func init() {
	sealCmd := &cobra.Command{
		Use:   "seal <environment>",
		Short: "Seal an environment to prevent further patch application",
		Args:  cobra.ExactArgs(1),
		RunE:  runSeal,
	}
	sealCmd.Flags().StringVar(&sealReason, "reason", "", "Optional reason for sealing")

	unsealCmd := &cobra.Command{
		Use:   "unseal <environment>",
		Short: "Unseal a previously sealed environment",
		Args:  cobra.ExactArgs(1),
		RunE:  runUnseal,
	}

	rootCmd.AddCommand(sealCmd)
	rootCmd.AddCommand(unsealCmd)
}

func runSeal(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return err
	}
	st, err := state.Load(cfg.StateFile)
	if err != nil {
		return err
	}
	if err := env.Seal(st, args[0], sealReason); err != nil {
		return err
	}
	if err := st.Save(cfg.StateFile); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "environment %q sealed\n", args[0])
	if sealReason != "" {
		fmt.Fprintf(os.Stdout, "reason: %s\n", sealReason)
	}
	return nil
}

func runUnseal(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return err
	}
	st, err := state.Load(cfg.StateFile)
	if err != nil {
		return err
	}
	if err := env.Unseal(st, args[0]); err != nil {
		return err
	}
	if err := st.Save(cfg.StateFile); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "environment %q unsealed\n", args[0])
	return nil
}
