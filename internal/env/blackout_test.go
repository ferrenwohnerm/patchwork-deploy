package env

import (
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func blackoutBaseState() *state.State {
	st := state.NewInMemory()
	st.AddRecord("staging", state.Record{Patch: "001-init.sql", AppliedAt: time.Now()})
	return st
}

func TestSetBlackout_Success(t *testing.T) {
	st := blackoutBaseState()
	start := time.Now().Add(1 * time.Hour)
	end := start.Add(2 * time.Hour)
	if err := SetBlackout(st, "staging", start, end, "planned maintenance"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	bw, err := GetBlackout(st, "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bw == nil {
		t.Fatal("expected blackout window, got nil")
	}
	if bw.Note != "planned maintenance" {
		t.Errorf("expected note %q, got %q", "planned maintenance", bw.Note)
	}
}

func TestSetBlackout_MissingEnvReturnsError(t *testing.T) {
	st := blackoutBaseState()
	start := time.Now()
	err := SetBlackout(st, "ghost", start, start.Add(time.Hour), "")
	if err == nil {
		t.Fatal("expected error for missing env")
	}
}

func TestSetBlackout_InvalidRangeReturnsError(t *testing.T) {
	st := blackoutBaseState()
	now := time.Now()
	err := SetBlackout(st, "staging", now.Add(time.Hour), now, "")
	if err == nil {
		t.Fatal("expected error when end is before start")
	}
}

func TestIsBlackedOut_TrueWhenInsideWindow(t *testing.T) {
	st := blackoutBaseState()
	now := time.Now()
	_ = SetBlackout(st, "staging", now.Add(-1*time.Hour), now.Add(1*time.Hour), "")
	ok, err := IsBlackedOut(st, "staging", now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected environment to be blacked out")
	}
}

func TestIsBlackedOut_FalseWhenOutsideWindow(t *testing.T) {
	st := blackoutBaseState()
	now := time.Now()
	_ = SetBlackout(st, "staging", now.Add(2*time.Hour), now.Add(4*time.Hour), "")
	ok, _ := IsBlackedOut(st, "staging", now)
	if ok {
		t.Error("expected environment not to be blacked out")
	}
}

func TestClearBlackout_RemovesWindow(t *testing.T) {
	st := blackoutBaseState()
	now := time.Now()
	_ = SetBlackout(st, "staging", now, now.Add(time.Hour), "temp")
	_ = ClearBlackout(st, "staging")
	bw, err := GetBlackout(st, "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bw != nil {
		t.Error("expected nil after clear")
	}
}
