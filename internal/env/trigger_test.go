package env

import (
	"testing"

	"github.com/patchwork-deploy/internal/state"
)

func triggerBaseState() *state.State {
	st := state.New()
	st.AddRecord("prod", state.Record{Patch: "001-init"})
	st.AddRecord("prod", state.Record{Patch: "002-schema"})
	return st
}

func TestSetTrigger_Success(t *testing.T) {
	st := triggerBaseState()
	if err := SetTrigger(st, "prod", "001-init", "post-apply", "notify-slack"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	triggers, _ := ListTriggers(st, "prod")
	if len(triggers) != 1 {
		t.Fatalf("expected 1 trigger, got %d", len(triggers))
	}
	if triggers[0].Event != "post-apply" || triggers[0].Action != "notify-slack" {
		t.Errorf("unexpected trigger: %+v", triggers[0])
	}
}

func TestSetTrigger_MissingEnvReturnsError(t *testing.T) {
	st := triggerBaseState()
	if err := SetTrigger(st, "staging", "001-init", "post-apply", "notify"); err == nil {
		t.Fatal("expected error for missing env")
	}
}

func TestSetTrigger_MissingPatchReturnsError(t *testing.T) {
	st := triggerBaseState()
	if err := SetTrigger(st, "prod", "999-nope", "post-apply", "notify"); err == nil {
		t.Fatal("expected error for missing patch")
	}
}

func TestSetTrigger_EmptyEventReturnsError(t *testing.T) {
	st := triggerBaseState()
	if err := SetTrigger(st, "prod", "001-init", "", "notify"); err == nil {
		t.Fatal("expected error for empty event")
	}
}

func TestRemoveTrigger_ClearsEntry(t *testing.T) {
	st := triggerBaseState()
	_ = SetTrigger(st, "prod", "001-init", "post-apply", "notify-slack")
	_ = RemoveTrigger(st, "prod", "001-init", "post-apply")
	triggers, _ := ListTriggers(st, "prod")
	if len(triggers) != 0 {
		t.Errorf("expected 0 triggers after removal, got %d", len(triggers))
	}
}

func TestListTriggers_MissingEnvReturnsError(t *testing.T) {
	st := triggerBaseState()
	_, err := ListTriggers(st, "ghost")
	if err == nil {
		t.Fatal("expected error for missing env")
	}
}
