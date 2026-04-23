package env

import (
	"testing"

	"github.com/your-org/patchwork-deploy/internal/state"
	"github.com/your-org/patchwork-deploy/internal/state/record"
)

func retentionBaseState() *state.State {
	st := state.New()
	st.Add(record.Record{Environment: "prod", Patch: "001-init.sql"})
	return st
}

func TestSetRetention_Success(t *testing.T) {
	st := retentionBaseState()
	if err := SetRetention(st, "prod", 10); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	n, ok := GetRetention(st, "prod")
	if !ok || n != 10 {
		t.Errorf("expected limit=10, got %d ok=%v", n, ok)
	}
}

func TestSetRetention_MissingEnvReturnsError(t *testing.T) {
	st := retentionBaseState()
	if err := SetRetention(st, "ghost", 5); err == nil {
		t.Fatal("expected error for missing env")
	}
}

func TestSetRetention_ZeroLimitReturnsError(t *testing.T) {
	st := retentionBaseState()
	if err := SetRetention(st, "prod", 0); err == nil {
		t.Fatal("expected error for zero limit")
	}
}

func TestGetRetention_NoneSet(t *testing.T) {
	st := retentionBaseState()
	_, ok := GetRetention(st, "prod")
	if ok {
		t.Error("expected no retention limit")
	}
}

func TestClearRetention_RemovesLimit(t *testing.T) {
	st := retentionBaseState()
	_ = SetRetention(st, "prod", 5)
	if err := ClearRetention(st, "prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := GetRetention(st, "prod")
	if ok {
		t.Error("expected retention to be cleared")
	}
}

func TestClearRetention_MissingEnvReturnsError(t *testing.T) {
	st := retentionBaseState()
	if err := ClearRetention(st, "ghost"); err == nil {
		t.Fatal("expected error for missing env")
	}
}
