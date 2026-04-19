package env

import (
	"testing"

	"github.com/patchwork-deploy/internal/state"
)

func groupBaseState() *state.State {
	st := state.New()
	st.AddEnvironment("staging")
	st.AddEnvironment("production")
	st.AddEnvironment("dev")
	return st
}

func TestAddToGroup_Success(t *testing.T) {
	st := groupBaseState()
	if err := AddToGroup(st, "web", "staging"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	members, err := ListGroup(st, "web")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(members) != 1 || members[0] != "staging" {
		t.Errorf("expected [staging], got %v", members)
	}
}

func TestAddToGroup_Idempotent(t *testing.T) {
	st := groupBaseState()
	_ = AddToGroup(st, "web", "staging")
	_ = AddToGroup(st, "web", "staging")
	members, _ := ListGroup(st, "web")
	if len(members) != 1 {
		t.Errorf("expected 1 member, got %d", len(members))
	}
}

func TestAddToGroup_MissingEnvReturnsError(t *testing.T) {
	st := groupBaseState()
	if err := AddToGroup(st, "web", "ghost"); err == nil {
		t.Error("expected error for missing environment")
	}
}

func TestAddToGroup_InvalidGroupNameReturnsError(t *testing.T) {
	st := groupBaseState()
	if err := AddToGroup(st, "bad name!", "staging"); err == nil {
		t.Error("expected error for invalid group name")
	}
}

func TestRemoveFromGroup_RemovesMember(t *testing.T) {
	st := groupBaseState()
	_ = AddToGroup(st, "web", "staging")
	_ = AddToGroup(st, "web", "production")
	_ = RemoveFromGroup(st, "web", "staging")
	members, _ := ListGroup(st, "web")
	if len(members) != 1 || members[0] != "production" {
		t.Errorf("expected [production], got %v", members)
	}
}

func TestListAllGroups_ReturnsAllGroups(t *testing.T) {
	st := groupBaseState()
	_ = AddToGroup(st, "web", "staging")
	_ = AddToGroup(st, "infra", "production")
	groups := ListAllGroups(st)
	if len(groups) != 2 {
		t.Errorf("expected 2 groups, got %d", len(groups))
	}
}

func TestListGroup_UnknownGroupReturnsError(t *testing.T) {
	st := groupBaseState()
	if _, err := ListGroup(st, "nope"); err == nil {
		t.Error("expected error for unknown group")
	}
}
