package env_test

import (
	"testing"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
)

func scopeBaseState() *state.State {
	st := state.New()
	st.AddRecord("staging", state.Record{Patch: "001-init"})
	st.AddRecord("staging", state.Record{Patch: "002-schema"})
	return st
}

func TestSetScope_Success(t *testing.T) {
	st := scopeBaseState()
	if err := env.SetScope(st, "staging", "001-init", "infra"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := env.GetScope(st, "staging", "001-init"); got != "infra" {
		t.Errorf("expected scope %q, got %q", "infra", got)
	}
}

func TestSetScope_InvalidNameReturnsError(t *testing.T) {
	st := scopeBaseState()
	if err := env.SetScope(st, "staging", "001-init", "bad scope!"); err == nil {
		t.Fatal("expected error for invalid scope name, got nil")
	}
}

func TestSetScope_MissingEnvReturnsError(t *testing.T) {
	st := scopeBaseState()
	if err := env.SetScope(st, "ghost", "001-init", "infra"); err == nil {
		t.Fatal("expected error for missing environment, got nil")
	}
}

func TestSetScope_MissingPatchReturnsError(t *testing.T) {
	st := scopeBaseState()
	if err := env.SetScope(st, "staging", "999-nope", "infra"); err == nil {
		t.Fatal("expected error for missing patch, got nil")
	}
}

func TestClearScope_RemovesEntry(t *testing.T) {
	st := scopeBaseState()
	_ = env.SetScope(st, "staging", "001-init", "app")
	if err := env.ClearScope(st, "staging", "001-init"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := env.GetScope(st, "staging", "001-init"); got != "" {
		t.Errorf("expected empty scope after clear, got %q", got)
	}
}

func TestListScopes_ReturnsAllEntries(t *testing.T) {
	st := scopeBaseState()
	_ = env.SetScope(st, "staging", "001-init", "infra")
	_ = env.SetScope(st, "staging", "002-schema", "db")
	scopes := env.ListScopes(st, "staging")
	if len(scopes) != 2 {
		t.Fatalf("expected 2 scopes, got %d", len(scopes))
	}
	if scopes["001-init"] != "infra" {
		t.Errorf("expected infra for 001-init, got %q", scopes["001-init"])
	}
	if scopes["002-schema"] != "db" {
		t.Errorf("expected db for 002-schema, got %q", scopes["002-schema"])
	}
}

func TestGetScope_NoneSet(t *testing.T) {
	st := scopeBaseState()
	if got := env.GetScope(st, "staging", "001-init"); got != "" {
		t.Errorf("expected empty string when no scope set, got %q", got)
	}
}
