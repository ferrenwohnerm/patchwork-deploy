package env

import (
	"testing"

	"github.com/patchwork-deploy/internal/state"
)

func intentBaseState() *state.State {
	st := state.NewInMemory()
	st.AddRecord("staging", state.Record{Patch: "001-init"})
	st.AddRecord("staging", state.Record{Patch: "002-feature"})
	return st
}

func TestSetIntent_Success(t *testing.T) {
	st := intentBaseState()
	if err := SetIntent(st, "staging", "001-init", "Bootstrap database schema"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := GetIntent(st, "staging", "001-init")
	if !ok || v != "Bootstrap database schema" {
		t.Errorf("expected intent to be set, got %q", v)
	}
}

func TestSetIntent_MissingEnvReturnsError(t *testing.T) {
	st := intentBaseState()
	if err := SetIntent(st, "prod", "001-init", "some intent"); err == nil {
		t.Fatal("expected error for missing env")
	}
}

func TestSetIntent_MissingPatchReturnsError(t *testing.T) {
	st := intentBaseState()
	if err := SetIntent(st, "staging", "999-missing", "some intent"); err == nil {
		t.Fatal("expected error for missing patch")
	}
}

func TestSetIntent_NewlineRejected(t *testing.T) {
	st := intentBaseState()
	if err := SetIntent(st, "staging", "001-init", "bad\nintent"); err == nil {
		t.Fatal("expected error for newline in intent")
	}
}

func TestRemoveIntent_ClearsValue(t *testing.T) {
	st := intentBaseState()
	_ = SetIntent(st, "staging", "001-init", "to remove")
	if err := RemoveIntent(st, "staging", "001-init"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := GetIntent(st, "staging", "001-init")
	if ok {
		t.Error("expected intent to be cleared")
	}
}

func TestListIntents_ReturnsAll(t *testing.T) {
	st := intentBaseState()
	_ = SetIntent(st, "staging", "001-init", "init intent")
	_ = SetIntent(st, "staging", "002-feature", "feature intent")
	result, err := ListIntents(st, "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 intents, got %d", len(result))
	}
	if result["001-init"] != "init intent" {
		t.Errorf("wrong value for 001-init: %q", result["001-init"])
	}
}
