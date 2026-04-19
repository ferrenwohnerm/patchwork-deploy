package env

import (
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func rolloutBaseState() *state.State {
	st := state.New()
	now := time.Now()
	st.Add(state.Record{Environment: "staging", Patch: "001-init", AppliedAt: now})
	st.Add(state.Record{Environment: "staging", Patch: "002-feature", AppliedAt: now})
	return st
}

func TestSetRollout_Success(t *testing.T) {
	st := rolloutBaseState()
	if err := SetRollout(st, "staging", "001-init", RolloutCanary); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cfg, ok := GetRollout(st, "staging", "001-init")
	if !ok {
		t.Fatal("expected rollout to be set")
	}
	if cfg.Strategy != RolloutCanary {
		t.Errorf("expected canary, got %s", cfg.Strategy)
	}
}

func TestSetRollout_MissingEnvReturnsError(t *testing.T) {
	st := rolloutBaseState()
	if err := SetRollout(st, "prod", "001-init", RolloutCanary); err == nil {
		t.Fatal("expected error for missing env")
	}
}

func TestSetRollout_MissingPatchReturnsError(t *testing.T) {
	st := rolloutBaseState()
	if err := SetRollout(st, "staging", "999-missing", RolloutCanary); err == nil {
		t.Fatal("expected error for missing patch")
	}
}

func TestSetRollout_InvalidStrategyReturnsError(t *testing.T) {
	st := rolloutBaseState()
	if err := SetRollout(st, "staging", "001-init", RolloutStrategy("rolling")); err == nil {
		t.Fatal("expected error for invalid strategy")
	}
}

func TestClearRollout_RemovesEntry(t *testing.T) {
	st := rolloutBaseState()
	_ = SetRollout(st, "staging", "001-init", RolloutBlueGreen)
	ClearRollout(st, "staging", "001-init")
	if _, ok := GetRollout(st, "staging", "001-init"); ok {
		t.Fatal("expected rollout to be cleared")
	}
}

func TestListRollouts_ReturnsAll(t *testing.T) {
	st := rolloutBaseState()
	_ = SetRollout(st, "staging", "001-init", RolloutCanary)
	_ = SetRollout(st, "staging", "002-feature", RolloutImmediate)
	list := ListRollouts(st, "staging")
	if len(list) != 2 {
		t.Errorf("expected 2 rollouts, got %d", len(list))
	}
}
