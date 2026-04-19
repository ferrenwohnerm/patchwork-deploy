package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/patchwork-deploy/internal/state"
)

func tempCapabilityDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "capability-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func writeCapabilityState(t *testing.T, dir string) {
	t.Helper()
	st := state.New()
	st.AddRecord("staging", state.Record{Patch: "001-init.sql"})
	if err := state.Save(dir, st); err != nil {
		t.Fatal(err)
	}
}

func TestCapabilityAdd_And_List(t *testing.T) {
	dir := tempCapabilityDir(t)
	writeCapabilityState(t, dir)
	t.Setenv("PATCHWORK_DIR", dir)

	if err := runCapabilityAdd("staging", "rollback"); err != nil {
		t.Fatalf("add failed: %v", err)
	}

	out := captureStdout(t, func() {
		_ = runCapabilityList("staging")
	})
	if !strings.Contains(out, "rollback") {
		t.Errorf("expected 'rollback' in output, got: %s", out)
	}
}

func TestCapabilityRemove_ClearsEntry(t *testing.T) {
	dir := tempCapabilityDir(t)
	writeCapabilityState(t, dir)
	t.Setenv("PATCHWORK_DIR", dir)

	_ = runCapabilityAdd("staging", "dry-run")
	if err := runCapabilityRemove("staging", "dry-run"); err != nil {
		t.Fatalf("remove failed: %v", err)
	}

	out := captureStdout(t, func() {
		_ = runCapabilityList("staging")
	})
	if strings.Contains(out, "dry-run") {
		t.Errorf("expected 'dry-run' to be removed, got: %s", out)
	}
}

func TestCapabilityAdd_UnknownEnvFails(t *testing.T) {
	dir := tempCapabilityDir(t)
	writeCapabilityState(t, dir)
	t.Setenv("PATCHWORK_DIR", dir)

	if err := runCapabilityAdd("ghost", "rollback"); err == nil {
		t.Fatal("expected error for unknown environment")
	}
}

func TestCapabilityStateFile_ExistsOnDisk(t *testing.T) {
	dir := tempCapabilityDir(t)
	writeCapabilityState(t, dir)
	t.Setenv("PATCHWORK_DIR", dir)

	_ = runCapabilityAdd("staging", "rollback")
	if _, err := os.Stat(filepath.Join(dir, "state.json")); err != nil {
		t.Fatalf("state file missing: %v", err)
	}
}
