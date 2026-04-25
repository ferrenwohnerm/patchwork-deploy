package env

import (
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func chainBaseState() *state.State {
	st := state.New()
	st.Add("prod", state.Record{Patch: "001-init", AppliedAt: time.Now()})
	st.Add("prod", state.Record{Patch: "002-schema", AppliedAt: time.Now()})
	st.Add("prod", state.Record{Patch: "003-data", AppliedAt: time.Now()})
	return st
}

func TestSetChain_Success(t *testing.T) {
	st := chainBaseState()
	if err := SetChain(st, "prod", "001-init", "002-schema"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := GetChain(st, "prod", "001-init")
	if !ok || v != "002-schema" {
		t.Errorf("expected successor 002-schema, got %q (ok=%v)", v, ok)
	}
}

func TestSetChain_SelfReturnsError(t *testing.T) {
	st := chainBaseState()
	if err := SetChain(st, "prod", "001-init", "001-init"); err == nil {
		t.Fatal("expected error for self-chain")
	}
}

func TestSetChain_MissingEnvReturnsError(t *testing.T) {
	st := chainBaseState()
	if err := SetChain(st, "staging", "001-init", "002-schema"); err == nil {
		t.Fatal("expected error for missing env")
	}
}

func TestSetChain_MissingPatchReturnsError(t *testing.T) {
	st := chainBaseState()
	if err := SetChain(st, "prod", "999-missing", "002-schema"); err == nil {
		t.Fatal("expected error for missing patch")
	}
}

func TestRemoveChain_ClearsEntry(t *testing.T) {
	st := chainBaseState()
	_ = SetChain(st, "prod", "001-init", "002-schema")
	if err := RemoveChain(st, "prod", "001-init"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := GetChain(st, "prod", "001-init")
	if ok {
		t.Error("expected chain to be removed")
	}
}

func TestListChains_ReturnsAllEntries(t *testing.T) {
	st := chainBaseState()
	_ = SetChain(st, "prod", "001-init", "002-schema")
	_ = SetChain(st, "prod", "002-schema", "003-data")
	chains := ListChains(st, "prod")
	if len(chains) != 2 {
		t.Errorf("expected 2 chains, got %d", len(chains))
	}
	if chains["001-init"] != "002-schema" {
		t.Errorf("unexpected chain for 001-init: %q", chains["001-init"])
	}
}
