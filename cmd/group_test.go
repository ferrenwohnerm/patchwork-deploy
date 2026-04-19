package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/patchwork-deploy/internal/state"
)

func tempGroupDir(t *testing.T) string {
	t.Helper()
	d := t.TempDir()
	return d
}

func writeGroupState(t *testing.T, dir string, st *state.State) {
	t.Helper()
	b, _ := json.Marshal(st)
	_ = os.WriteFile(filepath.Join(dir, "state.json"), b, 0644)
}

func TestGroupAdd_And_List(t *testing.T) {
	dir := tempGroupDir(t)
	st := state.New()
	st.AddEnvironment("staging")
	st.AddEnvironment("production")
	writeGroupState(t, dir, st)
	writeConfig(t, dir)

	out, err := executeCmd(t, dir, "group", "add", "web", "staging")
	if err != nil {
		t.Fatalf("add failed: %v", err)
	}
	if !strings.Contains(out, "staging") {
		t.Errorf("expected staging in output, got %q", out)
	}

	out, err = executeCmd(t, dir, "group", "list", "web")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if !strings.Contains(out, "staging") {
		t.Errorf("expected staging in list output, got %q", out)
	}
}

func TestGroupRemove_RemovesMember(t *testing.T) {
	dir := tempGroupDir(t)
	st := state.New()
	st.AddEnvironment("staging")
	writeGroupState(t, dir, st)
	writeConfig(t, dir)

	_, _ = executeCmd(t, dir, "group", "add", "web", "staging")
	_, err := executeCmd(t, dir, "group", "remove", "web", "staging")
	if err != nil {
		t.Fatalf("remove failed: %v", err)
	}

	_, err = executeCmd(t, dir, "group", "list", "web")
	if err == nil {
		t.Error("expected error listing empty/removed group")
	}
}

func TestGroupAdd_UnknownEnvFails(t *testing.T) {
	dir := tempGroupDir(t)
	st := state.New()
	writeGroupState(t, dir, st)
	writeConfig(t, dir)

	_, err := executeCmd(t, dir, "group", "add", "web", "ghost")
	if err == nil {
		t.Error("expected error for unknown environment")
	}
}
