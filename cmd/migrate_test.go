package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMigrateCmd_Import(t *testing.T) {
	dir := t.TempDir()

	cfgPath := filepath.Join(dir, "patchwork.yaml")
	patchDir := filepath.Join(dir, "patches")
	_ = os.MkdirAll(patchDir, 0755)

	cfgContent := fmt.Sprintf(`
state_file: %s/state.json
environments:
  - name: staging
    patch_dir: %s
`, dir, patchDir)
	_ = os.WriteFile(cfgPath, []byte(cfgContent), 0644)

	legacy := filepath.Join(dir, "legacy.txt")
	_ = os.WriteFile(legacy, []byte("patch-001\npatch-002\n"), 0644)

	args := []string{
		"--config", cfgPath,
		"migrate",
		"--env", "staging",
		"--import", legacy,
	}

	out, err := executeCommand(args...)
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "Imported 2") {
		t.Errorf("expected import count in output, got: %s", out)
	}
}

func TestMigrateCmd_Prune(t *testing.T) {
	dir := t.TempDir()
	patchDir := filepath.Join(dir, "patches")
	_ = os.MkdirAll(patchDir, 0755)
	// only patch-001 exists on disk
	_ = os.WriteFile(filepath.Join(patchDir, "patch-001.sh"), []byte("#!/bin/sh"), 0644)

	cfgPath := filepath.Join(dir, "patchwork.yaml")
	cfgContent := fmt.Sprintf(`
state_file: %s/state.json
environments:
  - name: prod
    patch_dir: %s
`, dir, patchDir)
	_ = os.WriteFile(cfgPath, []byte(cfgContent), 0644)

	// seed state with two records
	statePath := filepath.Join(dir, "state.json")
	_ = os.WriteFile(statePath, []byte(`{"records":[{"environment":"prod","patch_id":"patch-001","applied_at":""},{"environment":"prod","patch_id":"patch-stale","applied_at":""}]}`), 0644)

	args := []string{
		"--config", cfgPath,
		"migrate",
		"--env", "prod",
		"--prune",
	}

	out, err := executeCommand(args...)
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "Pruned 1") {
		t.Errorf("expected prune count in output, got: %s", out)
	}
}
