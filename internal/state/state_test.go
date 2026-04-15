package state

import (
	"os"
	"path/filepath"
	"testing"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "patchwork-state-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestLoad_EmptyWhenMissing(t *testing.T) {
	dir := tempDir(t)
	s, err := Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Records) != 0 {
		t.Errorf("expected 0 records, got %d", len(s.Records))
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := tempDir(t)
	s, _ := Load(dir)
	s.Add("001-init", "staging")
	s.Add("002-feature", "staging")

	if err := s.Save(); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := Load(dir)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	if len(loaded.Records) != 2 {
		t.Errorf("expected 2 records, got %d", len(loaded.Records))
	}
}

func TestHas(t *testing.T) {
	dir := tempDir(t)
	s, _ := Load(dir)
	s.Add("001-init", "production")

	if !s.Has("001-init", "production") {
		t.Error("expected Has to return true")
	}
	if s.Has("001-init", "staging") {
		t.Error("expected Has to return false for different env")
	}
	if s.Has("002-other", "production") {
		t.Error("expected Has to return false for unknown patch")
	}
}

func TestForEnvironment(t *testing.T) {
	dir := tempDir(t)
	s, _ := Load(dir)
	s.Add("001-init", "staging")
	s.Add("002-feature", "production")
	s.Add("003-fix", "staging")

	records := s.ForEnvironment("staging")
	if len(records) != 2 {
		t.Errorf("expected 2 staging records, got %d", len(records))
	}
}

func TestStateFilePath(t *testing.T) {
	dir := tempDir(t)
	s, _ := Load(dir)
	s.Add("001-init", "dev")
	_ = s.Save()

	expected := filepath.Join(dir, stateFileName)
	if _, err := os.Stat(expected); os.IsNotExist(err) {
		t.Errorf("state file not found at %s", expected)
	}
}
