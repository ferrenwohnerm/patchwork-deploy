package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"patchwork-deploy/internal/notify"
)

var (
	notifyEnv   string
	notifyLevel string
	notifyMsg   string
)

func init() {
	notifyCmd := &cobra.Command{
		Use:   "notify",
		Short: "Emit a deployment notification event",
		RunE:  runNotify,
	}

	notifyCmd.Flags().StringVarP(&notifyEnv, "env", "e", "", "target environment (required)")
	notifyCmd.Flags().StringVarP(&notifyLevel, "level", "l", "info", "event level: info|warn|error")
	notifyCmd.Flags().StringVarP(&notifyMsg, "message", "m", "", "notification message (required)")
	_ = notifyCmd.MarkFlagRequired("env")
	_ = notifyCmd.MarkFlagRequired("message")

	rootCmd.AddCommand(notifyCmd)
}

func runNotify(cmd *cobra.Command, _ []string) error {
	n := notify.New(os.Stdout)

	var lvl notify.Level
	switch notifyLevel {
	case "warn":
		lvl = notify.LevelWarn
	case "error":
		lvl = notify.LevelError
	default:
		lvl = notify.LevelInfo
	}

	e := n.Send(notifyEnv, lvl, notifyMsg)
	fmt.Fprintf(cmd.OutOrStdout(), "event recorded at %s\n", e.Timestamp.Format("15:04:05"))
	return nil
}
