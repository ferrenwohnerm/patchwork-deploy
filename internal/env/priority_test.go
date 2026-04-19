package env

import (
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func priorityBaseState() *state.State {
	st := state.New()
	st.Add(state.Record{Environment: "prod", Patch: "001-init", AppliedAt: time.Now()})
	st.Add(state.Record{Environment: "prod", Patch: "002-schema", AppliedAt: time.Now()})
	return st
}

func TestSetPriority_Success(t *testing.T) {
	st := priorityBaseState()
	if err := SetPriority(st, "prod", "001-init", 10); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	n, ok := GetPriority(st, "prod", "001-init")
	if !ok || n != 10 {
		t.Fatalf("expected priority 10, got %d ok=%v", n, ok)
	}
}

func TestSetPriority_NegativeReturnsError(t *testing.T) {
	st := priorityBaseState()
	if err := SetPriority(st, "prod", "001-init", -1); err == nil {
		t.Fatal("expected error for negative priority")
	}
}

func TestSetPriority_MissingEnvReturnsError(t *testing.T) {
	st := priorityBaseState()
	if err := SetPriority(st, "staging", "001-init", 5); err == nil {
		t.Fatal("expected error for missing env")
	}
}

func TestSetPriority_MissingPatchReturnsError(t *testing.T) {
	st := priorityBaseState()
	if err := SetPriority(st, "prod", "999-nope", 5); err == nil {
		t.Fatal("expected error for missing patch")
	}
}

func TestGetPriority_NoneSet(t *testing.T) {
	st := priorityBaseState()
	_, ok := GetPriority(st, "prod", "001-init")
	if ok {
		t.Fatal("expected no priority to be set")
	}
}

func TestClearPriority_RemovesSetting(t *testing.T) {
	st := priorityBaseState()
	_ = SetPriority(st, "prod", "001-init", 7)
	ClearPriority(st, "prod", "001-init")
	_, ok := GetPriority(st, "prod", "001-init")
	if ok {
		t.Fatal("expected priority to be cleared")
	}
}

func TestListPriorities_ReturnsAll(t *testing.T) {
	st := priorityBaseState()
	_ = SetPriority(st, "prod", "001-init", 1)
	_ = SetPriority(st, "prod", "002-schema", 5)
	m := ListPriorities(st, "prod")
	if len(m) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(m))
	}
	if m["001-init"] != 1 || m["002-schema"] != 5 {
		t.Fatalf("unexpected priorities: %v", m)
	}
}
