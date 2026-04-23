package env

import (
	"testing"

	"github.com/patchwork-deploy/internal/state"
	"github.com/patchwork-deploy/internal/state/record"
)

func shadowBaseState() *state.State {
	st := state.New()
	st.Add("staging", record.Record{Patch: "001-init.sql"})
	st.Add("staging", record.Record{Patch: "002-users.sql"})
	st.Add("mirror", record.Record{Patch: "001-init.sql"})
	return st
}

func TestSetShadow_Success(t *testing.T) {
	st := shadowBaseState()
	if err := SetShadow(st, "staging", "mirror", "001-init.sql"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	src := GetShadow(st, "mirror", "001-init.sql")
	if src != "staging" {
		t.Errorf("expected source=staging, got %q", src)
	}
}

func TestSetShadow_SameEnvReturnsError(t *testing.T) {
	st := shadowBaseState()
	if err := SetShadow(st, "staging", "staging", "001-init.sql"); err == nil {
		t.Fatal("expected error for same env")
	}
}

func TestSetShadow_MissingSourceReturnsError(t *testing.T) {
	st := shadowBaseState()
	if err := SetShadow(st, "nonexistent", "mirror", "001-init.sql"); err == nil {
		t.Fatal("expected error for missing source env")
	}
}

func TestSetShadow_MissingPatchReturnsError(t *testing.T) {
	st := shadowBaseState()
	if err := SetShadow(st, "staging", "mirror", "999-missing.sql"); err == nil {
		t.Fatal("expected error for missing patch")
	}
}

func TestRemoveShadow_ClearsSentinel(t *testing.T) {
	st := shadowBaseState()
	_ = SetShadow(st, "staging", "mirror", "001-init.sql")
	if err := RemoveShadow(st, "mirror", "001-init.sql"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src := GetShadow(st, "mirror", "001-init.sql"); src != "" {
		t.Errorf("expected empty after removal, got %q", src)
	}
}

func TestListShadows_ReturnsAllEntries(t *testing.T) {
	st := shadowBaseState()
	_ = SetShadow(st, "staging", "mirror", "001-init.sql")
	_ = SetShadow(st, "staging", "mirror", "002-users.sql")
	shadows := ListShadows(st, "mirror")
	if len(shadows) != 2 {
		t.Errorf("expected 2 shadows, got %d", len(shadows))
	}
	for patch, src := range shadows {
		if src != "staging" {
			t.Errorf("patch %q: expected source=staging, got %q", patch, src)
		}
	}
}

func TestGetShadow_NoneSet(t *testing.T) {
	st := shadowBaseState()
	if src := GetShadow(st, "mirror", "001-init.sql"); src != "" {
		t.Errorf("expected empty string, got %q", src)
	}
}
