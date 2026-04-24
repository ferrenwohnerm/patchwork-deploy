package env

import (
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func windowBaseState() *state.State {
	st := state.NewInMemory()
	st.Add("prod", state.Record{Patch: "001-init", AppliedAt: time.Now()})
	st.Add("prod", state.Record{Patch: "002-schema", AppliedAt: time.Now()})
	return st
}

func TestSetWindow_Success(t *testing.T) {
	st := windowBaseState()
	if err := SetWindow(st, "prod", "001-init", "08:00", "18:00"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s, e, ok := GetWindow(st, "prod", "001-init")
	if !ok {
		t.Fatal("expected window to be set")
	}
	if s != "08:00" || e != "18:00" {
		t.Errorf("got %s-%s, want 08:00-18:00", s, e)
	}
}

func TestSetWindow_MissingEnvReturnsError(t *testing.T) {
	st := windowBaseState()
	if err := SetWindow(st, "staging", "001-init", "08:00", "18:00"); err == nil {
		t.Fatal("expected error for missing env")
	}
}

func TestSetWindow_MissingPatchReturnsError(t *testing.T) {
	st := windowBaseState()
	if err := SetWindow(st, "prod", "999-missing", "08:00", "18:00"); err == nil {
		t.Fatal("expected error for missing patch")
	}
}

func TestSetWindow_InvalidTimeFormatReturnsError(t *testing.T) {
	st := windowBaseState()
	if err := SetWindow(st, "prod", "001-init", "8am", "18:00"); err == nil {
		t.Fatal("expected error for invalid start time")
	}
	if err := SetWindow(st, "prod", "001-init", "08:00", "6pm"); err == nil {
		t.Fatal("expected error for invalid end time")
	}
}

func TestSetWindow_StartNotBeforeEndReturnsError(t *testing.T) {
	st := windowBaseState()
	if err := SetWindow(st, "prod", "001-init", "18:00", "08:00"); err == nil {
		t.Fatal("expected error when start >= end")
	}
}

func TestClearWindow_RemovesEntry(t *testing.T) {
	st := windowBaseState()
	_ = SetWindow(st, "prod", "001-init", "09:00", "17:00")
	ClearWindow(st, "prod", "001-init")
	_, _, ok := GetWindow(st, "prod", "001-init")
	if ok {
		t.Fatal("expected window to be cleared")
	}
}

func TestListWindows_ReturnsAllEntries(t *testing.T) {
	st := windowBaseState()
	_ = SetWindow(st, "prod", "001-init", "08:00", "12:00")
	_ = SetWindow(st, "prod", "002-schema", "13:00", "17:00")
	win := ListWindows(st, "prod")
	if len(win) != 2 {
		t.Fatalf("expected 2 windows, got %d", len(win))
	}
	if win["001-init"] != [2]string{"08:00", "12:00"} {
		t.Errorf("unexpected window for 001-init: %v", win["001-init"])
	}
}
