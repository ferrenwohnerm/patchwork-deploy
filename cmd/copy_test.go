package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func tempCopyDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "copy-cmd-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestCopyCmd_CopiesPatches(t *testing.T) {
	dir := tempCopyDir(t)
	statePath := filepath.Join(dir, "state.json")

	st := state.New()
	base := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	st.Add(state.Record{Environment: "staging", PatchID: "001-init", AppliedAt: base})
	st.Add(state.Record{Environment: "staging", PatchID: "002-schema", AppliedAt: base.Add(time.Hour)})
	if err := state.Save(statePath, st); err != nil {
		t.Fatal(err)
	}

	cfgPath := writeTempConfig(t, dir, "staging", "prod")

	out, err := executeRoot(t, []string{"--config", cfgPath, "copy", "staging", "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "copied") {
		t.Errorf("expected 'copied' in output, got: %s", out)
	}

	loaded, _ := state.Load(statePath)
	prodRecs := loaded.ForEnvironment("prod")
	if len(prodRecs) != 2 {
		t.Errorf("expected 2 prod records after copy, got %d", len(prodRecs))
	}
}

func TestCopyCmd_SameEnvFails(t *testing.T) {
	dir := tempCopyDir(t)
	statePath := filepath.Join(dir, "state.json")

	st := state.New()
	st.Add(state.Record{Environment: "staging", PatchID: "001-init", AppliedAt: time.Now()})
	if err := state.Save(statePath, st); err != nil {
		t.Fatal(err)
	}

	cfgPath := writeTempConfig(t, dir, "staging")
	_, err := executeRoot(t, []string{"--config", cfgPath, "copy", "staging", "staging"})
	if err == nil {
		t.Fatal("expected error for same src/dst, got nil")
	}
}
