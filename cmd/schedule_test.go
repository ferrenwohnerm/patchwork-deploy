package cmd

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

func tempScheduleDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "schedule-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func writeScheduleState(t *testing.T, dir string, st *state.State) {
	t.Helper()
	if err := state.Save(dir, st); err != nil {
		t.Fatal(err)
	}
}

func TestScheduleAdd_And_List(t *testing.T) {
	dir := tempScheduleDir(t)
	st := state.NewInMemory()
	now := time.Now()
	st.Add("prod", state.Record{Environment: "prod", Patch: "001-init.sql", AppliedAt: now})
	writeScheduleState(t, dir, st)

	runAt := time.Now().Add(48 * time.Hour).UTC().Format(time.RFC3339)
	cfgPath := writeTempConfig(t, dir, []string{"prod"})

	root := newRootCmd()
	root.SetArgs([]string{"--config", cfgPath, "schedule", "add", "prod", "001-init.sql", runAt})
	if err := root.Execute(); err != nil {
		t.Fatalf("schedule add failed: %v", err)
	}

	root2 := newRootCmd()
	root2.SetArgs([]string{"--config", cfgPath, "schedule", "list", "prod"})
	if err := root2.Execute(); err != nil {
		t.Fatalf("schedule list failed: %v", err)
	}
}

func TestScheduleCancel_RemovesEntry(t *testing.T) {
	dir := tempScheduleDir(t)
	st := state.NewInMemory()
	now := time.Now()
	st.Add("prod", state.Record{Environment: "prod", Patch: "001-init.sql", AppliedAt: now})
	writeScheduleState(t, dir, st)

	runAt := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)
	cfgPath := writeTempConfig(t, dir, []string{"prod"})

	for _, args := range [][]string{
		{"--config", cfgPath, "schedule", "add", "prod", "001-init.sql", runAt},
		{"--config", cfgPath, "schedule", "cancel", "prod", "001-init.sql"},
	} {
		r := newRootCmd()
		r.SetArgs(args)
		if err := r.Execute(); err != nil {
			t.Fatalf("command %v failed: %v", args, err)
		}
	}
}

func TestScheduleAdd_UnknownEnvFails(t *testing.T) {
	dir := tempScheduleDir(t)
	st := state.NewInMemory()
	writeScheduleState(t, dir, st)
	cfgPath := writeTempConfig(t, dir, []string{"prod"})

	runAt := time.Now().Add(time.Hour).UTC().Format(time.RFC3339)
	r := newRootCmd()
	r.SetArgs([]string{"--config", cfgPath, "schedule", "add", "staging", "001.sql", runAt})
	if err := r.Execute(); err == nil {
		t.Fatal("expected error for unknown environment")
	}
}

func writeTempConfig(t *testing.T, dir string, envs []string) string {
	t.Helper()
	lines := "environments:\n"
	for _, e := range envs {
		lines += fmt.Sprintf("  - name: %s\n    patch_dir: %s\n", e, filepath.Join(dir, e))
	}
	p := filepath.Join(dir, "patchwork.yaml")
	if err := os.WriteFile(p, []byte(lines), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}
