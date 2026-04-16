package env

import (
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func resetBaseState() *state.State {
	st := state.New()
	now := time.Now()
	st.Add(state.Record{Environment: "staging", Patch: "001-init", AppliedAt: now})
	st.Add(state.Record{Environment: "staging", Patch: "002-schema", AppliedAt: now})
	st.Add(state.Record{Environment: "prod", Patch: "001-init", AppliedAt: now})
	return st
}

func TestReset_RemovesAllRecordsForEnv(t *testing.T) {
	st := resetBaseState()
	res, err := Reset(st, "staging", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Removed != 2 {
		t.Errorf("expected 2 removed, got %d", res.Removed)
	}
	if len(st.ForEnvironment("staging")) != 0 {
		t.Error("expected staging records to be empty")
	}
	if len(st.ForEnvironment("prod")) != 1 {
		t.Error("prod records should be untouched")
	}
}

func TestReset_DryRunDoesNotMutate(t *testing.T) {
	st := resetBaseState()
	res, err := Reset(st, "staging", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Removed != 2 {
		t.Errorf("expected 2 removed in dry run, got %d", res.Removed)
	}
	if len(st.ForEnvironment("staging")) != 2 {
		t.Error("dry run should not mutate state")
	}
}

func TestReset_MissingEnvReturnsError(t *testing.T) {
	st := resetBaseState()
	_, err := Reset(st, "unknown", false)
	if err == nil {
		t.Error("expected error for unknown environment")
	}
}
