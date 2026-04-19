package env

import (
	"testing"

	"github.com/patchwork-deploy/internal/state"
)

func weightBaseState() *state.State {
	st := state.New()
	st.Add(state.Record{Environment: "staging", Patch: "001-init"})
	st.Add(state.Record{Environment: "staging", Patch: "002-schema"})
	st.Add(state.Record{Environment: "prod", Patch: "001-init"})
	return st
}

func TestSetWeight_Success(t *testing.T) {
	st := weightBaseState()
	if err := SetWeight(st, "staging", "001-init", 10); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	w, ok := GetWeight(st, "staging", "001-init")
	if !ok || w != 10 {
		t.Errorf("expected weight 10, got %d (ok=%v)", w, ok)
	}
}

func TestSetWeight_NegativeReturnsError(t *testing.T) {
	st := weightBaseState()
	if err := SetWeight(st, "staging", "001-init", -1); err == nil {
		t.Error("expected error for negative weight")
	}
}

func TestSetWeight_MissingEnvReturnsError(t *testing.T) {
	st := weightBaseState()
	if err := SetWeight(st, "unknown", "001-init", 5); err == nil {
		t.Error("expected error for missing environment")
	}
}

func TestSetWeight_MissingPatchReturnsError(t *testing.T) {
	st := weightBaseState()
	if err := SetWeight(st, "staging", "999-nope", 5); err == nil {
		t.Error("expected error for missing patch")
	}
}

func TestGetWeight_NoneSet(t *testing.T) {
	st := weightBaseState()
	_, ok := GetWeight(st, "staging", "002-schema")
	if ok {
		t.Error("expected no weight to be set")
	}
}

func TestClearWeight_RemovesEntry(t *testing.T) {
	st := weightBaseState()
	_ = SetWeight(st, "staging", "001-init", 7)
	ClearWeight(st, "staging", "001-init")
	_, ok := GetWeight(st, "staging", "001-init")
	if ok {
		t.Error("expected weight to be cleared")
	}
}

func TestListWeights_ReturnsAllForEnv(t *testing.T) {
	st := weightBaseState()
	_ = SetWeight(st, "staging", "001-init", 3)
	_ = SetWeight(st, "staging", "002-schema", 7)
	_ = SetWeight(st, "prod", "001-init", 99)
	weights := ListWeights(st, "staging")
	if len(weights) != 2 {
		t.Fatalf("expected 2 weights for staging, got %d", len(weights))
	}
	if weights["001-init"] != 3 || weights["002-schema"] != 7 {
		t.Errorf("unexpected weights: %v", weights)
	}
}
