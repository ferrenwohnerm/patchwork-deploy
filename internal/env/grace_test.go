package env

import (
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func graceBaseState() *state.State {
	st := state.NewInMemory()
	st.Add("prod", state.Record{Patch: "001-init", AppliedAt: time.Now()})
	st.Add("prod", state.Record{Patch: "002-users", AppliedAt: time.Now()})
	return st
}

func TestSetGracePeriod_Success(t *testing.T) {
	st := graceBaseState()
	err := SetGracePeriod(st, "prod", "001-init", 30*time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok, err := GetGracePeriod(st, "prod", "001-init")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected grace period to be set")
	}
}

func TestSetGracePeriod_MissingEnvReturnsError(t *testing.T) {
	st := graceBaseState()
	err := SetGracePeriod(st, "staging", "001-init", time.Minute)
	if err == nil {
		t.Fatal("expected error for missing env")
	}
}

func TestSetGracePeriod_MissingPatchReturnsError(t *testing.T) {
	st := graceBaseState()
	err := SetGracePeriod(st, "prod", "999-missing", time.Minute)
	if err == nil {
		t.Fatal("expected error for missing patch")
	}
}

func TestSetGracePeriod_ZeroDurationReturnsError(t *testing.T) {
	st := graceBaseState()
	err := SetGracePeriod(st, "prod", "001-init", 0)
	if err == nil {
		t.Fatal("expected error for zero duration")
	}
}

func TestInGracePeriod_TrueWhenActive(t *testing.T) {
	st := graceBaseState()
	_ = SetGracePeriod(st, "prod", "001-init", time.Hour)
	in, err := InGracePeriod(st, "prod", "001-init")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !in {
		t.Fatal("expected patch to be in grace period")
	}
}

func TestInGracePeriod_FalseWhenNotSet(t *testing.T) {
	st := graceBaseState()
	in, err := InGracePeriod(st, "prod", "002-users")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if in {
		t.Fatal("expected patch to not be in grace period")
	}
}

func TestClearGracePeriod_RemovesSentinel(t *testing.T) {
	st := graceBaseState()
	_ = SetGracePeriod(st, "prod", "001-init", time.Hour)
	ClearGracePeriod(st, "prod", "001-init")
	_, ok, _ := GetGracePeriod(st, "prod", "001-init")
	if ok {
		t.Fatal("expected grace period to be cleared")
	}
}
