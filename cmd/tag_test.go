package cmd_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/patchwork-deploy/internal/state"
)

func tempTagDir(t *testing.T) string {
	t.Helper()
	d := t.TempDir()
	t.Setenv("PATCHWORK_DIR", d)
	return d
}

func writeTagState(t *testing.T, dir string) {
	t.Helper()
	st := state.New()
	st.Upsert(state.Record{Environment: "prod", Patch: "001-init"})
	st.Upsert(state.Record{Environment: "prod", Patch: "002-schema"})
	if err := state.Save(st, dir); err != nil {
		t.Fatal(err)
	}
}

func TestTagAdd_And_List(t *testing.T) {
	dir := tempTagDir(t)
	writeTagState(t, dir)

	out := executeCommand(t, "tag", "add", "prod", "001-init", "baseline")
	if strings.Contains(out, "error") {
		t.Fatalf("unexpected error: %s", out)
	}

	out = executeCommand(t, "tag", "list", "prod")
	if !strings.Contains(out, "baseline") {
		t.Errorf("expected baseline in output, got: %s", out)
	}
}

func TestTagRemove_ClearsTag(t *testing.T) {
	dir := tempTagDir(t)
	writeTagState(t, dir)

	_ = executeCommand(t, "tag", "add", "prod", "001-init", "v1")
	_ = executeCommand(t, "tag", "remove", "prod", "001-init")

	out := executeCommand(t, "tag", "list", "prod")
	if strings.Contains(out, "v1") {
		t.Errorf("expected tag removed, got: %s", out)
	}
}

func TestTagList_NoTags(t *testing.T) {
	dir := tempTagDir(t)
	writeTagState(t, dir)

	out := executeCommand(t, "tag", "list", "prod")
	if !strings.Contains(out, "no tags") {
		t.Errorf("expected 'no tags' message, got: %s", out)
	}
}

func TestTagAdd_InvalidTag(t *testing.T) {
	dir := tempTagDir(t)
	writeTagState(t, dir)

	out := executeCommand(t, "tag", "add", "prod", "001-init", "bad tag!")
	if !strings.Contains(out, "invalid") {
		t.Errorf("expected invalid tag error, got: %s", out)
	}
	_ = os.Getenv(filepath.Join(dir, "state.json")) // satisfy import
}
