package state

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTakeSnapshot_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	snapshotDir := filepath.Join(dir, "snapshots")

	s := &State{}
	s.records = []Record{
		{Environment: "staging", Patch: "001-init", AppliedAt: time.Now().UTC()},
		{Environment: "staging", Patch: "002-add-index", AppliedAt: time.Now().UTC()},
		{Environment: "prod", Patch: "001-init", AppliedAt: time.Now().UTC()},
	}

	path, err := TakeSnapshot(s, "staging", snapshotDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("snapshot file not found: %v", err)
	}
}

func TestTakeSnapshot_ContainsOnlyEnvRecords(t *testing.T) {
	dir := t.TempDir()

	s := &State{}
	s.records = []Record{
		{Environment: "staging", Patch: "001-init", AppliedAt: time.Now().UTC()},
		{Environment: "prod", Patch: "001-init", AppliedAt: time.Now().UTC()},
	}

	path, err := TakeSnapshot(s, "staging", dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	snap, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("load snapshot: %v", err)
	}

	if snap.Environment != "staging" {
		t.Errorf("expected environment staging, got %s", snap.Environment)
	}
	if len(snap.Records) != 1 {
		t.Errorf("expected 1 record, got %d", len(snap.Records))
	}
	if snap.Records[0].Patch != "001-init" {
		t.Errorf("unexpected patch name: %s", snap.Records[0].Patch)
	}
}

func TestLoadSnapshot_MissingFile(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/path.snapshot.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestTakeSnapshot_EmptyEnv(t *testing.T) {
	dir := t.TempDir()

	s := &State{}
	s.records = []Record{
		{Environment: "prod", Patch: "001-init", AppliedAt: time.Now().UTC()},
	}

	path, err := TakeSnapshot(s, "staging", dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	snap, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("load snapshot: %v", err)
	}

	if len(snap.Records) != 0 {
		t.Errorf("expected 0 records for empty env, got %d", len(snap.Records))
	}
}
