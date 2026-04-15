package env_test

import (
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
)

func renameBaseState() *state.State {
	st := state.New()
	now := time.Now()
	st.Add(state.Record{Environment: "staging", Patch: "001-init.sql", AppliedAt: now})
	st.Add(state.Record{Environment: "staging", Patch: "002-users.sql", AppliedAt: now})
	st.Add(state.Record{Environment: "prod", Patch: "001-init.sql", AppliedAt: now})
	return st
}

func TestRename_MovesRecordsToNewEnv(t *testing.T) {
	st := renameBaseState()
	res, err := env.Rename(st, "staging", "staging-v2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Migrated != 2 {
		t.Errorf("expected 2 migrated, got %d", res.Migrated)
	}
	if res.OldName != "staging" || res.NewName != "staging-v2" {
		t.Errorf("unexpected result names: %+v", res)
	}
	if got := st.ForEnvironment("staging-v2"); len(got) != 2 {
		t.Errorf("expected 2 records in staging-v2, got %d", len(got))
	}
	if got := st.ForEnvironment("staging"); len(got) != 0 {
		t.Errorf("expected 0 records in staging after rename, got %d", len(got))
	}
}

func TestRename_LeavesOtherEnvsUntouched(t *testing.T) {
	st := renameBaseState()
	_, err := env.Rename(st, "staging", "staging-v2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := st.ForEnvironment("prod"); len(got) != 1 {
		t.Errorf("expected prod to remain intact, got %d records", len(got))
	}
}

func TestRename_SameEnvReturnsError(t *testing.T) {
	st := renameBaseState()
	_, err := env.Rename(st, "staging", "staging")
	if err == nil {
		t.Fatal("expected error for identical env names, got nil")
	}
}

func TestRename_MissingSourceReturnsError(t *testing.T) {
	st := renameBaseState()
	_, err := env.Rename(st, "nonexistent", "new-env")
	if err == nil {
		t.Fatal("expected error for missing source env, got nil")
	}
}

func TestRename_ExistingDestinationReturnsError(t *testing.T) {
	st := renameBaseState()
	_, err := env.Rename(st, "staging", "prod")
	if err == nil {
		t.Fatal("expected error when destination already exists, got nil")
	}
}
