package env

import (
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func freezeBaseState() *state.State {
	st := &state.State{}
	st.Add(state.Record{Environment: "staging", Patch: "001-init", AppliedAt: time.Now()})
	st.Add(state.Record{Environment: "prod", Patch: "001-init", AppliedAt: time.Now()})
	return st
}

func TestFreeze_MarksEnvironment(t *testing.T) {
	st := freezeBaseState()
	res, err := Freeze(st, "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Frozen {
		t.Error("expected Frozen=true")
	}
	if !IsFrozen(st, "staging") {
		t.Error("expected staging to be frozen")
	}
}

func TestFreeze_IdempotentWhenAlreadyFrozen(t *testing.T) {
	st := freezeBaseState()
	_, _ = Freeze(st, "staging")
	_, err := Freeze(st, "staging")
	if err != nil {
		t.Fatalf("unexpected error on second freeze: %v", err)
	}
	count := 0
	for _, r := range st.ForEnvironment("staging") {
		if r.Patch == freezePrefix {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected 1 freeze sentinel, got %d", count)
	}
}

func TestUnfreeze_RemovesSentinel(t *testing.T) {
	st := freezeBaseState()
	_, _ = Freeze(st, "staging")
	res, err := Unfreeze(st, "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Frozen {
		t.Error("expected Frozen=false")
	}
	if IsFrozen(st, "staging") {
		t.Error("expected staging to be unfrozen")
	}
}

func TestFreeze_MissingEnvReturnsError(t *testing.T) {
	st := freezeBaseState()
	_, err := Freeze(st, "nonexistent")
	if err == nil {
		t.Error("expected error for missing environment")
	}
}

func TestIsFrozen_FalseByDefault(t *testing.T) {
	st := freezeBaseState()
	if IsFrozen(st, "prod") {
		t.Error("expected prod to not be frozen")
	}
}
