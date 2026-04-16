package env

import (
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func copyBaseState() *state.State {
	st := state.New()
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	st.Add(state.Record{Environment: "staging", PatchID: "001-init", AppliedAt: base})
	st.Add(state.Record{Environment: "staging", PatchID: "002-users", AppliedAt: base.Add(time.Hour)})
	st.Add(state.Record{Environment: "prod", PatchID: "001-init", AppliedAt: base.Add(2 * time.Hour)})
	return st
}

func TestCopy_CopiesAllPatches(t *testing.T) {
	st := copyBaseState()
	res, err := Copy(st, "staging", "prod", CopyOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Copied) != 1 || res.Copied[0] != "002-users" {
		t.Errorf("expected [002-users] copied, got %v", res.Copied)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "001-init" {
		t.Errorf("expected [001-init] skipped, got %v", res.Skipped)
	}
	prodRecords := st.ForEnvironment("prod")
	if len(prodRecords) != 2 {
		t.Errorf("expected 2 prod records, got %d", len(prodRecords))
	}
}

func TestCopy_SelectivePatches(t *testing.T) {
	st := copyBaseState()
	res, err := Copy(st, "staging", "prod", CopyOptions{PatchIDs: []string{"001-init"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Copied) != 0 {
		t.Errorf("expected nothing copied, got %v", res.Copied)
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %v", res.Skipped)
	}
}

func TestCopy_SameEnvReturnsError(t *testing.T) {
	st := copyBaseState()
	_, err := Copy(st, "staging", "staging", CopyOptions{})
	if err == nil {
		t.Fatal("expected error for same src/dst, got nil")
	}
}

func TestCopy_MissingSourceReturnsError(t *testing.T) {
	st := copyBaseState()
	_, err := Copy(st, "nonexistent", "prod", CopyOptions{})
	if err == nil {
		t.Fatal("expected error for missing source, got nil")
	}
}

func TestCopy_SourceRecordsUnchanged(t *testing.T) {
	st := copyBaseState()
	before := len(st.ForEnvironment("staging"))
	_, err := Copy(st, "staging", "prod", CopyOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	after := len(st.ForEnvironment("staging"))
	if before != after {
		t.Errorf("source records changed: before=%d after=%d", before, after)
	}
}
