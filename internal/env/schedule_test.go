package env

import (
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func scheduleBaseState() *state.State {
	st := state.NewInMemory()
	now := time.Now()
	st.Add("prod", state.Record{Environment: "prod", Patch: "001-init.sql", AppliedAt: now})
	st.Add("prod", state.Record{Environment: "prod", Patch: "002-add-users.sql", AppliedAt: now})
	return st
}

func TestSchedule_AddsEntry(t *testing.T) {
	st := scheduleBaseState()
	runAt := time.Now().Add(24 * time.Hour)
	err := Schedule(st, "prod", "001-init.sql", runAt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries := ListScheduled(st, "prod")
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Patch != "001-init.sql" {
		t.Errorf("expected patch 001-init.sql, got %s", entries[0].Patch)
	}
}

func TestSchedule_MissingEnvReturnsError(t *testing.T) {
	st := scheduleBaseState()
	err := Schedule(st, "staging", "001-init.sql", time.Now())
	if err == nil {
		t.Fatal("expected error for missing environment")
	}
}

func TestSchedule_MissingPatchReturnsError(t *testing.T) {
	st := scheduleBaseState()
	err := Schedule(st, "prod", "999-missing.sql", time.Now())
	if err == nil {
		t.Fatal("expected error for missing patch")
	}
}

func TestCancelScheduled_RemovesEntry(t *testing.T) {
	st := scheduleBaseState()
	_ = Schedule(st, "prod", "001-init.sql", time.Now().Add(time.Hour))
	err := CancelScheduled(st, "prod", "001-init.sql")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ListScheduled(st, "prod")) != 0 {
		t.Error("expected no scheduled entries after cancel")
	}
}

func TestCancelScheduled_MissingEntryReturnsError(t *testing.T) {
	st := scheduleBaseState()
	err := CancelScheduled(st, "prod", "001-init.sql")
	if err == nil {
		t.Fatal("expected error when cancelling non-existent schedule")
	}
}

func TestListScheduled_EmptyWhenNoneScheduled(t *testing.T) {
	st := scheduleBaseState()
	entries := ListScheduled(st, "prod")
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}
