package state

import (
	"os"
	"path/filepath"
	"testing"
)

func TestImport_AddsNewRecords(t *testing.T) {
	s := &State{}

	legacy := filepath.Join(t.TempDir(), "applied.txt")
	err := os.WriteFile(legacy, []byte("patch-001\npatch-002\npatch-003\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	n, err := Import(s, "staging", legacy)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 3 {
		t.Errorf("expected 3 imported, got %d", n)
	}
	if !s.Has("staging", "patch-001") {
		t.Error("expected patch-001 to be present")
	}
}

func TestImport_SkipsDuplicates(t *testing.T) {
	s := &State{
		Records: []Record{
			{Environment: "staging", PatchID: "patch-001"},
		},
	}

	legacy := filepath.Join(t.TempDir(), "applied.txt")
	_ = os.WriteFile(legacy, []byte("patch-001\npatch-002\n"), 0644)

	n, err := Import(s, "staging", legacy)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 imported (skip duplicate), got %d", n)
	}
}

func TestImport_MissingFile(t *testing.T) {
	s := &State{}
	_, err := Import(s, "staging", "/nonexistent/file.txt")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestPrune_RemovesStaleRecords(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "patch-001.sh"), []byte("#!/bin/sh"), 0644)
	// patch-002 intentionally absent

	s := &State{
		Records: []Record{
			{Environment: "prod", PatchID: "patch-001"},
			{Environment: "prod", PatchID: "patch-002"},
		},
	}

	n, err := Prune(s, "prod", dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 pruned, got %d", n)
	}
	if s.Has("prod", "patch-002") {
		t.Error("expected patch-002 to be pruned")
	}
	if !s.Has("prod", "patch-001") {
		t.Error("expected patch-001 to remain")
	}
}

func TestPrune_LeavesOtherEnvsUntouched(t *testing.T) {
	dir := t.TempDir()
	// no patches on disk — all prod records would be pruned

	s := &State{
		Records: []Record{
			{Environment: "prod", PatchID: "patch-001"},
			{Environment: "staging", PatchID: "patch-001"},
		},
	}

	_, err := Prune(s, "prod", dir)
	if err != nil {
		t.Fatal(err)
	}
	if !s.Has("staging", "patch-001") {
		t.Error("staging record should not be pruned")
	}
}
