package env

import (
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func timeoutBaseState() *state.State {
	st := state.New()
	st.Add(state.Record{Environment: "staging", Patch: "001-init"})
	st.Add(state.Record{Environment: "staging", Patch: "002-schema"})
	return st
}

func TestSetTimeout_Success(t *testing.T) {
	st := timeoutBaseState()
	err := SetTimeout(st, "staging", "001-init", 30*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	d, err := GetTimeout(st, "staging", "001-init")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d != 30*time.Second {
		t.Errorf("expected 30s, got %v", d)
	}
}

func TestSetTimeout_ZeroDurationReturnsError(t *testing.T) {
	st := timeoutBaseState()
	if err := SetTimeout(st, "staging", "001-init", 0); err == nil {
		t.Error("expected error for zero duration")
	}
}

func TestSetTimeout_MissingEnvReturnsError(t *testing.T) {
	st := timeoutBaseState()
	if err := SetTimeout(st, "prod", "001-init", time.Minute); err == nil {
		t.Error("expected error for missing environment")
	}
}

func TestSetTimeout_MissingPatchReturnsError(t *testing.T) {
	st := timeoutBaseState()
	if err := SetTimeout(st, "staging", "999-nope", time.Minute); err == nil {
		t.Error("expected error for missing patch")
	}
}

func TestGetTimeout_NoneSet(t *testing.T) {
	st := timeoutBaseState()
	d, err := GetTimeout(st, "staging", "002-schema")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d != 0 {
		t.Errorf("expected 0, got %v", d)
	}
}

func TestClearTimeout_RemovesEntry(t *testing.T) {
	st := timeoutBaseState()
	_ = SetTimeout(st, "staging", "001-init", time.Minute)
	ClearTimeout(st, "staging", "001-init")
	d, _ := GetTimeout(st, "staging", "001-init")
	if d != 0 {
		t.Errorf("expected 0 after clear, got %v", d)
	}
}

func TestListTimeouts_ReturnsAll(t *testing.T) {
	st := timeoutBaseState()
	_ = SetTimeout(st, "staging", "001-init", 10*time.Second)
	_ = SetTimeout(st, "staging", "002-schema", 20*time.Second)
	m := ListTimeouts(st, "staging")
	if len(m) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(m))
	}
	if m["001-init"] != 10*time.Second {
		t.Errorf("wrong duration for 001-init: %v", m["001-init"])
	}
	if m["002-schema"] != 20*time.Second {
		t.Errorf("wrong duration for 002-schema: %v", m["002-schema"])
	}
}
