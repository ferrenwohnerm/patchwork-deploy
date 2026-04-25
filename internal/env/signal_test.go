package env

import (
	"testing"

	"github.com/patchwork-deploy/internal/state"
)

func signalBaseState() *state.State {
	st := state.NewInMemory()
	st.Add("prod", "001-init.sql", "")
	st.Add("prod", "002-schema.sql", "")
	return st
}

func TestSetSignal_Success(t *testing.T) {
	st := signalBaseState()
	if err := SetSignal(st, "prod", "001-init.sql", "ready", "true"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := GetSignal(st, "prod", "001-init.sql", "ready")
	if !ok || v != "true" {
		t.Errorf("expected signal value 'true', got %q (ok=%v)", v, ok)
	}
}

func TestSetSignal_MissingEnvReturnsError(t *testing.T) {
	st := signalBaseState()
	if err := SetSignal(st, "staging", "001-init.sql", "ready", "true"); err == nil {
		t.Fatal("expected error for missing environment")
	}
}

func TestSetSignal_MissingPatchReturnsError(t *testing.T) {
	st := signalBaseState()
	if err := SetSignal(st, "prod", "999-missing.sql", "ready", "true"); err == nil {
		t.Fatal("expected error for missing patch")
	}
}

func TestSetSignal_InvalidNameReturnsError(t *testing.T) {
	st := signalBaseState()
	if err := SetSignal(st, "prod", "001-init.sql", "bad name!", "v"); err == nil {
		t.Fatal("expected error for invalid signal name")
	}
}

func TestSetSignal_NewlineInValueRejected(t *testing.T) {
	st := signalBaseState()
	if err := SetSignal(st, "prod", "001-init.sql", "ready", "val\nue"); err == nil {
		t.Fatal("expected error for newline in value")
	}
}

func TestRemoveSignal_ClearsEntry(t *testing.T) {
	st := signalBaseState()
	_ = SetSignal(st, "prod", "001-init.sql", "ready", "true")
	RemoveSignal(st, "prod", "001-init.sql", "ready")
	_, ok := GetSignal(st, "prod", "001-init.sql", "ready")
	if ok {
		t.Error("expected signal to be removed")
	}
}

func TestListSignals_ReturnsAll(t *testing.T) {
	st := signalBaseState()
	_ = SetSignal(st, "prod", "001-init.sql", "ready", "true")
	_ = SetSignal(st, "prod", "001-init.sql", "acked", "ops-team")
	sigs := ListSignals(st, "prod", "001-init.sql")
	if len(sigs) != 2 {
		t.Errorf("expected 2 signals, got %d", len(sigs))
	}
	if sigs["ready"] != "true" || sigs["acked"] != "ops-team" {
		t.Errorf("unexpected signal values: %v", sigs)
	}
}
