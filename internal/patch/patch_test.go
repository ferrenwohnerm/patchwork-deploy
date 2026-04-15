package patch_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/patchwork-deploy/internal/patch"
)

func makeTempPatchDir(t *testing.T, files []string) string {
	t.Helper()
	dir := t.TempDir()
	for _, f := range files {
		path := filepath.Join(dir, f)
		if err := os.WriteFile(path, []byte("# patch\n"), 0644); err != nil {
			t.Fatalf("failed to write temp file %s: %v", f, err)
		}
	}
	return dir
}

func TestDiscover_ReturnsSortedPatches(t *testing.T) {
	dir := makeTempPatchDir(t, []string{
		"staging_002_add-feature.yaml",
		"staging_001_init-db.yaml",
		"staging_003_update-config.yaml",
	})

	loader := patch.NewLoader(dir)
	patches, err := loader.Discover("staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(patches) != 3 {
		t.Fatalf("expected 3 patches, got %d", len(patches))
	}

	expected := []string{
		"staging_001_init-db",
		"staging_002_add-feature",
		"staging_003_update-config",
	}
	for i, p := range patches {
		if p.Name != expected[i] {
			t.Errorf("patch[%d]: expected name %q, got %q", i, expected[i], p.Name)
		}
		if p.Environment != "staging" {
			t.Errorf("patch[%d]: expected env %q, got %q", i, "staging", p.Environment)
		}
	}
}

func TestDiscover_NoMatches(t *testing.T) {
	dir := makeTempPatchDir(t, []string{
		"production_001_init.yaml",
	})

	loader := patch.NewLoader(dir)
	patches, err := loader.Discover("staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(patches) != 0 {
		t.Errorf("expected 0 patches, got %d", len(patches))
	}
}

func TestExists(t *testing.T) {
	dir := makeTempPatchDir(t, []string{"staging_001_init.yaml"})

	existing := filepath.Join(dir, "staging_001_init.yaml")
	if !patch.Exists(existing) {
		t.Errorf("expected file to exist: %s", existing)
	}

	missing := filepath.Join(dir, "staging_999_missing.yaml")
	if patch.Exists(missing) {
		t.Errorf("expected file to not exist: %s", missing)
	}
}
