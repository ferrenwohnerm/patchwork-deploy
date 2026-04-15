package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/user/patchwork-deploy/internal/audit"
)

var auditEnvFilter string

func init() {
	auditCmd := &cobra.Command{
		Use:   "audit",
		Short: "Display the audit log of patch actions",
		RunE:  runAudit,
	}
	auditCmd.Flags().StringVarP(&auditEnvFilter, "env", "e", "", "filter by environment")
	rootCmd.AddCommand(auditCmd)
}

func runAudit(cmd *cobra.Command, args []string) error {
	dir, err := workDir()
	if err != nil {
		return err
	}

	l := audit.New(dir)
	entries, err := l.Read()
	if err != nil {
		return fmt.Errorf("reading audit log: %w", err)
	}

	if len(entries) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no audit entries found")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tENVIRONMENT\tPATCH\tACTION\tSTATUS\tMESSAGE")

	for _, e := range entries {
		if auditEnvFilter != "" && e.Environment != auditEnvFilter {
			continue
		}
		status := "ok"
		if !e.Success {
			status = "fail"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			e.Timestamp.Format("2006-01-02T15:04:05Z"),
			e.Environment,
			e.Patch,
			e.Action,
			status,
			e.Message,
		)
	}
	return w.Flush()
}
