package env

import (
	"testing"

	"github.com/patchwork-deploy/internal/state"
)

func flagBaseState() *state.State {
	st := state.New()
	st.AddRecord("staging", state.Record{Patch: "001-init.sql"})
	return st
}

func TestSetFlag_Success(t *testing.T) {
	st := flagBaseState()
	if err := SetFlag(st, "staging", "feature-x"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !HasFlag(st, "staging", "feature-x") {
		t.Fatal("expected flag to be set")
	}
}

func TestSetFlag_InvalidKey(t *testing.T) {
	st := flagBaseState()
	if err := SetFlag(st, "staging", "Bad Flag!"); err == nil {
		t.Fatal("expected error for invalid flag name")
	}
}

func TestSetFlag_MissingEnvReturnsError(t *testing.T) {
	st := flagBaseState()
	if err := SetFlag(st, "ghost", "ok-flag"); err == nil {
		t.Fatal("expected error for missing environment")
	}
}

func TestUnsetFlag_ClearsFlag(t *testing.T) {
	st := flagBaseState()
	_ = SetFlag(st, "staging", "feature-x")
	if err := UnsetFlag(st, "staging", "feature-x"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if HasFlag(st, "staging", "feature-x") {
		t.Fatal("expected flag to be cleared")
	}
}

func TestUnsetFlag_IdempotentWhenNotSet(t *testing.T) {
	st := flagBaseState()
	if err := UnsetFlag(st, "staging", "no-such-flag"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListFlags_ReturnsAllFlags(t *testing.T) {
	st := flagBaseState()
	_ = SetFlag(st, "staging", "alpha")
	_ = SetFlag(st, "staging", "beta")
	flags := ListFlags(st, "staging")
	if len(flags) != 2 {
		t.Fatalf("expected 2 flags, got %d", len(flags))
	}
}

func TestListFlags_DoesNotLeakOtherEnvs(t *testing.T) {
	st := state.New()
	st.AddRecord("staging", state.Record{Patch: "001.sql"})
	st.AddRecord("prod", state.Record{Patch: "001.sql"})
	_ = SetFlag(st, "prod", "prod-only")
	if flags := ListFlags(st, "staging"); len(flags) != 0 {
		t.Fatalf("expected no flags for staging, got %v", flags)
	}
}

func TestListFlags_EmptyForMissingEnv(t *testing.T) {
	st := flagBaseState()
	flags := ListFlags(st, "nonexistent")
	if len(flags) != 0 {
		t.Fatalf("expected no flags for missing env, got %v", flags)
	}
}
