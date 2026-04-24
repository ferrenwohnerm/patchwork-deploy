package env_test

import (
	"testing"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
)

func watermarkBaseState() *state.State {
	st := state.New()
	st.Add("staging", "001-init.sql", "ok")
	st.Add("staging", "002-users.sql", "ok")
	return st
}

func TestSetWatermark_Success(t *testing.T) {
	st := watermarkBaseState()
	if err := env.SetWatermark(st, "staging", 5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := env.GetWatermark(st, "staging"); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestSetWatermark_ZeroLimitReturnsError(t *testing.T) {
	st := watermarkBaseState()
	if err := env.SetWatermark(st, "staging", 0); err == nil {
		t.Fatal("expected error for zero limit")
	}
}

func TestSetWatermark_MissingEnvReturnsError(t *testing.T) {
	st := watermarkBaseState()
	if err := env.SetWatermark(st, "ghost", 3); err == nil {
		t.Fatal("expected error for missing environment")
	}
}

func TestGetWatermark_NoneSet(t *testing.T) {
	st := watermarkBaseState()
	if got := env.GetWatermark(st, "staging"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestClearWatermark_RemovesEntry(t *testing.T) {
	st := watermarkBaseState()
	_ = env.SetWatermark(st, "staging", 10)
	if err := env.ClearWatermark(st, "staging"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := env.GetWatermark(st, "staging"); got != 0 {
		t.Fatalf("expected 0 after clear, got %d", got)
	}
}

func TestCheckWatermark_NotExceeded(t *testing.T) {
	st := watermarkBaseState()
	_ = env.SetWatermark(st, "staging", 10)
	exceeded, current, limit := env.CheckWatermark(st, "staging")
	if exceeded {
		t.Fatal("expected not exceeded")
	}
	if current != 2 {
		t.Fatalf("expected current=2, got %d", current)
	}
	if limit != 10 {
		t.Fatalf("expected limit=10, got %d", limit)
	}
}

func TestCheckWatermark_Exceeded(t *testing.T) {
	st := watermarkBaseState()
	_ = env.SetWatermark(st, "staging", 2)
	exceeded, _, _ := env.CheckWatermark(st, "staging")
	if !exceeded {
		t.Fatal("expected watermark to be exceeded")
	}
}

func TestCheckWatermark_NoWatermarkReturnsFalse(t *testing.T) {
	st := watermarkBaseState()
	exceeded, current, limit := env.CheckWatermark(st, "staging")
	if exceeded || current != 0 || limit != 0 {
		t.Fatal("expected all zeros when no watermark set")
	}
}
