package env

import (
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func statusBaseState() *state.State {
	now := time.Now()
	st := state.New()
	st.Add(state.Record{Environment: "prod", Patch: "001-init", AppliedAt: now.Add(-2 * time.Hour), Tags: []string{"baseline"}})
	st.Add(state.Record{Environment: "prod", Patch: "002-add-index", AppliedAt: now.Add(-1 * time.Hour)})
	st.Add(state.Record{Environment: "staging", Patch: "001-init", AppliedAt: now})
	return st
}

func TestStatus_ReturnsCorrectTotal(t *testing.T) {
	st := statusBaseState()
	s, err := Status(st, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Total != 2 {
		t.Errorf("expected 2, got %d", s.Total)
	}
}

func TestStatus_ReturnsLatestPatch(t *testing.T) {
	st := statusBaseState()
	s, err := Status(st, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Latest != "002-add-index" {
		t.Errorf("expected 002-add-index, got %s", s.Latest)
	}
}

func TestStatus_CollectsTags(t *testing.T) {
	st := statusBaseState()
	s, err := Status(st, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Tags) != 1 || s.Tags[0] != "baseline" {
		t.Errorf("unexpected tags: %v", s.Tags)
	}
}

func TestStatus_MissingEnvReturnsError(t *testing.T) {
	st := statusBaseState()
	_, err := Status(st, "ghost")
	if err == nil {
		t.Fatal("expected error for missing env")
	}
}

func TestStatusAll_ReturnsAllEnvs(t *testing.T) {
	st := statusBaseState()
	all := StatusAll(st)
	if len(all) != 2 {
		t.Errorf("expected 2 envs, got %d", len(all))
	}
	if all[0].Environment != "prod" {
		t.Errorf("expected prod first, got %s", all[0].Environment)
	}
}

func TestStatus_String_NoPatches(t *testing.T) {
	s := EnvStatus{Environment: "empty", Total: 0}
	if s.String() != "empty: no patches applied" {
		t.Errorf("unexpected string: %s", s.String())
	}
}
