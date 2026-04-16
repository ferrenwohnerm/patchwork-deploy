package env

import (
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func archiveBaseState() *state.State {
	st := state.New()
	now := time.Now().UTC()
	st.Add(state.Record{Environment: "staging", Patch: "001-init.sql", AppliedAt: now})
	st.Add(state.Record{Environment: "staging", Patch: "002-users.sql", AppliedAt: now})
	st.Add(state.Record{Environment: "prod", Patch: "001-init.sql", AppliedAt: now})
	return st
}

func TestArchive_MovesRecordsToArchiveKey(t *testing.T) {
	st := archiveBaseState()
	res, err := Archive(st, "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.RecordsArchived != 2 {
		t.Errorf("expected 2 archived, got %d", res.RecordsArchived)
	}
	if len(st.ForEnvironment("staging")) != 0 {
		t.Error("expected live staging records to be removed")
	}
	archived := st.ForEnvironment("staging:archived")
	if len(archived) != 2 {
		t.Errorf("expected 2 records in archive, got %d", len(archived))
	}
}

func TestArchive_LeavesOtherEnvsUntouched(t *testing.T) {
	st := archiveBaseState()
	_, err := Archive(st, "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	prod := st.ForEnvironment("prod")
	if len(prod) != 1 {
		t.Errorf("expected prod to remain intact, got %d records", len(prod))
	}
}

func TestArchive_MissingEnvReturnsError(t *testing.T) {
	st := archiveBaseState()
	_, err := Archive(st, "nonexistent")
	if err == nil {
		t.Error("expected error for missing environment")
	}
}

func TestListArchived_ReturnsArchivedRecords(t *testing.T) {
	st := archiveBaseState()
	_, _ = Archive(st, "staging")
	records := ListArchived(st, "staging")
	if len(records) != 2 {
		t.Errorf("expected 2 archived records, got %d", len(records))
	}
}

func TestListArchived_EmptyWhenNoneArchived(t *testing.T) {
	st := archiveBaseState()
	records := ListArchived(st, "staging")
	if len(records) != 0 {
		t.Errorf("expected 0 archived records, got %d", len(records))
	}
}
