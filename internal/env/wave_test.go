package env

import (
	"testing"

	"github.com/patchwork-deploy/internal/state"
	"github.com/patchwork-deploy/internal/state/record"
)

func waveBaseState() *state.State {
	st := state.New()
	st.Add(record.Record{Environment: "staging", Patch: "001-init"})
	st.Add(record.Record{Environment: "staging", Patch: "002-schema"})
	st.Add(record.Record{Environment: "staging", Patch: "003-data"})
	return st
}

func TestSetWave_Success(t *testing.T) {
	st := waveBaseState()
	if err := SetWave(st, "staging", "001-init", 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := GetWave(st, "staging", "001-init"); got != 1 {
		t.Fatalf("expected wave 1, got %d", got)
	}
}

func TestSetWave_NegativeReturnsError(t *testing.T) {
	st := waveBaseState()
	if err := SetWave(st, "staging", "001-init", -1); err == nil {
		t.Fatal("expected error for negative wave")
	}
}

func TestSetWave_MissingEnvReturnsError(t *testing.T) {
	st := waveBaseState()
	if err := SetWave(st, "prod", "001-init", 1); err == nil {
		t.Fatal("expected error for missing environment")
	}
}

func TestSetWave_MissingPatchReturnsError(t *testing.T) {
	st := waveBaseState()
	if err := SetWave(st, "staging", "999-nope", 1); err == nil {
		t.Fatal("expected error for missing patch")
	}
}

func TestGetWave_NoneSet(t *testing.T) {
	st := waveBaseState()
	if got := GetWave(st, "staging", "001-init"); got != -1 {
		t.Fatalf("expected -1 when no wave set, got %d", got)
	}
}

func TestClearWave_RemovesEntry(t *testing.T) {
	st := waveBaseState()
	_ = SetWave(st, "staging", "001-init", 2)
	ClearWave(st, "staging", "001-init")
	if got := GetWave(st, "staging", "001-init"); got != -1 {
		t.Fatalf("expected -1 after clear, got %d", got)
	}
}

func TestListWaves_SortedByWaveThenPatch(t *testing.T) {
	st := waveBaseState()
	_ = SetWave(st, "staging", "003-data", 1)
	_ = SetWave(st, "staging", "001-init", 1)
	_ = SetWave(st, "staging", "002-schema", 2)

	entries := ListWaves(st, "staging")
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Patch != "001-init" || entries[0].Wave != 1 {
		t.Errorf("unexpected first entry: %+v", entries[0])
	}
	if entries[1].Patch != "003-data" || entries[1].Wave != 1 {
		t.Errorf("unexpected second entry: %+v", entries[1])
	}
	if entries[2].Wave != 2 {
		t.Errorf("unexpected third entry wave: %+v", entries[2])
	}
}
