package cmd

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func tempWindowDir(t *testing.T) string {
	t.Helper()
	return t.TempDir()
}

func writeWindowState(t *testing.T, dir string) {
	t.Helper()
	st := state.NewInMemory()
	st.Add("prod", state.Record{Patch: "001-init", AppliedAt: time.Now()})
	st.Add("prod", state.Record{Patch: "002-schema", AppliedAt: time.Now()})
	if err := state.Save(st, filepath.Join(dir, "state.json")); err != nil {
		t.Fatalf("writing state: %v", err)
	}
}

func TestWindowSet_And_List(t *testing.T) {
	dir := tempWindowDir(t)
	writeWindowState(t, dir)

	out, err := executeCommand(rootCmd,
		"--state", filepath.Join(dir, "state.json"),
		"window", "set", "prod", "001-init", "08:00", "18:00")
	if err != nil {
		t.Fatalf("set failed: %v\n%s", err, out)
	}

	out, err = executeCommand(rootCmd,
		"--state", filepath.Join(dir, "state.json"),
		"window", "list", "prod")
	if err != nil {
		t.Fatalf("list failed: %v\n%s", err, out)
	}
	if !contains(out, "001-init") || !contains(out, "08:00") {
		t.Errorf("expected window in output, got:\n%s", out)
	}
}

func TestWindowClear_RemovesEntry(t *testing.T) {
	dir := tempWindowDir(t)
	writeWindowState(t, dir)

	_, _ = executeCommand(rootCmd,
		"--state", filepath.Join(dir, "state.json"),
		"window", "set", "prod", "001-init", "09:00", "17:00")

	_, err := executeCommand(rootCmd,
		"--state", filepath.Join(dir, "state.json"),
		"window", "clear", "prod", "001-init")
	if err != nil {
		t.Fatalf("clear failed: %v", err)
	}

	out, _ := executeCommand(rootCmd,
		"--state", filepath.Join(dir, "state.json"),
		"window", "list", "prod")
	if contains(out, "001-init") {
		t.Errorf("expected cleared entry to be absent, got:\n%s", out)
	}
}

func TestWindowSet_InvalidTimeFails(t *testing.T) {
	dir := tempWindowDir(t)
	writeWindowState(t, dir)

	_, err := executeCommand(rootCmd,
		"--state", filepath.Join(dir, "state.json"),
		"window", "set", "prod", "001-init", "bad", "18:00")
	if err == nil {
		t.Fatal("expected error for invalid time format")
	}
}
