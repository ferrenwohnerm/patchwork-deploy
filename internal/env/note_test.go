package env

import (
	"strings"
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func noteBaseState() *state.State {
	st := state.NewInMemory()
	st.Add(state.Record{Environment: "staging", Patch: "001-init", AppliedAt: time.Now()})
	return st
}

func TestSetNote_Success(t *testing.T) {
	st := noteBaseState()
	if err := SetNote(st, "staging", "deployed by alice"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := GetNote(st, "staging")
	if !ok || v != "deployed by alice" {
		t.Fatalf("expected note 'deployed by alice', got %q (ok=%v)", v, ok)
	}
}

func TestSetNote_MissingEnvReturnsError(t *testing.T) {
	st := noteBaseState()
	if err := SetNote(st, "ghost", "hello"); err == nil {
		t.Fatal("expected error for missing env")
	}
}

func TestSetNote_NewlineRejected(t *testing.T) {
	st := noteBaseState()
	if err := SetNote(st, "staging", "line1\nline2"); err == nil {
		t.Fatal("expected error for newline in note")
	}
}

func TestSetNote_TooLongRejected(t *testing.T) {
	st := noteBaseState()
	long := strings.Repeat("x", 501)
	if err := SetNote(st, "staging", long); err == nil {
		t.Fatal("expected error for oversized note")
	}
}

func TestClearNote_RemovesEntry(t *testing.T) {
	st := noteBaseState()
	_ = SetNote(st, "staging", "temporary note")
	if err := ClearNote(st, "staging"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := GetNote(st, "staging")
	if ok {
		t.Fatal("expected note to be cleared")
	}
}

func TestGetNote_MissingEnvReturnsFalse(t *testing.T) {
	st := noteBaseState()
	_, ok := GetNote(st, "nonexistent")
	if ok {
		t.Fatal("expected false for missing env")
	}
}
