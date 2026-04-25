package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func tempChainDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "chain-cmd-test")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func writeChainState(t *testing.T, dir string) {
	t.Helper()
	st := state.New()
	st.Add("prod", state.Record{Patch: "001-init", AppliedAt: time.Now()})
	st.Add("prod", state.Record{Patch: "002-schema", AppliedAt: time.Now()})
	if err := state.Save(st, dir); err != nil {
		t.Fatal(err)
	}
}

func TestChainSet_And_List(t *testing.T) {
	dir := tempChainDir(t)
	writeChainState(t, dir)
	t.Setenv("PATCHWORK_DIR", dir)

	out := &strings.Builder{}
	rootCmd.SetOut(out)

	rootCmd.SetArgs([]string{"chain", "set", "prod", "001-init", "002-schema"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("set failed: %v", err)
	}

	out.Reset()
	rootCmd.SetArgs([]string{"chain", "list", "prod"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if !strings.Contains(out.String(), "001-init") {
		t.Errorf("expected 001-init in output, got: %s", out.String())
	}
}

func TestChainRemove_ClearsEntry(t *testing.T) {
	dir := tempChainDir(t)
	writeChainState(t, dir)
	t.Setenv("PATCHWORK_DIR", dir)

	rootCmd.SetArgs([]string{"chain", "set", "prod", "001-init", "002-schema"})
	_ = rootCmd.Execute()

	rootCmd.SetArgs([]string{"chain", "remove", "prod", "001-init"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("remove failed: %v", err)
	}

	// Verify state file no longer contains the chain key
	data, _ := os.ReadFile(filepath.Join(dir, "state.json"))
	if strings.Contains(string(data), "chain:prod:001-init") {
		t.Error("expected chain entry to be removed from state")
	}
}

func TestChainSet_SameEnvFails(t *testing.T) {
	dir := tempChainDir(t)
	writeChainState(t, dir)
	t.Setenv("PATCHWORK_DIR", dir)

	rootCmd.SetArgs([]string{"chain", "set", "prod", "001-init", "001-init"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error for self-chain")
	}
}
