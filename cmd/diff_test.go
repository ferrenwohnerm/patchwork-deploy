package cmd_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/patchwork-deploy/internal/state"
)

func tempDiffDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "diff-cmd-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestDiffCmd_ShowsDifferences(t *testing.T) {
	dir := tempDiffDir(t)
	stateFile := filepath.Join(dir, "state.json")
	cfgFile := filepath.Join(dir, "patchwork.yaml")

	cfgContent := `state_file: ` + stateFile + `
environments:
  - name: prod
    patch_dir: ` + dir + `/patches/prod
  - name: staging
    patch_dir: ` + dir + `/patches/staging
`
	if err := os.WriteFile(cfgFile, []byte(cfgContent), 0644); err != nil {
		t.Fatal(err)
	}

	st := &state.State{}
	st.Add(state.Record{Patch: "001_init.sql", Environment: "prod", AppliedAt: "2024-01-01T00:00:00Z"})
	st.Add(state.Record{Patch: "001_init.sql", Environment: "staging", AppliedAt: "2024-01-01T00:00:00Z"})
	st.Add(state.Record{Patch: "002_users.sql", Environment: "prod", AppliedAt: "2024-01-02T00:00:00Z"})
	if err := state.Save(stateFile, st); err != nil {
		t.Fatal(err)
	}

	rootCmd.SetArgs([]string{"--config", cfgFile, "diff", "--env-a", "prod", "--env-b", "staging"})
	var out strings.Builder
	rootCmd.SetOut(&out)

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDiffCmd_NoDifferences(t *testing.T) {
	dir := tempDiffDir(t)
	stateFile := filepath.Join(dir, "state.json")
	cfgFile := filepath.Join(dir, "patchwork.yaml")

	cfgContent := `state_file: ` + stateFile + `
environments:
  - name: prod
    patch_dir: ` + dir + `/patches/prod
  - name: staging
    patch_dir: ` + dir + `/patches/staging
`
	if err := os.WriteFile(cfgFile, []byte(cfgContent), 0644); err != nil {
		t.Fatal(err)
	}

	st := &state.State{}
	st.Add(state.Record{Patch: "001_init.sql", Environment: "prod", AppliedAt: "2024-01-01T00:00:00Z"})
	st.Add(state.Record{Patch: "001_init.sql", Environment: "staging", AppliedAt: "2024-01-01T00:00:00Z"})
	if err := state.Save(stateFile, st); err != nil {
		t.Fatal(err)
	}

	rootCmd.SetArgs([]string{"--config", cfgFile, "diff", "--env-a", "prod", "--env-b", "staging"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
