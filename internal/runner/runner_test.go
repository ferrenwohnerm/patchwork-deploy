package runner_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/patchwork-deploy/internal/config"
	"github.com/patchwork-deploy/internal/patch"
	"github.com/patchwork-deploy/internal/runner"
	"github.com/patchwork-deploy/internal/state"
)

func makeEnv(t *testing.T, patchDir string) config.Environment {
	t.Helper()
	return config.Environment{Name: "test", PatchDir: patchDir}
}

func writePatch(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestRun_AppliesNewPatches(t *testing.T) {
	dir := t.TempDir()
	writePatch(t, dir, "001_init.sh", "#!/bin/sh\necho ok")

	env := makeEnv(t, dir)
	loader := patch.NewLoader()
	applier := patch.NewApplier()
	st, _ := state.Load(filepath.Join(dir, "state.json"))

	r := runner.New(env, loader, applier, st)
	results, err := r.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || !results[0].Applied {
		t.Errorf("expected 1 applied result, got %+v", results)
	}
}

func TestRun_SkipsAlreadyApplied(t *testing.T) {
	dir := t.TempDir()
	writePatch(t, dir, "001_init.sh", "#!/bin/sh\necho ok")

	env := makeEnv(t, dir)
	loader := patch.NewLoader()
	applier := patch.NewApplier()
	st, _ := state.Load(filepath.Join(dir, "state.json"))
	st.Record("test", "001_init.sh")

	r := runner.New(env, loader, applier, st)
	results, err := r.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || !results[0].Skipped {
		t.Errorf("expected 1 skipped result, got %+v", results)
	}
}

func TestRun_EmptyPatchDir(t *testing.T) {
	dir := t.TempDir()
	env := makeEnv(t, dir)
	loader := patch.NewLoader()
	applier := patch.NewApplier()
	st, _ := state.Load(filepath.Join(dir, "state.json"))

	r := runner.New(env, loader, applier, st)
	results, err := r.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no results, got %d", len(results))
	}
}
