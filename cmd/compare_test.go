package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func tempCompareDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "compare-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func writeCompareState(t *testing.T, path string, records []state.Record) {
	t.Helper()
	data, _ := json.Marshal(records)
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
}

func TestCompareCmd_ShowsDifferences(t *testing.T) {
	dir := tempCompareDir(t)
	statePath := filepath.Join(dir, "state.json")
	now := time.Now()
	writeCompareState(t, statePath, []state.Record{
		{Environment: "staging", Patch: "001-init", AppliedAt: now},
		{Environment: "staging", Patch: "002-users", AppliedAt: now},
		{Environment: "prod", Patch: "001-init", AppliedAt: now},
	})

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"compare", "staging", "prod",
		"--config", "../patchwork.example.yaml",
		"--state", statePath,
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "002-users") {
		t.Errorf("expected 002-users in output, got: %s", out)
	}
}

func TestCompareCmd_SameEnvFails(t *testing.T) {
	dir := tempCompareDir(t)
	statePath := filepath.Join(dir, "state.json")
	now := time.Now()
	writeCompareState(t, statePath, []state.Record{
		{Environment: "staging", Patch: "001-init", AppliedAt: now},
	})

	rootCmd.SetArgs([]string{"compare", "staging", "staging",
		"--config", "../patchwork.example.yaml",
		"--state", statePath,
	})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error for same env")
	}
}
