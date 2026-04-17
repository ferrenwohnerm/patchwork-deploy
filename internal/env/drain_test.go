package env

import (
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func drainBaseState() *state.State {
	st := state.New()
	st.Add(state.Record{Environment: "staging", Patch: "001-init.sql", AppliedAt: time.Now()})
	st.Add(state.Record{Environment: "__scheduled__staging", Patch: "002-add-index.sql", AppliedAt: time.Now()})
	st.Add(state.Record{Environment: "__scheduled__staging", Patch: "003-cleanup.sql", AppliedAt: time.Now()})
	return st
}

func TestDrain_RemovesScheduledEntries(t *testing.T) {
	st := drainBaseState()
	res, err := Drain(st, "staging", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 2 {
		t.Fatalf("expected 2 removed, got %d", len(res.Removed))
	}
	if len(st.ForEnvironment("__scheduled__staging")) != 0 {
		t.Error("expected scheduled entries to be cleared")
	}
}

func TestDrain_DryRunDoesNotMutate(t *testing.T) {
	st := drainBaseState()
	res, err := Drain(st, "staging", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 2 {
		t.Fatalf("expected 2 in result, got %d", len(res.Removed))
	}
	if len(st.ForEnvironment("__scheduled__staging")) != 2 {
		t.Error("dry run should not remove scheduled entries")
	}
}

func TestDrain_MissingEnvReturnsError(t *testing.T) {
	st := drainBaseState()
	_, err := Drain(st, "production", false)
	if err == nil {
		t.Fatal("expected error for missing env")
	}
}

func TestIsDrained_TrueAfterDrain(t *testing.T) {
	st := drainBaseState()
	if IsDrained(st, "staging") {
		t.Fatal("should not be drained initially")
	}
	Drain(st, "staging", false)
	if !IsDrained(st, "staging") {
		t.Fatal("should be drained after drain call")
	}
}

func TestUndrain_ClearsSentinel(t *testing.T) {
	st := drainBaseState()
	Drain(st, "staging", false)
	if err := Undrain(st, "staging"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if IsDrained(st, "staging") {
		t.Error("expected drain sentinel to be cleared")
	}
}
