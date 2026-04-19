package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func tempPriorityDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "priority-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func writePriorityState(t *testing.T, dir string, st *state.State) {
	t.Helper()
	b, _ := json.Marshal(st)
	if err := os.WriteFile(filepath.Join(dir, "state.json"), b, 0644); err != nil {
		t.Fatal(err)
	}
}

func TestPrioritySet_And_List(t *testing.T) {
	dir := tempPriorityDir(t)
	st := state.New()
	st.Add(state.Record{Environment: "prod", Patch: "001-init", AppliedAt: time.Now()})
	writePriorityState(t, dir, st)

	root := newRootCmd()
	root.SetArgs([]string{"priority", "set", "prod", "001-init", "5", "--dir", dir})
	if err := root.Execute(); err != nil {
		t.Fatalf("set failed: %v", err)
	}

	out := captureOutput(t, func() {
		root2 := newRootCmd()
		root2.SetArgs([]string{"priority", "list", "prod", "--dir", dir})
		_ = root2.Execute()
	})
	if !contains(out, "001-init") {
		t.Errorf("expected patch in list output, got: %s", out)
	}
}

func TestPriorityClear_RemovesEntry(t *testing.T) {
	dir := tempPriorityDir(t)
	st := state.New()
	st.Add(state.Record{Environment: "prod", Patch: "001-init", AppliedAt: time.Now()})
	writePriorityState(t, dir, st)

	root := newRootCmd()
	root.SetArgs([]string{"priority", "set", "prod", "001-init", "3", "--dir", dir})
	_ = root.Execute()

	root2 := newRootCmd()
	root2.SetArgs([]string{"priority", "clear", "prod", "001-init", "--dir", dir})
	if err := root2.Execute(); err != nil {
		t.Fatalf("clear failed: %v", err)
	}

	out := captureOutput(t, func() {
		root3 := newRootCmd()
		root3.SetArgs([]string{"priority", "list", "prod", "--dir", dir})
		_ = root3.Execute()
	})
	if contains(out, "001-init") {
		t.Errorf("expected patch to be cleared from list, got: %s", out)
	}
}

func TestPrioritySet_InvalidNumber(t *testing.T) {
	dir := tempPriorityDir(t)
	root := newRootCmd()
	root.SetArgs([]string{"priority", "set", "prod", "001-init", "notanumber", "--dir", dir})
	if err := root.Execute(); err == nil {
		t.Fatal("expected error for non-integer priority")
	}
}
