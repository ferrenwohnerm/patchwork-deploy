package env_test

import (
	"testing"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
)

func affinityBaseState() *state.State {
	st := state.New()
	st.AddRecord("prod", state.Record{Patch: "001-init"})
	st.AddRecord("prod", state.Record{Patch: "002-schema"})
	return st
}

func TestSetAffinity_Success(t *testing.T) {
	st := affinityBaseState()
	if err := env.SetAffinity(st, "prod", "001-init", "db"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := env.GetAffinity(st, "prod", "001-init"); got != "db" {
		t.Errorf("expected %q, got %q", "db", got)
	}
}

func TestSetAffinity_InvalidGroupChars(t *testing.T) {
	st := affinityBaseState()
	if err := env.SetAffinity(st, "prod", "001-init", "bad group!"); err == nil {
		t.Fatal("expected error for invalid group name")
	}
}

func TestSetAffinity_MissingEnvReturnsError(t *testing.T) {
	st := affinityBaseState()
	if err := env.SetAffinity(st, "staging", "001-init", "db"); err == nil {
		t.Fatal("expected error for missing environment")
	}
}

func TestSetAffinity_MissingPatchReturnsError(t *testing.T) {
	st := affinityBaseState()
	if err := env.SetAffinity(st, "prod", "999-missing", "db"); err == nil {
		t.Fatal("expected error for missing patch")
	}
}

func TestRemoveAffinity_ClearsEntry(t *testing.T) {
	st := affinityBaseState()
	_ = env.SetAffinity(st, "prod", "001-init", "db")
	if err := env.RemoveAffinity(st, "prod", "001-init"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := env.GetAffinity(st, "prod", "001-init"); got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestListAffinities_ReturnsAll(t *testing.T) {
	st := affinityBaseState()
	_ = env.SetAffinity(st, "prod", "001-init", "db")
	_ = env.SetAffinity(st, "prod", "002-schema", "db")
	result := env.ListAffinities(st, "prod")
	if len(result) != 2 {
		t.Errorf("expected 2 affinities, got %d", len(result))
	}
	if result["001-init"] != "db" || result["002-schema"] != "db" {
		t.Errorf("unexpected affinity map: %v", result)
	}
}

func TestGetAffinity_NoneSet(t *testing.T) {
	st := affinityBaseState()
	if got := env.GetAffinity(st, "prod", "001-init"); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}
