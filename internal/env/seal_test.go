package env_test

import (
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
)

func sealBaseState() *state.State {
	st := state.New()
	st.Add(state.Record{Environment: "prod", Patch: "001-init.sql", AppliedAt: time.Now()})
	st.Add(state.Record{Environment: "staging", Patch: "001-init.sql", AppliedAt: time.Now()})
	return st
}

func TestSeal_MarksEnvironment(t *testing.T) {
	st := sealBaseState()
	if err := env.Seal(st, "prod", "release freeze"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !env.IsSealed(st, "prod") {
		t.Fatal("expected prod to be sealed")
	}
}

func TestSeal_StoresReason(t *testing.T) {
	st := sealBaseState()
	_ = env.Seal(st, "prod", "compliance hold")
	if got := env.SealReason(st, "prod"); got != "compliance hold" {
		t.Fatalf("expected 'compliance hold', got %q", got)
	}
}

func TestSeal_IdempotentWhenAlreadySealed(t *testing.T) {
	st := sealBaseState()
	_ = env.Seal(st, "prod", "first")
	if err := env.Seal(st, "prod", "second"); err != nil {
		t.Fatalf("unexpected error on second seal: %v", err)
	}
	count := 0
	for _, r := range st.ForEnvironment("prod") {
		if r.Patch == "__sealed__" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("expected exactly 1 seal sentinel, got %d", count)
	}
}

func TestUnseal_RemovesSentinel(t *testing.T) {
	st := sealBaseState()
	_ = env.Seal(st, "prod", "")
	if err := env.Unseal(st, "prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env.IsSealed(st, "prod") {
		t.Fatal("expected prod to be unsealed")
	}
}

func TestSeal_MissingEnvReturnsError(t *testing.T) {
	st := sealBaseState()
	if err := env.Seal(st, "unknown", ""); err == nil {
		t.Fatal("expected error for missing env")
	}
}

func TestIsSealed_FalseForOtherEnvs(t *testing.T) {
	st := sealBaseState()
	_ = env.Seal(st, "prod", "")
	if env.IsSealed(st, "staging") {
		t.Fatal("staging should not be sealed")
	}
}
