package env_test

import (
	"testing"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
	"github.com/patchwork-deploy/internal/state/record"
)

func relayBaseState() *state.State {
	st := state.New()
	st.Add(record.Record{Environment: "staging", Patch: "001-init"})
	st.Add(record.Record{Environment: "staging", Patch: "002-schema"})
	st.Add(record.Record{Environment: "prod", Patch: "001-init"})
	return st
}

func TestSetRelay_Success(t *testing.T) {
	st := relayBaseState()
	err := env.SetRelay(st, "staging", "001-init", "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	target, ok := env.GetRelay(st, "staging", "001-init")
	if !ok {
		t.Fatal("expected relay to be set")
	}
	if target != "prod" {
		t.Errorf("expected target %q, got %q", "prod", target)
	}
}

func TestSetRelay_SameEnvReturnsError(t *testing.T) {
	st := relayBaseState()
	err := env.SetRelay(st, "staging", "001-init", "staging")
	if err == nil {
		t.Fatal("expected error for same env relay")
	}
}

func TestSetRelay_MissingEnvReturnsError(t *testing.T) {
	st := relayBaseState()
	err := env.SetRelay(st, "ghost", "001-init", "prod")
	if err == nil {
		t.Fatal("expected error for missing environment")
	}
}

func TestSetRelay_MissingPatchReturnsError(t *testing.T) {
	st := relayBaseState()
	err := env.SetRelay(st, "staging", "999-missing", "prod")
	if err == nil {
		t.Fatal("expected error for missing patch")
	}
}

func TestRemoveRelay_ClearsEntry(t *testing.T) {
	st := relayBaseState()
	_ = env.SetRelay(st, "staging", "001-init", "prod")
	err := env.RemoveRelay(st, "staging", "001-init")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := env.GetRelay(st, "staging", "001-init")
	if ok {
		t.Error("expected relay to be cleared")
	}
}

func TestListRelays_ReturnsAllForEnv(t *testing.T) {
	st := relayBaseState()
	_ = env.SetRelay(st, "staging", "001-init", "prod")
	_ = env.SetRelay(st, "staging", "002-schema", "prod")
	relays := env.ListRelays(st, "staging")
	if len(relays) != 2 {
		t.Errorf("expected 2 relays, got %d", len(relays))
	}
	if relays["001-init"] != "prod" {
		t.Errorf("unexpected relay target: %q", relays["001-init"])
	}
}
