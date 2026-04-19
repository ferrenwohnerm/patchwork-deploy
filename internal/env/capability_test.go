package env

import (
	"testing"

	"github.com/patchwork-deploy/internal/state"
)

func capabilityBaseState() *state.State {
	st := state.New()
	st.AddRecord("staging", state.Record{Patch: "001-init.sql"})
	return st
}

func TestAddCapability_Success(t *testing.T) {
	st := capabilityBaseState()
	if err := AddCapability(st, "staging", "rollback"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !HasCapability(st, "staging", "rollback") {
		t.Fatal("expected capability to be set")
	}
}

func TestAddCapability_InvalidChars(t *testing.T) {
	st := capabilityBaseState()
	if err := AddCapability(st, "staging", "bad cap!"); err == nil {
		t.Fatal("expected error for invalid capability name")
	}
}

func TestAddCapability_MissingEnvReturnsError(t *testing.T) {
	st := capabilityBaseState()
	if err := AddCapability(st, "ghost", "rollback"); err == nil {
		t.Fatal("expected error for missing environment")
	}
}

func TestRemoveCapability_ClearsEntry(t *testing.T) {
	st := capabilityBaseState()
	_ = AddCapability(st, "staging", "rollback")
	_ = RemoveCapability(st, "staging", "rollback")
	if HasCapability(st, "staging", "rollback") {
		t.Fatal("expected capability to be removed")
	}
}

func TestListCapabilities_ReturnsAll(t *testing.T) {
	st := capabilityBaseState()
	_ = AddCapability(st, "staging", "rollback")
	_ = AddCapability(st, "staging", "dry-run")
	caps := ListCapabilities(st, "staging")
	if len(caps) != 2 {
		t.Fatalf("expected 2 capabilities, got %d", len(caps))
	}
}

func TestListCapabilities_EmptyWhenNone(t *testing.T) {
	st := capabilityBaseState()
	caps := ListCapabilities(st, "staging")
	if len(caps) != 0 {
		t.Fatalf("expected 0 capabilities, got %d", len(caps))
	}
}
