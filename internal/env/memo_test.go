package env

import (
	"testing"

	"github.com/patchwork-deploy/internal/state"
	"github.com/patchwork-deploy/internal/state/record"
)

func memoBaseState() *state.State {
	st := state.New()
	st.Add("staging", record.Record{Patch: "001-init.sql"})
	return st
}

func TestSetMemo_Success(t *testing.T) {
	st := memoBaseState()
	if err := SetMemo(st, "staging", "deployed by CI"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	val, err := GetMemo(st, "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "deployed by CI" {
		t.Errorf("expected memo %q, got %q", "deployed by CI", val)
	}
}

func TestSetMemo_MissingEnvReturnsError(t *testing.T) {
	st := memoBaseState()
	if err := SetMemo(st, "prod", "hello"); err == nil {
		t.Fatal("expected error for missing env")
	}
}

func TestSetMemo_NewlineRejected(t *testing.T) {
	st := memoBaseState()
	if err := SetMemo(st, "staging", "bad\nmemo"); err == nil {
		t.Fatal("expected error for newline in memo")
	}
}

func TestSetMemo_TooLongRejected(t *testing.T) {
	st := memoBaseState()
	long := string(make([]byte, 257))
	if err := SetMemo(st, "staging", long); err == nil {
		t.Fatal("expected error for memo exceeding limit")
	}
}

func TestGetMemo_NoneSet(t *testing.T) {
	st := memoBaseState()
	val, err := GetMemo(st, "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "" {
		t.Errorf("expected empty memo, got %q", val)
	}
}

func TestClearMemo_RemovesMemo(t *testing.T) {
	st := memoBaseState()
	_ = SetMemo(st, "staging", "to be removed")
	if err := ClearMemo(st, "staging"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	val, _ := GetMemo(st, "staging")
	if val != "" {
		t.Errorf("expected empty memo after clear, got %q", val)
	}
}
