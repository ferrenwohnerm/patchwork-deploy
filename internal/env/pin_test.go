package env

import (
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func pinBaseState() *state.State {
	st := state.New()
	st.Add(state.Record{Environment: "prod", Patch: "001-init", AppliedAt: time.Now()})
	st.Add(state.Record{Environment: "prod", Patch: "002-add-table", AppliedAt: time.Now()})
	st.Add(state.Record{Environment: "prod", Patch: "003-index", AppliedAt: time.Now()})
	return st
}

func TestPin_MarksPatch(t *testing.T) {
	st := pinBaseState()
	res, err := Pin(st, "prod", "002-add-table")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Patch != "002-add-table" {
		t.Errorf("expected pinned patch 002-add-table, got %s", res.Patch)
	}
	if PinnedPatch(st, "prod") != "002-add-table" {
		t.Errorf("PinnedPatch should return 002-add-table")
	}
}

func TestPin_MissingEnvReturnsError(t *testing.T) {
	st := pinBaseState()
	_, err := Pin(st, "staging", "001-init")
	if err == nil {
		t.Fatal("expected error for missing env")
	}
}

func TestPin_MissingPatchReturnsError(t *testing.T) {
	st := pinBaseState()
	_, err := Pin(st, "prod", "999-nope")
	if err == nil {
		t.Fatal("expected error for missing patch")
	}
}

func TestUnpin_ClearsSentinel(t *testing.T) {
	st := pinBaseState()
	_, _ = Pin(st, "prod", "001-init")
	if err := Unpin(st, "prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if PinnedPatch(st, "prod") != "" {
		t.Errorf("expected no pin after unpin")
	}
}

func TestUnpin_MissingEnvReturnsError(t *testing.T) {
	st := pinBaseState()
	if err := Unpin(st, "ghost"); err == nil {
		t.Fatal("expected error for missing env")
	}
}

func TestPinnedPatch_NoneReturnsEmpty(t *testing.T) {
	st := pinBaseState()
	if p := PinnedPatch(st, "prod"); p != "" {
		t.Errorf("expected empty pin, got %s", p)
	}
}
