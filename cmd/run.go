package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/patchwork-deploy/internal/config"
	"github.com/patchwork-deploy/internal/patch"
	"github.com/patchwork-deploy/internal/runner"
	"github.com/patchwork-deploy/internal/state"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run <environment>",
	Short: "Apply pending patches to the specified environment",
	Args:  cobra.ExactArgs(1),
	RunE:  runPatches,
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func runPatches(cmd *cobra.Command, args []string) error {
	envName := args[0]

	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	env, err := cfg.GetEnvironment(envName)
	if err != nil {
		return fmt.Errorf("unknown environment %q: %w", envName, err)
	}

	st, err := state.Load(cfg.StateFile)
	if err != nil {
		return fmt.Errorf("load state: %w", err)
	}

	loader := patch.NewLoader()
	applier := patch.NewApplier()
	r := runner.New(env, loader, applier, st)

	results, runErr := r.Run()
	summary := runner.Summarise(results)
	log.Printf("run complete: %s", summary)

	if saveErr := st.Save(cfg.StateFile); saveErr != nil {
		fmt.Fprintf(os.Stderr, "warning: could not save state: %v\n", saveErr)
	}

	if runErr != nil {
		return runErr
	}

	fmt.Println(summary)
	return nil
}
