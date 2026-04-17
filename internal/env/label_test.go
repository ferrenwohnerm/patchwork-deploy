package env

import (
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func labelBaseState() *state.State {
	st := state.New()
	now := time.Now()
	st.Add(state.Record{Environment: "staging", Patch: "001-init.sql", AppliedAt: now})
	st.Add(state.Record{Environment: "staging", Patch: "002-users.sql", AppliedAt: now})
	return st
}

func TestSetLabel_AddsLabel(t *testing.T) {
	st := labelBaseState()
	if err := SetLabel(st, "staging", "team", "backend"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	labels, err := ListLabels(st, "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if labels["team"] != "backend" {
		t.Errorf("expected label team=backend, got %v", labels)
	}
}

func TestSetLabel_InvalidKey(t *testing.T) {
	st := labelBaseState()
	if err := SetLabel(st, "staging", "bad key!", "val"); err == nil {
		t.Error("expected error for invalid key")
	}
}

func TestSetLabel_MissingEnv(t *testing.T) {
	st := labelBaseState()
	if err := SetLabel(st, "prod", "team", "ops"); err == nil {
		t.Error("expected error for missing environment")
	}
}

func TestRemoveLabel_ClearsLabel(t *testing.T) {
	st := labelBaseState()
	_ = SetLabel(st, "staging", "owner", "alice")
	if err := RemoveLabel(st, "staging", "owner"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	labels, _ := ListLabels(st, "staging")
	if _, ok := labels["owner"]; ok {
		t.Error("expected label to be removed")
	}
}

func TestRemoveLabel_NotFound(t *testing.T) {
	st := labelBaseState()
	if err := RemoveLabel(st, "staging", "nonexistent"); err == nil {
		t.Error("expected error when label not found")
	}
}

func TestListLabels_ReturnsAllLabels(t *testing.T) {
	st := labelBaseState()
	_ = SetLabel(st, "staging", "team", "backend")
	_ = SetLabel(st, "staging", "region", "us-east")
	labels, err := ListLabels(st, "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(labels) != 2 {
		t.Errorf("expected 2 labels, got %d", len(labels))
	}
}
