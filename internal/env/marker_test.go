package env

import (
	"testing"

	"github.com/patchwork-deploy/internal/state"
)

func markerBaseState() *state.State {
	st := state.NewInMemory()
	st.AddEnvironment("staging")
	return st
}

func TestSetMarker_Success(t *testing.T) {
	st := markerBaseState()
	if err := SetMarker(st, "staging", "build", "42"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := GetMarker(st, "staging", "build")
	if !ok || v != "42" {
		t.Errorf("expected marker value '42', got %q (ok=%v)", v, ok)
	}
}

func TestSetMarker_MissingEnvReturnsError(t *testing.T) {
	st := markerBaseState()
	if err := SetMarker(st, "prod", "build", "1"); err == nil {
		t.Fatal("expected error for missing environment")
	}
}

func TestSetMarker_InvalidNameReturnsError(t *testing.T) {
	st := markerBaseState()
	if err := SetMarker(st, "staging", "bad name", "v"); err == nil {
		t.Fatal("expected error for whitespace in name")
	}
}

func TestSetMarker_NewlineInValueRejected(t *testing.T) {
	st := markerBaseState()
	if err := SetMarker(st, "staging", "info", "line1\nline2"); err == nil {
		t.Fatal("expected error for newline in value")
	}
}

func TestRemoveMarker_ClearsEntry(t *testing.T) {
	st := markerBaseState()
	_ = SetMarker(st, "staging", "build", "99")
	if err := RemoveMarker(st, "staging", "build"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := GetMarker(st, "staging", "build")
	if ok {
		t.Error("expected marker to be removed")
	}
}

func TestListMarkers_ReturnsAll(t *testing.T) {
	st := markerBaseState()
	_ = SetMarker(st, "staging", "build", "1")
	_ = SetMarker(st, "staging", "deploy", "2")
	m := ListMarkers(st, "staging")
	if len(m) != 2 {
		t.Errorf("expected 2 markers, got %d", len(m))
	}
	if m["build"] != "1" || m["deploy"] != "2" {
		t.Errorf("unexpected marker values: %v", m)
	}
}
