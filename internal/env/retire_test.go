package env

import (
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func retireBaseState() *state.State {
	st := state.New()
	now := time.Now()
	st.Add(state.Record{Environment: "prod", Patch: "001-init", AppliedAt: now})
	st.Add(state.Record{Environment: "prod", Patch: "002-add-table", AppliedAt: now})
	st.Add(state.Record{Environment: "staging", Patch: "001-init", AppliedAt: now})
	return st
}

func TestRetire_MovesRecordsToRetiredKey(t *testing.T) {
	st := retireBaseState()
	if err := Retire(st, "prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(st.ForEnvironment("prod")) != 0 {
		t.Error("expected prod records to be removed")
	}
	retired := st.ForEnvironment("__retired__prod")
	if len(retired) == 0 {
		t.Error("expected retired records to exist")
	}
}

func TestRetire_LeavesOtherEnvsUntouched(t *testing.T) {
	st := retireBaseState()
	_ = Retire(st, "prod")
	if len(st.ForEnvironment("staging")) == 0 {
		t.Error("staging should be unaffected")
	}
}

func TestRetire_MissingEnvReturnsError(t *testing.T) {
	st := retireBaseState()
	if err := Retire(st, "nonexistent"); err == nil {
		t.Error("expected error for missing environment")
	}
}

func TestIsRetired_TrueAfterRetire(t *testing.T) {
	st := retireBaseState()
	_ = Retire(st, "prod")
	if !IsRetired(st, "prod") {
		t.Error("expected prod to be retired")
	}
}

func TestIsRetired_FalseForActiveEnv(t *testing.T) {
	st := retireBaseState()
	if IsRetired(st, "staging") {
		t.Error("staging should not be retired")
	}
}

func TestListRetired_ReturnsRetiredNames(t *testing.T) {
	st := retireBaseState()
	_ = Retire(st, "prod")
	names := ListRetired(st)
	if len(names) != 1 || names[0] != "prod" {
		t.Errorf("expected [prod], got %v", names)
	}
}
