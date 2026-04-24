package env

import (
	"testing"

	"github.com/patchwork-deploy/internal/state"
)

func fenceBaseState() *state.State {
	st := state.NewInMemory()
	st.Add("prod", "001-init.sql", "ok")
	st.Add("prod", "002-schema.sql", "ok")
	return st
}

func TestSetFence_Success(t *testing.T) {
	st := fenceBaseState()
	if err := SetFence(st, "prod", "001-init.sql", "gate-A"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := GetFence(st, "prod", "001-init.sql"); got != "gate-A" {
		t.Errorf("expected gate-A, got %q", got)
	}
}

func TestSetFence_InvalidNameReturnsError(t *testing.T) {
	st := fenceBaseState()
	if err := SetFence(st, "prod", "001-init.sql", "bad name!"); err == nil {
		t.Fatal("expected error for invalid fence name")
	}
}

func TestSetFence_MissingEnvReturnsError(t *testing.T) {
	st := fenceBaseState()
	if err := SetFence(st, "staging", "001-init.sql", "gate-A"); err == nil {
		t.Fatal("expected error for missing environment")
	}
}

func TestSetFence_MissingPatchReturnsError(t *testing.T) {
	st := fenceBaseState()
	if err := SetFence(st, "prod", "999-missing.sql", "gate-A"); err == nil {
		t.Fatal("expected error for missing patch")
	}
}

func TestClearFence_RemovesSentinel(t *testing.T) {
	st := fenceBaseState()
	_ = SetFence(st, "prod", "001-init.sql", "gate-A")
	if err := ClearFence(st, "prod", "001-init.sql"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if IsFenced(st, "prod", "001-init.sql") {
		t.Error("expected fence to be cleared")
	}
}

func TestListFences_ReturnsAllForEnv(t *testing.T) {
	st := fenceBaseState()
	_ = SetFence(st, "prod", "001-init.sql", "gate-A")
	_ = SetFence(st, "prod", "002-schema.sql", "gate-B")
	fences := ListFences(st, "prod")
	if len(fences) != 2 {
		t.Fatalf("expected 2 fences, got %d", len(fences))
	}
	if fences["001-init.sql"] != "gate-A" {
		t.Errorf("expected gate-A for 001-init.sql, got %q", fences["001-init.sql"])
	}
	if fences["002-schema.sql"] != "gate-B" {
		t.Errorf("expected gate-B for 002-schema.sql, got %q", fences["002-schema.sql"])
	}
}
