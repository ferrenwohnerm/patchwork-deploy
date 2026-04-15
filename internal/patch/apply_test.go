package patch

import (
	"os"
	"path/filepath"
	"testing"
)

func TestApply_Success(t *testing.T) {
	patchDir := makeTempPatchDir(t, []string{"001_init.yaml", "002_add_service.yaml"})
	loader := NewLoader(patchDir)
	applier := NewApplier(loader)

	destDir := t.TempDir()
	if err := applier.Apply("001_init.yaml", "staging", destDir); err != nil {
		t.Fatalf("Apply() unexpected error: %v", err)
	}

	dest := filepath.Join(destDir, "001_init.yaml")
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		t.Errorf("expected applied patch file to exist at %s", dest)
	}
}

func TestApply_PatchNotFound(t *testing.T) {
	patchDir := makeTempPatchDir(t, []string{"001_init.yaml"})
	loader := NewLoader(patchDir)
	applier := NewApplier(loader)

	err := applier.Apply("999_missing.yaml", "staging", t.TempDir())
	if err == nil {
		t.Fatal("expected error for missing patch, got nil")
	}
}

func TestWasApplied(t *testing.T) {
	patchDir := makeTempPatchDir(t, []string{"001_init.yaml"})
	loader := NewLoader(patchDir)
	applier := NewApplier(loader)

	if applier.WasApplied("001_init.yaml", "production") {
		t.Error("expected WasApplied to return false before applying")
	}

	if err := applier.Apply("001_init.yaml", "production", t.TempDir()); err != nil {
		t.Fatalf("Apply() unexpected error: %v", err)
	}

	if !applier.WasApplied("001_init.yaml", "production") {
		t.Error("expected WasApplied to return true after applying")
	}

	if applier.WasApplied("001_init.yaml", "staging") {
		t.Error("WasApplied should be false for a different environment")
	}
}

func TestApplied_ReturnsAllRecords(t *testing.T) {
	patchDir := makeTempPatchDir(t, []string{"001_init.yaml", "002_add_service.yaml"})
	loader := NewLoader(patchDir)
	applier := NewApplier(loader)
	destDir := t.TempDir()

	_ = applier.Apply("001_init.yaml", "staging", destDir)
	_ = applier.Apply("002_add_service.yaml", "staging", destDir)

	records := applier.Applied()
	if len(records) != 2 {
		t.Errorf("expected 2 records, got %d", len(records))
	}
}
