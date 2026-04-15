package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourorg/patchwork-deploy/internal/lock"
)

var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Manage the deployment lock",
}

var lockAcquireCmd = &cobra.Command{
	Use:   "acquire",
	Short: "Acquire the deployment lock",
	RunE:  runLockAcquire,
}

var lockReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Release the deployment lock",
	RunE:  runLockRelease,
}

var lockStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Report whether the deployment lock is currently held",
	RunE:  runLockStatus,
}

func init() {
	lockCmd.AddCommand(lockAcquireCmd)
	lockCmd.AddCommand(lockReleaseCmd)
	lockCmd.AddCommand(lockStatusCmd)
	rootCmd.AddCommand(lockCmd)
}

func runLockAcquire(cmd *cobra.Command, _ []string) error {
	dir := workDir(cmd)
	l := lock.New(dir)
	if err := l.Acquire(); err != nil {
		if errors.Is(err, lock.ErrLocked) {
			fmt.Fprintln(os.Stderr, "lock already held:", err)
			os.Exit(1)
		}
		return err
	}
	fmt.Println("lock acquired")
	return nil
}

func runLockRelease(cmd *cobra.Command, _ []string) error {
	l := lock.New(workDir(cmd))
	if err := l.Release(); err != nil {
		return err
	}
	fmt.Println("lock released")
	return nil
}

func runLockStatus(cmd *cobra.Command, _ []string) error {
	l := lock.New(workDir(cmd))
	if l.IsHeld() {
		fmt.Println("status: locked")
	} else {
		fmt.Println("status: unlocked")
	}
	return nil
}

func workDir(cmd *cobra.Command) string {
	if d, _ := cmd.Flags().GetString("dir"); d != "" {
		return d
	}
	wd, _ := os.Getwd()
	return wd
}
