package env

import (
	"testing"

	"github.com/patchwork-deploy/internal/state"
)

func badgeBaseState() *state.State {
	st := state.New()
	st.Add("prod", state.Record{Patch: "001-init.sql"})
	return st
}

func TestSetBadge_Success(t *testing.T) {
	st := badgeBaseState()
	if err := SetBadge(st, "prod", "status", "stable"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := GetBadge(st, "prod", "status")
	if !ok || v != "stable" {
		t.Fatalf("expected badge value 'stable', got %q (ok=%v)", v, ok)
	}
}

func TestSetBadge_InvalidKey(t *testing.T) {
	st := badgeBaseState()
	if err := SetBadge(st, "prod", "bad key!", "v"); err == nil {
		t.Fatal("expected error for invalid key")
	}
}

func TestSetBadge_NewlineInValueRejected(t *testing.T) {
	st := badgeBaseState()
	if err := SetBadge(st, "prod", "info", "line1\nline2"); err == nil {
		t.Fatal("expected error for newline in value")
	}
}

func TestSetBadge_MissingEnvReturnsError(t *testing.T) {
	st := badgeBaseState()
	if err := SetBadge(st, "staging", "status", "ok"); err == nil {
		t.Fatal("expected error for missing env")
	}
}

func TestRemoveBadge_ClearsBadge(t *testing.T) {
	st := badgeBaseState()
	_ = SetBadge(st, "prod", "tier", "gold")
	if err := RemoveBadge(st, "prod", "tier"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := GetBadge(st, "prod", "tier")
	if ok {
		t.Fatal("expected badge to be removed")
	}
}

func TestListBadges_ReturnsAll(t *testing.T) {
	st := badgeBaseState()
	_ = SetBadge(st, "prod", "status", "stable")
	_ = SetBadge(st, "prod", "owner", "team-a")
	badges := ListBadges(st, "prod")
	if len(badges) != 2 {
		t.Fatalf("expected 2 badges, got %d", len(badges))
	}
	if badges["status"] != "stable" {
		t.Errorf("expected status=stable")
	}
	if badges["owner"] != "team-a" {
		t.Errorf("expected owner=team-a")
	}
}
