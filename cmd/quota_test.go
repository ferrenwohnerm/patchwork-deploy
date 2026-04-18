package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func tempQuotaDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "quota-cmd-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func writeQuotaState(t *testing.T, dir string) {
	t.Helper()
	st := state.New()
	now := time.Now()
	st.Add(state.Record{Environment: "staging", Patch: "001-init", AppliedAt: now})
	st.Add(state.Record{Environment: "staging", Patch: "002-users", AppliedAt: now})
	if err := state.Save(filepath.Join(dir, "state.json"), st); err != nil {
		t.Fatal(err)
	}
}

func TestQuotaSet_And_Check(t *testing.T) {
	dir := tempQuotaDir(t)
	writeQuotaState(t, dir)
	t.Setenv("PATCHWORK_DIR", dir)

	if err := runQuotaSet("staging", "10"); err != nil {
		t.Fatalf("set quota: %v", err)
	}

	out := captureStdout(t, func() {
		if err := runQuotaCheck("staging"); err != nil {
			t.Fatalf("check quota: %v", err)
		}
	})
	if !strings.Contains(out, "applied=2") || !strings.Contains(out, "limit=10") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestQuotaCheck_Exceeded(t *testing.T) {
	dir := tempQuotaDir(t)
	writeQuotaState(t, dir)
	t.Setenv("PATCHWORK_DIR", dir)

	_ = runQuotaSet("staging", "1")
	out := captureStdout(t, func() {
		_ = runQuotaCheck("staging")
	})
	if !strings.Contains(out, "EXCEEDED") {
		t.Errorf("expected EXCEEDED in output, got %q", out)
	}
}

func TestQuotaRemove_ClearsLimit(t *testing.T) {
	dir := tempQuotaDir(t)
	writeQuotaState(t, dir)
	t.Setenv("PATCHWORK_DIR", dir)

	_ = runQuotaSet("staging", "5")
	if err := runQuotaRemove("staging"); err != nil {
		t.Fatalf("remove quota: %v", err)
	}
	out := captureStdout(t, func() {
		_ = runQuotaCheck("staging")
	})
	if !strings.Contains(out, "limit=none") {
		t.Errorf("expected limit=none after remove, got %q", out)
	}
}

func TestQuotaSet_InvalidLimit(t *testing.T) {
	dir := tempQuotaDir(t)
	writeQuotaState(t, dir)
	t.Setenv("PATCHWORK_DIR", dir)
	if err := runQuotaSet("staging", "abc"); err == nil {
		t.Fatal("expected error for non-numeric limit")
	}
}
