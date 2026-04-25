package env

import (
	"testing"

	"github.com/patchwork-deploy/internal/state"
	"github.com/patchwork-deploy/internal/state/record"
)

func strandBaseState() *state.State {
	st := state.New()
	st.Add(record.Record{Environment: "staging", Patch: "001-init"})
	st.Add(record.Record{Environment: "staging", Patch: "002-schema"})
	return st
}

func TestSetStrand_Success(t *testing.T) {
	st := strandBaseState()
	if err := SetStrand(st, "staging", "001-init", "fast-track"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := GetStrand(st, "staging", "001-init"); got != "fast-track" {
		t.Errorf("expected fast-track, got %q", got)
	}
}

func TestSetStrand_InvalidNameReturnsError(t *testing.T) {
	st := strandBaseState()
	if err := SetStrand(st, "staging", "001-init", "bad name!"); err == nil {
		t.Fatal("expected error for invalid strand name")
	}
}

func TestSetStrand_MissingEnvReturnsError(t *testing.T) {
	st := strandBaseState()
	if err := SetStrand(st, "prod", "001-init", "main"); err == nil {
		t.Fatal("expected error for missing environment")
	}
}

func TestSetStrand_MissingPatchReturnsError(t *testing.T) {
	st := strandBaseState()
	if err := SetStrand(st, "staging", "999-nope", "main"); err == nil {
		t.Fatal("expected error for missing patch")
	}
}

func TestClearStrand_RemovesEntry(t *testing.T) {
	st := strandBaseState()
	_ = SetStrand(st, "staging", "001-init", "hotfix")
	ClearStrand(st, "staging", "001-init")
	if got := GetStrand(st, "staging", "001-init"); got != "" {
		t.Errorf("expected empty after clear, got %q", got)
	}
}

func TestListStrands_ReturnsAllForEnv(t *testing.T) {
	st := strandBaseState()
	_ = SetStrand(st, "staging", "001-init", "main")
	_ = SetStrand(st, "staging", "002-schema", "hotfix")
	m := ListStrands(st, "staging")
	if len(m) != 2 {
		t.Fatalf("expected 2 strands, got %d", len(m))
	}
	if m["001-init"] != "main" || m["002-schema"] != "hotfix" {
		t.Errorf("unexpected strand map: %v", m)
	}
}
