package env

import (
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func pauseBaseState() *state.State {
	st := state.New()
	st.Add(state.Record{Environment: "staging", Patch: "001-init.sql", AppliedAt: time.Now()})
	st.Add(state.Record{Environment: "prod", Patch: "001-init.sql", AppliedAt: time.Now()})
	return st
}

func TestPause_MarksEnvironment(t *testing.T) {
	st := pauseBaseState()
	if err := Pause(st, "staging"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !IsPaused(st, "staging") {
		t.Error("expected staging to be paused")
	}
}

func TestPause_IdempotentWhenAlreadyPaused(t *testing.T) {
	st := pauseBaseState()
	_ = Pause(st, "staging")
	if err := Pause(st, "staging"); err != nil {
		t.Fatalf("unexpected error on second pause: %v", err)
	}
	count := 0
	for _, r := range st.ForEnvironment("staging") {
		if r.Patch == pauseSentinel {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected 1 sentinel, got %d", count)
	}
}

func TestUnpause_RemovesSentinel(t *testing.T) {
	st := pauseBaseState()
	_ = Pause(st, "staging")
	if err := Unpause(st, "staging"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if IsPaused(st, "staging") {
		t.Error("expected staging to be unpaused")
	}
}

func TestPause_MissingEnvReturnsError(t *testing.T) {
	st := pauseBaseState()
	if err := Pause(st, "unknown"); err == nil {
		t.Error("expected error for missing environment")
	}
}

func TestPause_DoesNotAffectOtherEnvs(t *testing.T) {
	st := pauseBaseState()
	_ = Pause(st, "staging")
	if IsPaused(st, "prod") {
		t.Error("prod should not be paused")
	}
}
