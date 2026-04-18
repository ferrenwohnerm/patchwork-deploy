package env

import (
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func quotaBaseState() *state.State {
	st := state.New()
	now := time.Now()
	st.Add(state.Record{Environment: "staging", Patch: "001-init", AppliedAt: now})
	st.Add(state.Record{Environment: "staging", Patch: "002-users", AppliedAt: now})
	st.Add(state.Record{Environment: "prod", Patch: "001-init", AppliedAt: now})
	return st
}

func TestSetQuota_Success(t *testing.T) {
	st := quotaBaseState()
	if err := SetQuota(st, "staging", 5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	limit, ok := GetQuota(st, "staging")
	if !ok || limit != 5 {
		t.Errorf("expected quota 5, got %d (ok=%v)", limit, ok)
	}
}

func TestSetQuota_MissingEnvReturnsError(t *testing.T) {
	st := quotaBaseState()
	if err := SetQuota(st, "ghost", 3); err == nil {
		t.Fatal("expected error for missing env")
	}
}

func TestSetQuota_ZeroLimitReturnsError(t *testing.T) {
	st := quotaBaseState()
	if err := SetQuota(st, "staging", 0); err == nil {
		t.Fatal("expected error for zero quota")
	}
}

func TestCheckQuota_NotExceeded(t *testing.T) {
	st := quotaBaseState()
	_ = SetQuota(st, "staging", 10)
	res, err := CheckQuota(st, "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Exceeded {
		t.Error("expected quota not exceeded")
	}
	if res.Applied != 2 {
		t.Errorf("expected 2 applied, got %d", res.Applied)
	}
}

func TestCheckQuota_Exceeded(t *testing.T) {
	st := quotaBaseState()
	_ = SetQuota(st, "staging", 1)
	res, err := CheckQuota(st, "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Exceeded {
		t.Error("expected quota to be exceeded")
	}
}

func TestRemoveQuota_ClearsLimit(t *testing.T) {
	st := quotaBaseState()
	_ = SetQuota(st, "staging", 5)
	_ = RemoveQuota(st, "staging")
	_, ok := GetQuota(st, "staging")
	if ok {
		t.Error("expected quota to be removed")
	}
}
