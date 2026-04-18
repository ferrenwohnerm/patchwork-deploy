package env

import (
	"testing"

	"github.com/patchwork-deploy/internal/state"
)

func depsBaseState() *state.State {
	st := state.New()
	now := baseTime()
	st.Add(state.Record{Environment: "staging", Patch: "001-init", AppliedAt: now})
	st.Add(state.Record{Environment: "staging", Patch: "002-schema", AppliedAt: now})
	st.Add(state.Record{Environment: "staging", Patch: "003-seed", AppliedAt: now})
	return st
}

func TestSetDep_Success(t *testing.T) {
	st := depsBaseState()
	if err := SetDep(st, "staging", "002-schema", "001-init"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	deps, _ := GetDeps(st, "staging", "002-schema")
	if len(deps) != 1 || deps[0] != "001-init" {
		t.Errorf("expected [001-init], got %v", deps)
	}
}

func TestSetDep_Idempotent(t *testing.T) {
	st := depsBaseState()
	SetDep(st, "staging", "002-schema", "001-init")
	SetDep(st, "staging", "002-schema", "001-init")
	deps, _ := GetDeps(st, "staging", "002-schema")
	if len(deps) != 1 {
		t.Errorf("expected 1 dep, got %d", len(deps))
	}
}

func TestSetDep_MissingEnvReturnsError(t *testing.T) {
	st := depsBaseState()
	if err := SetDep(st, "prod", "002-schema", "001-init"); err == nil {
		t.Error("expected error for missing env")
	}
}

func TestSetDep_MissingPatchReturnsError(t *testing.T) {
	st := depsBaseState()
	if err := SetDep(st, "staging", "999-nope", "001-init"); err == nil {
		t.Error("expected error for missing patch")
	}
}

func TestGetDeps_NoneSet(t *testing.T) {
	st := depsBaseState()
	deps, err := GetDeps(st, "staging", "001-init")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(deps) != 0 {
		t.Errorf("expected no deps, got %v", deps)
	}
}

func TestRemoveDep_ClearsEdge(t *testing.T) {
	st := depsBaseState()
	SetDep(st, "staging", "003-seed", "001-init")
	SetDep(st, "staging", "003-seed", "002-schema")
	RemoveDep(st, "staging", "003-seed", "001-init")
	deps, _ := GetDeps(st, "staging", "003-seed")
	if len(deps) != 1 || deps[0] != "002-schema" {
		t.Errorf("expected [002-schema], got %v", deps)
	}
}
