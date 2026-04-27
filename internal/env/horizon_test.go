package env_test

import (
	"testing"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
)

func horizonBaseState() *state.State {
	st := state.NewInMemory()
	st.AddRecord("prod", state.Record{Patch: "001-init.sql"})
	st.AddRecord("prod", state.Record{Patch: "002-add-users.sql"})
	return st
}

func TestSetHorizon_Success(t *testing.T) {
	st := horizonBaseState()
	if err := env.SetHorizon(st, "prod", "001-init.sql", 5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	limit, ok := env.GetHorizon(st, "prod", "001-init.sql")
	if !ok || limit != 5 {
		t.Errorf("expected limit 5, got %d (ok=%v)", limit, ok)
	}
}

func TestSetHorizon_ZeroLimitReturnsError(t *testing.T) {
	st := horizonBaseState()
	if err := env.SetHorizon(st, "prod", "001-init.sql", 0); err == nil {
		t.Error("expected error for zero limit")
	}
}

func TestSetHorizon_MissingEnvReturnsError(t *testing.T) {
	st := horizonBaseState()
	if err := env.SetHorizon(st, "staging", "001-init.sql", 3); err == nil {
		t.Error("expected error for missing environment")
	}
}

func TestSetHorizon_MissingPatchReturnsError(t *testing.T) {
	st := horizonBaseState()
	if err := env.SetHorizon(st, "prod", "999-unknown.sql", 3); err == nil {
		t.Error("expected error for missing patch")
	}
}

func TestGetHorizon_NoneSet(t *testing.T) {
	st := horizonBaseState()
	_, ok := env.GetHorizon(st, "prod", "001-init.sql")
	if ok {
		t.Error("expected no horizon to be set")
	}
}

func TestClearHorizon_RemovesEntry(t *testing.T) {
	st := horizonBaseState()
	_ = env.SetHorizon(st, "prod", "001-init.sql", 10)
	env.ClearHorizon(st, "prod", "001-init.sql")
	_, ok := env.GetHorizon(st, "prod", "001-init.sql")
	if ok {
		t.Error("expected horizon to be cleared")
	}
}

func TestListHorizons_ReturnsAllEntries(t *testing.T) {
	st := horizonBaseState()
	_ = env.SetHorizon(st, "prod", "001-init.sql", 3)
	_ = env.SetHorizon(st, "prod", "002-add-users.sql", 7)
	horizons := env.ListHorizons(st, "prod")
	if len(horizons) != 2 {
		t.Fatalf("expected 2 horizons, got %d", len(horizons))
	}
	if horizons["001-init.sql"] != 3 {
		t.Errorf("expected 3 for 001-init.sql, got %d", horizons["001-init.sql"])
	}
	if horizons["002-add-users.sql"] != 7 {
		t.Errorf("expected 7 for 002-add-users.sql, got %d", horizons["002-add-users.sql"])
	}
}
