package env

import (
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func expireBaseState() *state.State {
	st := state.New()
	now := time.Now()
	st.Add(state.Record{Environment: "staging", Patch: "001-init.sql", AppliedAt: now})
	st.Add(state.Record{Environment: "prod", Patch: "001-init.sql", AppliedAt: now})
	return st
}

func TestSetExpiry_Success(t *testing.T) {
	st := expireBaseState()
	at := time.Now().Add(24 * time.Hour)
	if err := SetExpiry(st, "staging", at); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok, err := GetExpiry(st, "staging")
	if err != nil || !ok {
		t.Fatalf("expected expiry to be set, ok=%v err=%v", ok, err)
	}
	if got.Unix() != at.UTC().Truncate(time.Second).Unix() {
		t.Errorf("expiry mismatch: got %v", got)
	}
}

func TestSetExpiry_MissingEnvReturnsError(t *testing.T) {
	st := expireBaseState()
	err := SetExpiry(st, "ghost", time.Now().Add(time.Hour))
	if err == nil {
		t.Fatal("expected error for missing env")
	}
}

func TestIsExpired_FalseWhenFuture(t *testing.T) {
	st := expireBaseState()
	_ = SetExpiry(st, "staging", time.Now().Add(time.Hour))
	expired, err := IsExpired(st, "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if expired {
		t.Error("expected not expired")
	}
}

func TestIsExpired_TrueWhenPast(t *testing.T) {
	st := expireBaseState()
	_ = SetExpiry(st, "staging", time.Now().Add(-time.Hour))
	expired, err := IsExpired(st, "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !expired {
		t.Error("expected expired")
	}
}

func TestClearExpiry_RemovesSentinel(t *testing.T) {
	st := expireBaseState()
	_ = SetExpiry(st, "staging", time.Now().Add(time.Hour))
	_ = ClearExpiry(st, "staging")
	_, ok, _ := GetExpiry(st, "staging")
	if ok {
		t.Error("expected expiry to be cleared")
	}
}

func TestListExpired_ReturnsExpiredEnvs(t *testing.T) {
	st := expireBaseState()
	_ = SetExpiry(st, "staging", time.Now().Add(-time.Hour))
	_ = SetExpiry(st, "prod", time.Now().Add(time.Hour))
	list, err := ListExpired(st)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 1 || list[0] != "staging" {
		t.Errorf("expected [staging], got %v", list)
	}
}
