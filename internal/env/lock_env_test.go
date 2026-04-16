package env

import (
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func lockBaseState() *state.State {
	st := state.New()
	st.Add(state.Record{Environment: "staging", Patch: "001-init", AppliedAt: time.Now()})
	return st
}

func TestLockEnv_MarksEnvironment(t *testing.T) {
	st := lockBaseState()
	if err := LockEnvironment(st, "staging"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !IsEnvLocked(st, "staging") {
		t.Error("expected environment to be locked")
	}
}

func TestLockEnv_IdempotentWhenAlreadyLocked(t *testing.T) {
	st := lockBaseState()
	_ = LockEnvironment(st, "staging")
	if err := LockEnvironment(st, "staging"); err != nil {
		t.Fatalf("expected no error on second lock, got: %v", err)
	}
}

func TestUnlockEnv_RemovesSentinel(t *testing.T) {
	st := lockBaseState()
	_ = LockEnvironment(st, "staging")
	if err := UnlockEnvironment(st, "staging"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if IsEnvLocked(st, "staging") {
		t.Error("expected environment to be unlocked")
	}
}

func TestUnlockEnv_IdempotentWhenNotLocked(t *testing.T) {
	st := lockBaseState()
	if err := UnlockEnvironment(st, "staging"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLockEnv_MissingEnvReturnsError(t *testing.T) {
	st := state.New()
	if err := LockEnvironment(st, "ghost"); err == nil {
		t.Error("expected error for missing environment")
	}
}
