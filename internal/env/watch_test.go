package env

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func watchBaseState() *state.State {
	st := state.New()
	st.Add(state.Record{Environment: "prod", Patch: "001-init", AppliedAt: time.Now()})
	st.Add(state.Record{Environment: "prod", Patch: "002-users", AppliedAt: time.Now()})
	return st
}

func TestWatch_ReturnsCorrectCount(t *testing.T) {
	st := watchBaseState()
	res, err := Watch(st, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.PatchCount != 2 {
		t.Errorf("expected 2 patches, got %d", res.PatchCount)
	}
	if res.LastApplied != "002-users" {
		t.Errorf("expected last patch 002-users, got %s", res.LastApplied)
	}
}

func TestWatch_MissingEnvReturnsError(t *testing.T) {
	st := watchBaseState()
	_, err := Watch(st, "staging")
	if err == nil {
		t.Fatal("expected error for missing environment")
	}
}

func TestWatch_NoDriftWhenNoSnapshot(t *testing.T) {
	st := watchBaseState()
	res, err := Watch(st, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Drifted {
		t.Error("expected no drift when snapshot is absent")
	}
}

func TestWatch_DetectsDrift(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("PATCHWORK_STATE_DIR", dir)

	st := watchBaseState()

	// Write snapshot with only first patch.
	snapPath := filepath.Join(dir, "snapshot-prod.json")
	data := `[{"environment":"prod","patch":"001-init","applied_at":"2024-01-01T00:00:00Z"}]`
	if err := os.WriteFile(snapPath, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}

	res, err := Watch(st, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Drifted {
		t.Error("expected drift to be detected")
	}
	if res.DriftDetail == "" {
		t.Error("expected non-empty drift detail")
	}
}
