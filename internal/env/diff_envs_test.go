package env_test

import (
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
)

func diffBaseState() *state.State {
	now := time.Now()
	st := &state.State{}
	st.Add(state.Record{Environment: "staging", Patch: "001-init.sql", AppliedAt: now})
	st.Add(state.Record{Environment: "staging", Patch: "002-users.sql", AppliedAt: now})
	st.Add(state.Record{Environment: "production", Patch: "001-init.sql", AppliedAt: now})
	return st
}

func TestDiffEnvironments_ShowsMissingInTarget(t *testing.T) {
	st := diffBaseState()

	result, err := env.DiffEnvironments(st, "staging", "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Source != "staging" || result.Target != "production" {
		t.Errorf("unexpected source/target: %s / %s", result.Source, result.Target)
	}

	if len(result.Result.OnlyInA) != 1 {
		t.Errorf("expected 1 patch only in staging, got %d", len(result.Result.OnlyInA))
	}

	if result.Result.OnlyInA[0].Patch != "002-users.sql" {
		t.Errorf("expected 002-users.sql only in staging, got %s", result.Result.OnlyInA[0].Patch)
	}
}

func TestDiffEnvironments_SameEnvReturnsError(t *testing.T) {
	st := diffBaseState()

	_, err := env.DiffEnvironments(st, "staging", "staging")
	if err == nil {
		t.Fatal("expected error for same source and target")
	}
}

func TestDiffEnvironments_MissingSourceReturnsError(t *testing.T) {
	st := diffBaseState()

	_, err := env.DiffEnvironments(st, "nonexistent", "production")
	if err == nil {
		t.Fatal("expected error for missing source environment")
	}
}

func TestDiffEnvironments_MissingTargetReturnsError(t *testing.T) {
	st := diffBaseState()

	_, err := env.DiffEnvironments(st, "staging", "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing target environment")
	}
}

func TestDiffEnvironments_BothHaveSamePatches(t *testing.T) {
	now := time.Now()
	st := &state.State{}
	st.Add(state.Record{Environment: "staging", Patch: "001-init.sql", AppliedAt: now})
	st.Add(state.Record{Environment: "production", Patch: "001-init.sql", AppliedAt: now})

	result, err := env.DiffEnvironments(st, "staging", "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Result.OnlyInA) != 0 || len(result.Result.OnlyInB) != 0 {
		t.Errorf("expected no differences, got onlyInA=%d onlyInB=%d",
			len(result.Result.OnlyInA), len(result.Result.OnlyInB))
	}

	if len(result.Result.InBoth) != 1 {
		t.Errorf("expected 1 shared patch, got %d", len(result.Result.InBoth))
	}
}
