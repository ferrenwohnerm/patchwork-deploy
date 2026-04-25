package env

import (
	"testing"

	"github.com/patchwork-deploy/internal/state"
)

func slotBaseState() *state.State {
	st := state.NewInMemory()
	st.AddEnvironment("staging")
	st.AddRecord("staging", state.Record{Patch: "001-init"})
	return st
}

func TestSetSlot_Success(t *testing.T) {
	st := slotBaseState()
	if err := SetSlot(st, "staging", "001-init", "primary", "us-east-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := GetSlot(st, "staging", "001-init", "primary")
	if !ok || v != "us-east-1" {
		t.Errorf("expected slot value us-east-1, got %q (ok=%v)", v, ok)
	}
}

func TestSetSlot_MissingEnvReturnsError(t *testing.T) {
	st := slotBaseState()
	if err := SetSlot(st, "prod", "001-init", "primary", "us-east-1"); err == nil {
		t.Fatal("expected error for missing environment")
	}
}

func TestSetSlot_MissingPatchReturnsError(t *testing.T) {
	st := slotBaseState()
	if err := SetSlot(st, "staging", "999-nope", "primary", "us-east-1"); err == nil {
		t.Fatal("expected error for missing patch")
	}
}

func TestSetSlot_InvalidNameReturnsError(t *testing.T) {
	st := slotBaseState()
	if err := SetSlot(st, "staging", "001-init", "BAD SLOT!", "val"); err == nil {
		t.Fatal("expected error for invalid slot name")
	}
}

func TestSetSlot_NewlineInValueRejected(t *testing.T) {
	st := slotBaseState()
	if err := SetSlot(st, "staging", "001-init", "zone", "val\ninjected"); err == nil {
		t.Fatal("expected error for newline in value")
	}
}

func TestRemoveSlot_ClearsEntry(t *testing.T) {
	st := slotBaseState()
	_ = SetSlot(st, "staging", "001-init", "zone", "eu-west-1")
	RemoveSlot(st, "staging", "001-init", "zone")
	_, ok := GetSlot(st, "staging", "001-init", "zone")
	if ok {
		t.Error("expected slot to be removed")
	}
}

func TestListSlots_ReturnsAll(t *testing.T) {
	st := slotBaseState()
	_ = SetSlot(st, "staging", "001-init", "zone", "eu-west-1")
	_ = SetSlot(st, "staging", "001-init", "tier", "blue")
	slots := ListSlots(st, "staging", "001-init")
	if len(slots) != 2 {
		t.Errorf("expected 2 slots, got %d", len(slots))
	}
	if slots["zone"] != "eu-west-1" {
		t.Errorf("unexpected zone value: %q", slots["zone"])
	}
	if slots["tier"] != "blue" {
		t.Errorf("unexpected tier value: %q", slots["tier"])
	}
}
