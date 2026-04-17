package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func tempPromoteDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "promote-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func writePromoteState(t *testing.T, dir, env, patch string) {
	t.Helper()
	content := env + "," + patch + ",2024-01-01T00:00:00Z,ok,\n"
	err := os.WriteFile(filepath.Join(dir, "state.csv"), []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPromoteCmd_SameEnvFails(t *testing.T) {
	dir := tempPromoteDir(t)
	writePromoteState(t, dir, "staging", "001-init.sql")

	cmd := &cobra.Command{}
	cmd.Flags().String("state", filepath.Join(dir, "state.csv"), "")
	cmd.Flags().String("source", "staging", "")
	cmd.Flags().String("target", "staging", "")
	cmd.Flags().Bool("dry-run", false, "")

	err := runPromote(cmd, []string{})
	if err == nil {
		t.Fatal("expected error when source == target")
	}
}

func TestPromoteCmd_MissingSourceFails(t *testing.T) {
	dir := tempPromoteDir(t)

	cmd := &cobra.Command{}
	cmd.Flags().String("state", filepath.Join(dir, "state.csv"), "")
	cmd.Flags().String("source", "dev", "")
	cmd.Flags().String("target", "staging", "")
	cmd.Flags().Bool("dry-run", false, "")

	err := runPromote(cmd, []string{})
	if err == nil {
		t.Fatal("expected error for missing source environment")
	}
}
