package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func tempWatchDir(t *testing.T) string {
	t.Helper()
	return t.TempDir()
}

func writeWatchState(t *testing.T, dir string, records []state.Record) {
	t.Helper()
	data, _ := json.Marshal(records)
	if err := os.WriteFile(filepath.Join(dir, "state.json"), data, 0644); err != nil {
		t.Fatal(err)
	}
}

func TestWatchCmd_NoDrift(t *testing.T) {
	dir := tempWatchDir(t)
	records := []state.Record{
		{Environment: "prod", Patch: "001-init", AppliedAt: time.Now()},
	}
	writeWatchState(t, dir, records)

	cfgPath := writeTempConfig(t, dir, "prod")

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"watch", "prod", "--config", cfgPath})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out := buf.String(); out == "" {
		t.Error("expected output from watch command")
	}
}

func TestWatchCmd_UnknownEnv(t *testing.T) {
	dir := tempWatchDir(t)
	writeWatchState(t, dir, nil)
	cfgPath := writeTempConfig(t, dir, "prod")

	rootCmd.SetArgs([]string{"watch", "staging", "--config", cfgPath})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for unknown environment")
	}
}
