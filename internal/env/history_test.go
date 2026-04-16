package env

import (
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func historyBaseState() *state.State {
	st := state.New()
	st.Add(state.Record{Environment: "prod", Patch: "001-init", AppliedAt: epoch})
	st.Add(state.Record{Environment: "prod", Patch: "002-feature", AppliedAt: epoch.Add(time.Hour)})
	st.Add(state.Record{Environment: "prod", Patch: "003-fix", AppliedAt: epoch.Add(2 * time.Hour)})
	st.Add(state.Record{Environment: "staging", Patch: "001-init", AppliedAt: epoch})
	return st
}

func TestHistory_ReturnsSortedRecords(t *testing.T) {
	st := historyBaseState()
	recs, err := History(st, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recs) != 3 {
		t.Fatalf("expected 3 records, got %d", len(recs))
	}
	if recs[0].Patch != "001-init" || recs[2].Patch != "003-fix" {
		t.Errorf("unexpected order: %v", recs)
	}
}

func TestHistory_MissingEnvReturnsError(t *testing.T) {
	st := historyBaseState()
	_, err := History(st, "ghost")
	if err == nil {
		t.Fatal("expected error for missing env")
	}
}

func TestHistorySince_FiltersCorrectly(t *testing.T) {
	st := historyBaseState()
	cutoff := epoch.Add(30 * time.Minute)
	recs, err := HistorySince(st, "prod", cutoff)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recs) != 2 {
		t.Fatalf("expected 2 records after cutoff, got %d", len(recs))
	}
	if recs[0].Patch != "002-feature" {
		t.Errorf("expected 002-feature first, got %s", recs[0].Patch)
	}
}

func TestHistorySince_AllIncluded(t *testing.T) {
	st := historyBaseState()
	recs, err := HistorySince(st, "prod", epoch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recs) != 3 {
		t.Errorf("expected all 3 records, got %d", len(recs))
	}
}
