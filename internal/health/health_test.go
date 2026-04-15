package health_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/health"
	"github.com/patchwork-deploy/internal/state"
)

func tempDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "health-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func TestCheck_NoStateFile(t *testing.T) {
	d := tempDir(t)
	s := health.Check("staging", filepath.Join(d, "missing.json"), filepath.Join(d, "lock"))
	if s.StateOK {
		t.Error("expected StateOK=false for missing state file")
	}
	if s.PatchCount != 0 {
		t.Errorf("expected PatchCount=0, got %d", s.PatchCount)
	}
	if s.Locked {
		t.Error("expected Locked=false")
	}
}

func TestCheck_WithPatches(t *testing.T) {
	d := tempDir(t)
	stateFile := filepath.Join(d, "state.json")

	st := state.New()
	st.Add(state.Record{Environment: "prod", Patch: "001", AppliedAt: time.Now().Add(-time.Hour)})
	st.Add(state.Record{Environment: "prod", Patch: "002", AppliedAt: time.Now()})
	st.Add(state.Record{Environment: "staging", Patch: "001", AppliedAt: time.Now()})
	if err := state.Save(stateFile, st); err != nil {
		t.Fatalf("Save: %v", err)
	}

	s := health.Check("prod", stateFile, filepath.Join(d, "lock"))
	if !s.StateOK {
		t.Error("expected StateOK=true")
	}
	if s.PatchCount != 2 {
		t.Errorf("expected PatchCount=2, got %d", s.PatchCount)
	}
	if s.LastApplied.IsZero() {
		t.Error("expected LastApplied to be set")
	}
}

func TestCheck_LockedState(t *testing.T) {
	d := tempDir(t)
	lockFile := filepath.Join(d, "patchwork.lock")
	if err := os.WriteFile(lockFile, []byte("42"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	s := health.Check("prod", filepath.Join(d, "state.json"), lockFile)
	if !s.Locked {
		t.Error("expected Locked=true")
	}
}

func TestStatus_String(t *testing.T) {
	s := health.Status{
		Environment: "staging",
		,
		LastApplied: time.Time{},
		StateOK:     true,
		Locked:      false,
	}
	out := s.String()
	if !strings.Contains(out, "staging") {
		t.Errorf("expected env name in output: %s", out)
	}
	if !strings.Contains(out, "patches=3") {
		t.Errorf("expected patch count in output: %s", out)
	}
	if !strings.Contains(out, "never") {
		t.Errorf("expected 'never' for zero LastApplied: %s", out)
	}
}
