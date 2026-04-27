package env_test

import (
	"testing"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
)

func surgeBaseState() *state.State {
	st := state.NewInMemory()
	st.AddRecord("prod", state.Record{Patch: "001-init"})
	st.AddRecord("prod", state.Record{Patch: "002-add-index"})
	return st
}

func TestSetSurge_Success(t *testing.T) {
	st := surgeBaseState()
	if err := env.SetSurge(st, "prod", "001-init", 3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	n, ok := env.GetSurge(st, "prod", "001-init")
	if !ok {
		t.Fatal("expected surge to be set")
	}
	if n != 3 {
		t.Errorf("expected 3, got %d", n)
	}
}

func TestSetSurge_ZeroLimitReturnsError(t *testing.T) {
	st := surgeBaseState()
	if err := env.SetSurge(st, "prod", "001-init", 0); err == nil {
		t.Fatal("expected error for zero limit")
	}
}

func TestSetSurge_MissingEnvReturnsError(t *testing.T) {
	st := surgeBaseState()
	if err := env.SetSurge(st, "staging", "001-init", 2); err == nil {
		t.Fatal("expected error for missing environment")
	}
}

func TestSetSurge_MissingPatchReturnsError(t *testing.T) {
	st := surgeBaseState()
	if err := env.SetSurge(st, "prod", "999-missing", 2); err == nil {
		t.Fatal("expected error for missing patch")
	}
}

func TestGetSurge_NoneSet(t *testing.T) {
	st := surgeBaseState()
	_, ok := env.GetSurge(st, "prod", "001-init")
	if ok {
		t.Fatal("expected no surge to be set")
	}
}

func TestClearSurge_RemovesEntry(t *testing.T) {
	st := surgeBaseState()
	_ = env.SetSurge(st, "prod", "001-init", 5)
	env.ClearSurge(st, "prod", "001-init")
	_, ok := env.GetSurge(st, "prod", "001-init")
	if ok {
		t.Fatal("expected surge to be cleared")
	}
}

func TestListSurges_ReturnsAllEntries(t *testing.T) {
	st := surgeBaseState()
	_ = env.SetSurge(st, "prod", "001-init", 2)
	_ = env.SetSurge(st, "prod", "002-add-index", 4)
	entries := env.ListSurges(st, "prod")
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries["001-init"] != 2 {
		t.Errorf("expected 2 for 001-init, got %d", entries["001-init"])
	}
	if entries["002-add-index"] != 4 {
		t.Errorf("expected 4 for 002-add-index, got %d", entries["002-add-index"])
	}
}
