package env

import (
	"testing"

	"github.com/patchwork-deploy/internal/state"
	"github.com/patchwork-deploy/internal/state/record"
)

func aliasBaseState() *state.State {
	st := state.NewInMemory()
	st.Add(record.Record{Environment: "production", Patch: "001-init.sql"})
	st.Add(record.Record{Environment: "staging", Patch: "001-init.sql"})
	return st
}

func TestSetAlias_Success(t *testing.T) {
	st := aliasBaseState()
	if err := SetAlias(st, "production", "prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := GetAlias(st, "production"); got != "prod" {
		t.Errorf("expected prod, got %q", got)
	}
}

func TestSetAlias_InvalidChars(t *testing.T) {
	st := aliasBaseState()
	if err := SetAlias(st, "production", "PROD!"); err == nil {
		t.Error("expected error for invalid alias chars")
	}
}

func TestSetAlias_MissingEnvReturnsError(t *testing.T) {
	st := aliasBaseState()
	if err := SetAlias(st, "ghost", "g"); err == nil {
		t.Error("expected error for missing env")
	}
}

func TestRemoveAlias_ClearsAlias(t *testing.T) {
	st := aliasBaseState()
	_ = SetAlias(st, "staging", "stg")
	if err := RemoveAlias(st, "staging"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := GetAlias(st, "staging"); got != "" {
		t.Errorf("expected empty alias after removal, got %q", got)
	}
}

func TestResolveAlias_ReturnsEnvName(t *testing.T) {
	st := aliasBaseState()
	_ = SetAlias(st, "production", "prod")
	if got := ResolveAlias(st, "prod"); got != "production" {
		t.Errorf("expected production, got %q", got)
	}
}

func TestResolveAlias_UnknownReturnsEmpty(t *testing.T) {
	st := aliasBaseState()
	if got := ResolveAlias(st, "nope"); got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}
