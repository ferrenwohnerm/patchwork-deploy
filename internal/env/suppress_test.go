package env

import (
	"testing"

	"github.com/patchwork-deploy/internal/state"
)

func suppressBaseState() *state.State {
	st := state.NewInMemory()
	st.Add("staging", state.Record{Patch: "001-init"})
	st.Add("staging", state.Record{Patch: "002-users"})
	return st
}

func TestSuppress_MarkesPatch(t *testing.T) {
	st := suppressBaseState()
	if err := Suppress(st, "staging", "001-init", "known issue"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !IsSuppressed(st, "staging", "001-init") {
		t.Fatal("expected patch to be suppressed")
	}
}

func TestSuppress_MissingEnvReturnsError(t *testing.T) {
	st := suppressBaseState()
	if err := Suppress(st, "prod", "001-init", "reason"); err == nil {
		t.Fatal("expected error for missing env")
	}
}

func TestSuppress_MissingPatchReturnsError(t *testing.T) {
	st := suppressBaseState()
	if err := Suppress(st, "staging", "999-missing", "reason"); err == nil {
		t.Fatal("expected error for missing patch")
	}
}

func TestSuppress_NewlineInReasonRejected(t *testing.T) {
	st := suppressBaseState()
	if err := Suppress(st, "staging", "001-init", "bad\nreason"); err == nil {
		t.Fatal("expected error for newline in reason")
	}
}

func TestUnsuppress_ClearsSentinel(t *testing.T) {
	st := suppressBaseState()
	_ = Suppress(st, "staging", "001-init", "temp")
	if err := Unsuppress(st, "staging", "001-init"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if IsSuppressed(st, "staging", "001-init") {
		t.Fatal("expected patch to no longer be suppressed")
	}
}

func TestListSuppressed_ReturnsAll(t *testing.T) {
	st := suppressBaseState()
	_ = Suppress(st, "staging", "001-init", "reason-a")
	_ = Suppress(st, "staging", "002-users", "reason-b")
	list := ListSuppressed(st, "staging")
	if len(list) != 2 {
		t.Fatalf("expected 2 suppressed patches, got %d", len(list))
	}
	if list["001-init"] != "reason-a" {
		t.Errorf("unexpected reason for 001-init: %s", list["001-init"])
	}
}
