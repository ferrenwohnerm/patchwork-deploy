package env

import (
	"testing"

	"github.com/patchwork-deploy/internal/state"
)

func ceilingBaseState() *state.State {
	st := state.NewInMemory()
	st.AddRecord("prod", state.Record{Patch: "001-init.sql"})
	st.AddRecord("prod", state.Record{Patch: "002-users.sql"})
	return st
}

func TestSetCeiling_Success(t *testing.T) {
	st := ceilingBaseState()
	if err := SetCeiling(st, "prod", 5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, err := GetCeiling(st, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != 5 {
		t.Errorf("expected 5, got %d", v)
	}
}

func TestSetCeiling_ZeroLimitReturnsError(t *testing.T) {
	st := ceilingBaseState()
	if err := SetCeiling(st, "prod", 0); err == nil {
		t.Fatal("expected error for zero limit")
	}
}

func TestSetCeiling_MissingEnvReturnsError(t *testing.T) {
	st := ceilingBaseState()
	if err := SetCeiling(st, "staging", 3); err == nil {
		t.Fatal("expected error for missing environment")
	}
}

func TestGetCeiling_NoneSet(t *testing.T) {
	st := ceilingBaseState()
	v, err := GetCeiling(st, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != 0 {
		t.Errorf("expected 0 when no ceiling set, got %d", v)
	}
}

func TestClearCeiling_RemovesLimit(t *testing.T) {
	st := ceilingBaseState()
	_ = SetCeiling(st, "prod", 10)
	if err := ClearCeiling(st, "prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, _ := GetCeiling(st, "prod")
	if v != 0 {
		t.Errorf("expected 0 after clear, got %d", v)
	}
}

func TestCheckCeiling_NotExceeded(t *testing.T) {
	st := ceilingBaseState()
	_ = SetCeiling(st, "prod", 5)
	if err := CheckCeiling(st, "prod"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCheckCeiling_Exceeded(t *testing.T) {
	st := ceilingBaseState()
	_ = SetCeiling(st, "prod", 2) // prod already has 2 records
	if err := CheckCeiling(st, "prod"); err == nil {
		t.Fatal("expected ceiling exceeded error")
	}
}

func TestCheckCeiling_NoCeilingAllowsAny(t *testing.T) {
	st := ceilingBaseState()
	if err := CheckCeiling(st, "prod"); err != nil {
		t.Errorf("unexpected error when no ceiling set: %v", err)
	}
}
