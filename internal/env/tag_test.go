package env_test

import (
	"testing"

	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
)

func tagBaseState() *state.State {
	st := state.New()
	st.Upsert(state.Record{Environment: "prod", Patch: "001-init"})
	st.Upsert(state.Record{Environment: "prod", Patch: "002-schema"})
	return st
}

func TestAddTag_Success(t *testing.T) {
	st := tagBaseState()
	if err := env.AddTag(st, "prod", "001-init", "baseline"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tags, _ := env.ListTags(st, "prod")
	if len(tags) != 1 || tags[0].Tag != "baseline" {
		t.Errorf("expected tag baseline, got %+v", tags)
	}
}

func TestAddTag_InvalidChars(t *testing.T) {
	st := tagBaseState()
	if err := env.AddTag(st, "prod", "001-init", "bad tag!"); err == nil {
		t.Error("expected error for invalid tag")
	}
}

func TestAddTag_MissingEnv(t *testing.T) {
	st := tagBaseState()
	if err := env.AddTag(st, "staging", "001-init", "v1"); err == nil {
		t.Error("expected error for missing env")
	}
}

func TestAddTag_MissingPatch(t *testing.T) {
	st := tagBaseState()
	if err := env.AddTag(st, "prod", "999-nope", "v1"); err == nil {
		t.Error("expected error for missing patch")
	}
}

func TestRemoveTag_ClearsTag(t *testing.T) {
	st := tagBaseState()
	_ = env.AddTag(st, "prod", "001-init", "baseline")
	if err := env.RemoveTag(st, "prod", "001-init"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tags, _ := env.ListTags(st, "prod")
	if len(tags) != 0 {
		t.Errorf("expected no tags, got %+v", tags)
	}
}

func TestListTags_MissingEnv(t *testing.T) {
	st := tagBaseState()
	_, err := env.ListTags(st, "ghost")
	if err == nil {
		t.Error("expected error for missing env")
	}
}
