package env_test

import (
	"testing"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
)

func vectorBaseState() *state.State {
	st := state.New()
	st.AddRecord(state.Record{Environment: "staging", Patch: "001-init"})
	st.AddRecord(state.Record{Environment: "staging", Patch: "002-schema"})
	return st
}

func TestSetVector_Success(t *testing.T) {
	st := vectorBaseState()
	if err := env.SetVector(st, "staging", "001-init", "canary"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := env.GetVector(st, "staging", "001-init"); got != "canary" {
		t.Errorf("expected %q, got %q", "canary", got)
	}
}

func TestSetVector_InvalidNameReturnsError(t *testing.T) {
	st := vectorBaseState()
	if err := env.SetVector(st, "staging", "001-init", "BAD NAME!"); err == nil {
		t.Error("expected error for invalid vector name")
	}
}

func TestSetVector_MissingEnvReturnsError(t *testing.T) {
	st := vectorBaseState()
	if err := env.SetVector(st, "prod", "001-init", "stable"); err == nil {
		t.Error("expected error for missing environment")
	}
}

func TestSetVector_MissingPatchReturnsError(t *testing.T) {
	st := vectorBaseState()
	if err := env.SetVector(st, "staging", "999-missing", "stable"); err == nil {
		t.Error("expected error for missing patch")
	}
}

func TestClearVector_RemovesEntry(t *testing.T) {
	st := vectorBaseState()
	_ = env.SetVector(st, "staging", "001-init", "canary")
	if err := env.ClearVector(st, "staging", "001-init"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := env.GetVector(st, "staging", "001-init"); got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestListVectors_ReturnsAllEntries(t *testing.T) {
	st := vectorBaseState()
	_ = env.SetVector(st, "staging", "001-init", "canary")
	_ = env.SetVector(st, "staging", "002-schema", "stable")
	vecs := env.ListVectors(st, "staging")
	if len(vecs) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(vecs))
	}
	if vecs["001-init"] != "canary" {
		t.Errorf("expected canary, got %q", vecs["001-init"])
	}
	if vecs["002-schema"] != "stable" {
		t.Errorf("expected stable, got %q", vecs["002-schema"])
	}
}

func TestGetVector_NoneSet(t *testing.T) {
	st := vectorBaseState()
	if got := env.GetVector(st, "staging", "001-init"); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}
