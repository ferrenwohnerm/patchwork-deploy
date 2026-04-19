package env

import (
	"testing"

	"github.com/patchwork-deploy/internal/state"
)

func throttleBaseState() *state.State {
	st := state.NewInMemory()
	st.AddRecord("staging", state.Record{Patch: "001-init.sql"})
	st.AddRecord("prod", state.Record{Patch: "001-init.sql"})
	return st
}

func TestSetThrottle_Success(t *testing.T) {
	st := throttleBaseState()
	if err := SetThrottle(st, "staging", 5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	n, err := GetThrottle(st, "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 5 {
		t.Errorf("expected 5, got %d", n)
	}
}

func TestSetThrottle_ZeroLimitReturnsError(t *testing.T) {
	st := throttleBaseState()
	if err := SetThrottle(st, "staging", 0); err == nil {
		t.Fatal("expected error for zero limit")
	}
}

func TestSetThrottle_MissingEnvReturnsError(t *testing.T) {
	st := throttleBaseState()
	if err := SetThrottle(st, "ghost", 3); err == nil {
		t.Fatal("expected error for missing env")
	}
}

func TestGetThrottle_NoneSet(t *testing.T) {
	st := throttleBaseState()
	n, err := GetThrottle(st, "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0, got %d", n)
	}
}

func TestClearThrottle_RemovesLimit(t *testing.T) {
	st := throttleBaseState()
	_ = SetThrottle(st, "staging", 10)
	if err := ClearThrottle(st, "staging"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	n, _ := GetThrottle(st, "staging")
	if n != 0 {
		t.Errorf("expected 0 after clear, got %d", n)
	}
}

func TestListThrottles_ReturnsAllEntries(t *testing.T) {
	st := throttleBaseState()
	_ = SetThrottle(st, "staging", 3)
	_ = SetThrottle(st, "prod", 7)
	m := ListThrottles(st)
	if m["staging"] != 3 {
		t.Errorf("expected staging=3, got %d", m["staging"])
	}
	if m["prod"] != 7 {
		t.Errorf("expected prod=7, got %d", m["prod"])
	}
}
