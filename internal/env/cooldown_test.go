package env

import (
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
	"github.com/patchwork-deploy/internal/state/record"
)

func cooldownBaseState() *state.State {
	st := state.NewInMemory()
	st.Add("staging", record.Record{Patch: "001-init.sql", AppliedAt: time.Now()})
	st.Add("staging", record.Record{Patch: "002-users.sql", AppliedAt: time.Now()})
	return st
}

func TestSetCooldown_Success(t *testing.T) {
	st := cooldownBaseState()
	err := SetCooldown(st, "staging", "001-init.sql", 30*time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok, err := GetCooldown(st, "staging", "001-init.sql")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected cooldown to be set")
	}
}

func TestSetCooldown_MissingEnvReturnsError(t *testing.T) {
	st := cooldownBaseState()
	err := SetCooldown(st, "production", "001-init.sql", 10*time.Minute)
	if err == nil {
		t.Fatal("expected error for missing environment")
	}
}

func TestSetCooldown_MissingPatchReturnsError(t *testing.T) {
	st := cooldownBaseState()
	err := SetCooldown(st, "staging", "999-missing.sql", 10*time.Minute)
	if err == nil {
		t.Fatal("expected error for missing patch")
	}
}

func TestSetCooldown_ZeroDurationReturnsError(t *testing.T) {
	st := cooldownBaseState()
	err := SetCooldown(st, "staging", "001-init.sql", 0)
	if err == nil {
		t.Fatal("expected error for zero duration")
	}
}

func TestInCooldown_TrueWhenActive(t *testing.T) {
	st := cooldownBaseState()
	_ = SetCooldown(st, "staging", "001-init.sql", 1*time.Hour)
	active, err := InCooldown(st, "staging", "001-init.sql")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !active {
		t.Fatal("expected cooldown to be active")
	}
}

func TestInCooldown_FalseWhenNoneSet(t *testing.T) {
	st := cooldownBaseState()
	active, err := InCooldown(st, "staging", "001-init.sql")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if active {
		t.Fatal("expected no cooldown to be active")
	}
}

func TestClearCooldown_RemovesEntry(t *testing.T) {
	st := cooldownBaseState()
	_ = SetCooldown(st, "staging", "001-init.sql", 1*time.Hour)
	err := ClearCooldown(st, "staging", "001-init.sql")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok, _ := GetCooldown(st, "staging", "001-init.sql")
	if ok {
		t.Fatal("expected cooldown to be cleared")
	}
}
