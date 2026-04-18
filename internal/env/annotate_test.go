package env

import (
	"testing"

	"github.com/patchwork-deploy/internal/state"
)

func annotateBaseState() *state.State {
	st := state.New()
	st.Add(state.Record{Environment: "staging", Patch: "001-init"})
	st.Add(state.Record{Environment: "staging", Patch: "002-schema"})
	return st
}

func TestSetAnnotation_Success(t *testing.T) {
	st := annotateBaseState()
	if err := SetAnnotation(st, "staging", "001-init", "bootstraps the db"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := GetAnnotation(st, "staging", "001-init")
	if !ok || v != "bootstraps the db" {
		t.Fatalf("expected annotation, got %q %v", v, ok)
	}
}

func TestSetAnnotation_MissingEnvReturnsError(t *testing.T) {
	st := annotateBaseState()
	if err := SetAnnotation(st, "prod", "001-init", "note"); err == nil {
		t.Fatal("expected error for missing env")
	}
}

func TestSetAnnotation_MissingPatchReturnsError(t *testing.T) {
	st := annotateBaseState()
	if err := SetAnnotation(st, "staging", "999-nope", "note"); err == nil {
		t.Fatal("expected error for missing patch")
	}
}

func TestSetAnnotation_NewlineRejected(t *testing.T) {
	st := annotateBaseState()
	if err := SetAnnotation(st, "staging", "001-init", "line1\nline2"); err == nil {
		t.Fatal("expected error for newline in text")
	}
}

func TestRemoveAnnotation_ClearsEntry(t *testing.T) {
	st := annotateBaseState()
	_ = SetAnnotation(st, "staging", "001-init", "temp note")
	RemoveAnnotation(st, "staging", "001-init")
	_, ok := GetAnnotation(st, "staging", "001-init")
	if ok {
		t.Fatal("expected annotation to be removed")
	}
}

func TestListAnnotations_ReturnsAll(t *testing.T) {
	st := annotateBaseState()
	_ = SetAnnotation(st, "staging", "001-init", "first")
	_ = SetAnnotation(st, "staging", "002-schema", "second")
	list := ListAnnotations(st, "staging")
	if len(list) != 2 {
		t.Fatalf("expected 2 annotations, got %d", len(list))
	}
	if list["001-init"] != "first" || list["002-schema"] != "second" {
		t.Fatalf("unexpected annotation values: %v", list)
	}
}
