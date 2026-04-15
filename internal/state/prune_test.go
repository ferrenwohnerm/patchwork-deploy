package state

import (
	"testing"
	"time"
)

var (
	now   = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	old   = now.Add(-48 * time.Hour)
	cutoff = now.Add(-24 * time.Hour)
)

func stateWithAgedRecords() *State {
	return &State{
		Records: []Record{
			{Environment: "prod", Patch: "001-init.sql", AppliedAt: old},
			{Environment: "prod", Patch: "002-add-index.sql", AppliedAt: now},
			{Environment: "staging", Patch: "001-init.sql", AppliedAt: old},
		},
	}
}

func TestPruneByAge_RemovesOldRecords(t *testing.T) {
	st := stateWithAgedRecords()
	res := PruneByAge(st, PruneOptions{OlderThan: cutoff})

	if len(res.Removed) != 2 {
		t.Fatalf("expected 2 removed, got %d", len(res.Removed))
	}
	if len(res.Retained) != 1 {
		t.Fatalf("expected 1 retained, got %d", len(res.Retained))
	}
	if len(st.Records) != 1 {
		t.Fatalf("expected state to have 1 record after prune, got %d", len(st.Records))
	}
}

func TestPruneByAge_DryRunDoesNotMutate(t *testing.T) {
	st := stateWithAgedRecords()
	res := PruneByAge(st, PruneOptions{OlderThan: cutoff, DryRun: true})

	if len(res.Removed) != 2 {
		t.Fatalf("expected 2 would-be removed, got %d", len(res.Removed))
	}
	if len(st.Records) != 3 {
		t.Fatalf("state should be unchanged in dry run, got %d records", len(st.Records))
	}
}

func TestPruneByAge_FiltersByEnvironment(t *testing.T) {
	st := stateWithAgedRecords()
	res := PruneByAge(st, PruneOptions{OlderThan: cutoff, Environment: "prod"})

	if len(res.Removed) != 1 {
		t.Fatalf("expected 1 removed for prod, got %d", len(res.Removed))
	}
	if res.Removed[0].Environment != "prod" {
		t.Errorf("expected removed record to be prod, got %s", res.Removed[0].Environment)
	}
	if len(st.Records) != 2 {
		t.Fatalf("expected 2 remaining records, got %d", len(st.Records))
	}
}

func TestPruneByAge_NoMatchLeavesStateIntact(t *testing.T) {
	st := stateWithAgedRecords()
	futureCutoff := old.Add(-1 * time.Hour)
	res := PruneByAge(st, PruneOptions{OlderThan: futureCutoff})

	if len(res.Removed) != 0 {
		t.Fatalf("expected nothing removed, got %d", len(res.Removed))
	}
	if len(st.Records) != 3 {
		t.Fatalf("expected state unchanged, got %d records", len(st.Records))
	}
}
