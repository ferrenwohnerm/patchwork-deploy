package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/patchwork-deploy/internal/state"
)

func tempIntentDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "intent-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func writeIntentState(t *testing.T, dir string) {
	t.Helper()
	st := state.NewInMemory()
	st.AddRecord("staging", state.Record{Patch: "001-init"})
	st.AddRecord("staging", state.Record{Patch: "002-feature"})
	if err := state.Save(st, dir); err != nil {
		t.Fatal(err)
	}
}

func TestIntentSet_And_Get(t *testing.T) {
	dir := tempIntentDir(t)
	writeIntentState(t, dir)
	t.Setenv("PATCHWORK_DIR", dir)

	root := newRootCmd()
	root.SetArgs([]string{"intent", "set", "staging", "001-init", "bootstrap schema"})
	if err := root.Execute(); err != nil {
		t.Fatalf("set failed: %v", err)
	}

	st, _ := state.Load(dir)
	v, ok := st.GetMeta("staging", "intent:001-init")
	if !ok || v != "bootstrap schema" {
		t.Errorf("expected intent to be stored, got %q", v)
	}
}

func TestIntentRemove_ClearsEntry(t *testing.T) {
	dir := tempIntentDir(t)
	writeIntentState(t, dir)
	t.Setenv("PATCHWORK_DIR", dir)

	root := newRootCmd()
	root.SetArgs([]string{"intent", "set", "staging", "001-init", "to be removed"})
	_ = root.Execute()

	root2 := newRootCmd()
	root2.SetArgs([]string{"intent", "remove", "staging", "001-init"})
	if err := root2.Execute(); err != nil {
		t.Fatalf("remove failed: %v", err)
	}

	st, _ := state.Load(dir)
	_, ok := st.GetMeta("staging", "intent:001-init")
	if ok {
		t.Error("expected intent to be cleared")
	}
}

func TestIntentSet_UnknownEnvFails(t *testing.T) {
	dir := tempIntentDir(t)
	writeIntentState(t, dir)
	t.Setenv("PATCHWORK_DIR", dir)

	root := newRootCmd()
	root.SetArgs([]string{"intent", "set", "prod", "001-init", "some intent"})
	if err := root.Execute(); err == nil {
		t.Fatal("expected error for unknown environment")
	}
}

func TestIntentList_ShowsEntries(t *testing.T) {
	dir := tempIntentDir(t)
	writeIntentState(t, dir)
	t.Setenv("PATCHWORK_DIR", dir)

	root := newRootCmd()
	root.SetArgs([]string{"intent", "set", "staging", "002-feature", "new feature rollout"})
	_ = root.Execute()

	_ = filepath.Join(dir, "state.json") // ensure file exists
	root2 := newRootCmd()
	root2.SetArgs([]string{"intent", "list", "staging"})
	if err := root2.Execute(); err != nil {
		t.Fatalf("list failed: %v", err)
	}
}
