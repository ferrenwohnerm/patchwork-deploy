package cmd

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func tempRolloutDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "rollout-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func writeRolloutState(t *testing.T, path string) {
	t.Helper()
	st := state.New()
	now := time.Now()
	st.Add(state.Record{Environment: "staging", Patch: "001-init", AppliedAt: now})
	if err := state.Save(path, st); err != nil {
		t.Fatal(err)
	}
}

func TestRolloutSet_And_List(t *testing.T) {
	dir := tempRolloutDir(t)
	sf := filepath.Join(dir, "state.json")
	writeRolloutState(t, sf)

	if err := runRolloutSet(sf, "staging", "001-init", "canary"); err != nil {
		t.Fatalf("set failed: %v", err)
	}
	if err := runRolloutList(sf, "staging"); err != nil {
		t.Fatalf("list failed: %v", err)
	}
}

func TestRolloutClear_RemovesEntry(t *testing.T) {
	dir := tempRolloutDir(t)
	sf := filepath.Join(dir, "state.json")
	writeRolloutState(t, sf)

	_ = runRolloutSet(sf, "staging", "001-init", "immediate")
	if err := runRolloutClear(sf, "staging", "001-init"); err != nil {
		t.Fatalf("clear failed: %v", err)
	}
}

func TestRolloutSet_InvalidStrategyFails(t *testing.T) {
	dir := tempRolloutDir(t)
	sf := filepath.Join(dir, "state.json")
	writeRolloutState(t, sf)

	if err := runRolloutSet(sf, "staging", "001-init", "rolling"); err == nil {
		t.Fatal("expected error for invalid strategy")
	}
}
