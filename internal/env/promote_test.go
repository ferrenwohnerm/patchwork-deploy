package env_test

import (
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
)

func baseState() *state.State {
	st := &state.State{}
	now := time.Now().UTC()
	for _, p := range []string{"001-init.sql", "002-seed.sql", "003-index.sql"} {
		st.Add(state.Record{Environment: "staging", Patch: p, AppliedAt: now})
	}
	// 001 already in production
	st.Add(state.Record{Environment: "production", Patch: "001-init.sql", AppliedAt: now})
	return st
}

func TestPromote_CopiesNewPatches(t *testing.T) {
	st := baseState()
	res, err := env.Promote(st, "staging", "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Promoted) != 2 {
		t.Errorf("expected 2 promoted, got %d", len(res.Promoted))
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(res.Skipped))
	}
}

func TestPromote_SameEnvReturnsError(t *testing.T) {
	st := baseState()
	_, err := env.Promote(st, "staging", "staging")
	if err == nil {
		t.Fatal("expected error for same source/target env")
	}
}

func TestPromote_MissingSourceReturnsError(t *testing.T) {
	st := baseState()
	_, err := env.Promote(st, "unknown", "production")
	if err == nil {
		t.Fatal("expected error for unknown source environment")
	}
}

func TestPromote_RecordsAreWrittenToState(t *testing.T) {
	st := baseState()
	_, err := env.Promote(st, "staging", "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	prod := st.ForEnvironment("production")
	if len(prod) != 3 {
		t.Errorf("expected 3 production records after promote, got %d", len(prod))
	}
}

func TestPromote_EmptyTargetPromotesAll(t *testing.T) {
	st := baseState()
	res, err := env.Promote(st, "staging", "canary")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Promoted) != 3 {
		t.Errorf("expected all 3 patches promoted, got %d", len(res.Promoted))
	}
	if len(res.Skipped) != 0 {
		t.Errorf("expected 0 skipped, got %d", len(res.Skipped))
	}
}
