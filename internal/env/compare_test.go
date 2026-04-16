package env

import (
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func compareBaseState() *state.State {
	now := time.Now()
	st := &state.State{}
	st.Add(state.Record{Environment: "staging", Patch: "001-init", AppliedAt: now})
	st.Add(state.Record{Environment: "staging", Patch: "002-users", AppliedAt: now})
	st.Add(state.Record{Environment: "prod", Patch: "001-init", AppliedAt: now})
	st.Add(state.Record{Environment: "prod", Patch: "003-payments", AppliedAt: now})
	return st
}

func TestCompareEnvironments_ShowsDifferences(t *testing.T) {
	st := compareBaseState()
	res, err := CompareEnvironments(st, "staging", "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.InBoth) != 1 || res.InBoth[0] != "001-init" {
		t.Errorf("expected InBoth=[001-init], got %v", res.InBoth)
	}
	if len(res.OnlyInSource) != 1 || res.OnlyInSource[0] != "002-users" {
		t.Errorf("expected OnlyInSource=[002-users], got %v", res.OnlyInSource)
	}
	if len(res.OnlyInTarget) != 1 || res.OnlyInTarget[0] != "003-payments" {
		t.Errorf("expected OnlyInTarget=[003-payments], got %v", res.OnlyInTarget)
	}
}

func TestCompareEnvironments_SameEnvReturnsError(t *testing.T) {
	st := compareBaseState()
	_, err := CompareEnvironments(st, "staging", "staging")
	if err == nil {
		t.Fatal("expected error for same env")
	}
}

func TestCompareEnvironments_MissingSourceReturnsError(t *testing.T) {
	st := compareBaseState()
	_, err := CompareEnvironments(st, "ghost", "prod")
	if err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestCompareEnvironments_MissingTargetReturnsError(t *testing.T) {
	st := compareBaseState()
	_, err := CompareEnvironments(st, "staging", "ghost")
	if err == nil {
		t.Fatal("expected error for missing target")
	}
}

func TestCompareEnvironments_IdenticalEnvs(t *testing.T) {
	st := compareBaseState()
	// add dev with same patches as staging
	now := time.Now()
	st.Add(state.Record{Environment: "dev", Patch: "001-init", AppliedAt: now})
	st.Add(state.Record{Environment: "dev", Patch: "002-users", AppliedAt: now})
	res, err := CompareEnvironments(st, "staging", "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.OnlyInSource) != 0 || len(res.OnlyInTarget) != 0 {
		t.Errorf("expected no differences, got src=%v tgt=%v", res.OnlyInSource, res.OnlyInTarget)
	}
	if len(res.InBoth) != 2 {
		t.Errorf("expected 2 in both, got %d", len(res.InBoth))
	}
}
