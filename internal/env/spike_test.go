package env

import (
	"testing"

	"github.com/patchwork-deploy/internal/state"
)

func spikeBaseState() *state.State {
	st := state.NewInMemory()
	st.AddEnvironment("staging")
	st.AddRecord("staging", state.Record{Patch: "001-init.sql"})
	st.AddRecord("staging", state.Record{Patch: "002-add-index.sql"})
	return st
}

func TestSetSpike_Success(t *testing.T) {
	st := spikeBaseState()
	if err := SetSpike(st, "staging", "001-init.sql", 3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	limit, ok := GetSpike(st, "staging", "001-init.sql")
	if !ok {
		t.Fatal("expected spike to be set")
	}
	if limit != 3 {
		t.Fatalf("expected limit 3, got %d", limit)
	}
}

func TestSetSpike_ZeroLimitReturnsError(t *testing.T) {
	st := spikeBaseState()
	if err := SetSpike(st, "staging", "001-init.sql", 0); err == nil {
		t.Fatal("expected error for zero limit")
	}
}

func TestSetSpike_MissingEnvReturnsError(t *testing.T) {
	st := spikeBaseState()
	if err := SetSpike(st, "prod", "001-init.sql", 2); err == nil {
		t.Fatal("expected error for missing env")
	}
}

func TestSetSpike_MissingPatchReturnsError(t *testing.T) {
	st := spikeBaseState()
	if err := SetSpike(st, "staging", "999-missing.sql", 2); err == nil {
		t.Fatal("expected error for missing patch")
	}
}

func TestGetSpike_NoneSet(t *testing.T) {
	st := spikeBaseState()
	_, ok := GetSpike(st, "staging", "001-init.sql")
	if ok {
		t.Fatal("expected no spike to be set")
	}
}

func TestClearSpike_RemovesEntry(t *testing.T) {
	st := spikeBaseState()
	_ = SetSpike(st, "staging", "001-init.sql", 5)
	ClearSpike(st, "staging", "001-init.sql")
	_, ok := GetSpike(st, "staging", "001-init.sql")
	if ok {
		t.Fatal("expected spike to be cleared")
	}
}

func TestListSpikes_ReturnsAll(t *testing.T) {
	st := spikeBaseState()
	_ = SetSpike(st, "staging", "001-init.sql", 2)
	_ = SetSpike(st, "staging", "002-add-index.sql", 4)
	spikes := ListSpikes(st, "staging")
	if len(spikes) != 2 {
		t.Fatalf("expected 2 spikes, got %d", len(spikes))
	}
	if spikes["001-init.sql"] != 2 {
		t.Errorf("expected limit 2 for 001-init.sql, got %d", spikes["001-init.sql"])
	}
	if spikes["002-add-index.sql"] != 4 {
		t.Errorf("expected limit 4 for 002-add-index.sql, got %d", spikes["002-add-index.sql"])
	}
}
