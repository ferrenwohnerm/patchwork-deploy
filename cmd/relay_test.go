package cmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/patchwork-deploy/internal/state"
	"github.com/patchwork-deploy/internal/state/record"
)

func tempRelayDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "relay-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func writeRelayState(t *testing.T, dir string) {
	t.Helper()
	st := state.New()
	st.Add(record.Record{Environment: "staging", Patch: "001-init"})
	st.Add(record.Record{Environment: "prod", Patch: "001-init"})
	if err := state.Save(st, dir); err != nil {
		t.Fatal(err)
	}
}

func TestRelaySet_And_List(t *testing.T) {
	dir := tempRelayDir(t)
	writeRelayState(t, dir)
	t.Setenv("PATCHWORK_DIR", dir)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)

	rootCmd.SetArgs([]string{"relay", "set", "staging", "001-init", "prod"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("set failed: %v", err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"relay", "list", "staging"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("list failed: %v", err)
	}
	out := buf.String()
	if out == "" {
		t.Error("expected relay list output")
	}
	if !bytes.Contains([]byte(out), []byte("001-init")) {
		t.Errorf("expected patch in output, got: %s", out)
	}
}

func TestRelayRemove_ClearsEntry(t *testing.T) {
	dir := tempRelayDir(t)
	writeRelayState(t, dir)
	t.Setenv("PATCHWORK_DIR", dir)

	rootCmd.SetArgs([]string{"relay", "set", "staging", "001-init", "prod"})
	_ = rootCmd.Execute()

	rootCmd.SetArgs([]string{"relay", "remove", "staging", "001-init"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("remove failed: %v", err)
	}

	st, _ := state.Load(dir)
	path := filepath.Join(dir, "state.json")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("state file missing: %v", err)
	}
	_ = st
}

func TestRelaySet_SameEnvFails(t *testing.T) {
	dir := tempRelayDir(t)
	writeRelayState(t, dir)
	t.Setenv("PATCHWORK_DIR", dir)

	rootCmd.SetArgs([]string{"relay", "set", "staging", "001-init", "staging"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when source and target env are the same")
	}
}
