package state

import (
	"testing"
	"time"
)

func baseTime() time.Time {
	return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
}

func stateWithRecords() *State {
	t0 := baseTime()
	s := &State{}
	s.Records = []Record{
		{Environment: "prod", Patch: "001-init", AppliedAt: t0},
		{Environment: "prod", Patch: "002-schema", AppliedAt: t0.Add(time.Hour)},
		{Environment: "prod", Patch: "003-seed", AppliedAt: t0.Add(2 * time.Hour)},
		{Environment: "staging", Patch: "001-init", AppliedAt: t0},
	}
	return s
}

func TestBuildRollbackPlan_FullHistory(t *testing.T) {
	s := stateWithRecords()
	plan, err := BuildRollbackPlan(s, "prod", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if plan.Environment != "prod" {
		t.Errorf("expected env prod, got %s", plan.Environment)
	}
	if len(plan.Patches) != 3 {
		t.Fatalf("expected 3 patches, got %d", len(plan.Patches))
	}
	// Should be in reverse order.
	if plan.Patches[0] != "003-seed" || plan.Patches[2] != "001-init" {
		t.Errorf("unexpected order: %v", plan.Patches)
	}
}

func TestBuildRollbackPlan_FromPatch(t *testing.T) {
	s := stateWithRecords()
	plan, err := BuildRollbackPlan(s, "prod", "002-schema")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(plan.Patches) != 2 {
		t.Fatalf("expected 2 patches, got %d", len(plan.Patches))
	}
	if plan.Patches[0] != "003-seed" || plan.Patches[1] != "002-schema" {
		t.Errorf("unexpected order: %v", plan.Patches)
	}
}

func TestBuildRollbackPlan_UnknownEnvironment(t *testing.T) {
	s := stateWithRecords()
	_, err := BuildRollbackPlan(s, "dev", "")
	if err == nil {
		t.Fatal("expected error for unknown environment")
	}
}

func TestBuildRollbackPlan_UnknownPatch(t *testing.T) {
	s := stateWithRecords()
	_, err := BuildRollbackPlan(s, "prod", "999-missing")
	if err == nil {
		t.Fatal("expected error for unknown patch")
	}
}
