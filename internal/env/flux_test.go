package env

import (
	"testing"

	"github.com/patchwork-deploy/internal/state"
)

func fluxBaseState() *state.State {
	st := state.NewInMemory()
	st.Add("dev", "001-init.sql", "", nil)
	st.Add("dev", "002-schema.sql", "", nil)
	return st
}

func TestSetFlux_Success(t *testing.T) {
	st := fluxBaseState()
	if err := SetFlux(st, "dev", "001-init.sql", "auto"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := GetFlux(st, "dev", "001-init.sql"); got != "auto" {
		t.Errorf("expected auto, got %q", got)
	}
}

func TestSetFlux_InvalidModeReturnsError(t *testing.T) {
	st := fluxBaseState()
	if err := SetFlux(st, "dev", "001-init.sql", "unknown"); err == nil {
		t.Fatal("expected error for invalid mode")
	}
}

func TestSetFlux_MissingEnvReturnsError(t *testing.T) {
	st := fluxBaseState()
	if err := SetFlux(st, "staging", "001-init.sql", "manual"); err == nil {
		t.Fatal("expected error for missing environment")
	}
}

func TestSetFlux_MissingPatchReturnsError(t *testing.T) {
	st := fluxBaseState()
	if err := SetFlux(st, "dev", "999-missing.sql", "gated"); err == nil {
		t.Fatal("expected error for missing patch")
	}
}

func TestGetFlux_NoneSet(t *testing.T) {
	st := fluxBaseState()
	if got := GetFlux(st, "dev", "001-init.sql"); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestClearFlux_RemovesEntry(t *testing.T) {
	st := fluxBaseState()
	_ = SetFlux(st, "dev", "001-init.sql", "gated")
	if err := ClearFlux(st, "dev", "001-init.sql"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := GetFlux(st, "dev", "001-init.sql"); got != "" {
		t.Errorf("expected empty after clear, got %q", got)
	}
}

func TestListFluxes_ReturnsAllEntries(t *testing.T) {
	st := fluxBaseState()
	_ = SetFlux(st, "dev", "001-init.sql", "auto")
	_ = SetFlux(st, "dev", "002-schema.sql", "manual")
	fluxes := ListFluxes(st, "dev")
	if len(fluxes) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(fluxes))
	}
	if fluxes["001-init.sql"] != "auto" {
		t.Errorf("expected auto for 001-init.sql, got %q", fluxes["001-init.sql"])
	}
	if fluxes["002-schema.sql"] != "manual" {
		t.Errorf("expected manual for 002-schema.sql, got %q", fluxes["002-schema.sql"])
	}
}
