package env

import (
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func baselineBaseState() *state.State {
	st := state.New()
	now := time.Now().UTC()
	st.Add(state.Record{Environment: "prod", Patch: "001_init", AppliedAt: now.Add(-2 * time.Hour)})
	st.Add(state.Record{Environment: "prod", Patch: "002_schema", AppliedAt: now.Add(-1 * time.Hour)})
	st.Add(state.Record{Environment: "prod", Patch: "003_data", AppliedAt: now})
	return st
}

func TestSetBaseline_Success(t *testing.T) {
	st := baselineBaseState()
	if err := SetBaseline(st, "prod", "002_schema"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, ok := GetBaseline(st, "prod")
	if !ok {
		t.Fatal("expected baseline to be set")
	}
	if b.Patch != "002_schema" {
		t.Errorf("expected 002_schema, got %s", b.Patch)
	}
	if b.Environment != "prod" {
		t.Errorf("expected prod, got %s", b.Environment)
	}
}

func TestSetBaseline_MissingEnvReturnsError(t *testing.T) {
	st := baselineBaseState()
	if err := SetBaseline(st, "staging", "001_init"); err == nil {
		t.Fatal("expected error for missing environment")
	}
}

func TestSetBaseline_MissingPatchReturnsError(t *testing.T) {
	st := baselineBaseState()
	if err := SetBaseline(st, "prod", "999_missing"); err == nil {
		t.Fatal("expected error for missing patch")
	}
}

func TestGetBaseline_NoneSet(t *testing.T) {
	st := baselineBaseState()
	_, ok := GetBaseline(st, "prod")
	if ok {
		t.Fatal("expected no baseline")
	}
}

func TestClearBaseline_RemovesEntry(t *testing.T) {
	st := baselineBaseState()
	_ = SetBaseline(st, "prod", "001_init")
	if err := ClearBaseline(st, "prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := GetBaseline(st, "prod")
	if ok {
		t.Fatal("expected baseline to be cleared")
	}
}

func TestClearBaseline_NoneSetReturnsError(t *testing.T) {
	st := baselineBaseState()
	if err := ClearBaseline(st, "prod"); err == nil {
		t.Fatal("expected error when no baseline is set")
	}
}
