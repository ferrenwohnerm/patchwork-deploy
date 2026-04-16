package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/patchwork-deploy/internal/config"
	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

var copyPatches []string

func init() {
	copyCmd := &cobra.Command{
		Use:   "copy <src> <dst>",
		Short: "Copy patch records from one environment to another",
		Args:  cobra.ExactArgs(2),
		RunE:  runCopy,
	}
	copyCmd.Flags().StringSliceVarP(&copyPatches, "patches", "p", nil, "comma-separated patch IDs to copy (default: all)")
	rootCmd.AddCommand(copyCmd)
}

func runCopy(cmd *cobra.Command, args []string) error {
	src, dst := args[0], args[1]

	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	st, err := state.Load(cfg.StateFile)
	if err != nil {
		return fmt.Errorf("load state: %w", err)
	}

	opts := env.CopyOptions{PatchIDs: copyPatches}
	res, err := env.Copy(st, src, dst, opts)
	if err != nil {
		return err
	}

	if err := state.Save(cfg.StateFile, st); err != nil {
		return fmt.Errorf("save state: %w", err)
	}

	if len(res.Copied) == 0 && len(res.Skipped) == 0 {
		fmt.Fprintln(os.Stdout, "nothing to copy")
		return nil
	}
	if len(res.Copied) > 0 {
		fmt.Fprintf(os.Stdout, "copied: %s\n", strings.Join(res.Copied, ", "))
	}
	if len(res.Skipped) > 0 {
		fmt.Fprintf(os.Stdout, "skipped (already present): %s\n", strings.Join(res.Skipped, ", "))
	}
	return nil
}
