package env_test

import (
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/config"
	"github.com/patchwork-deploy/internal/env"
	"github.com/patchwork-deploy/internal/state"
)

func cloneBaseState() *state.State {
	st := state.New()
	now := time.Now()
	st.Add(state.Record{Environment: "staging", Patch: "001-init", AppliedAt: now})
	st.Add(state.Record{Environment: "staging", Patch: "002-users", AppliedAt: now})
	return st
}

func cloneCfg() *config.Config {
	return &config.Config{
		Environments: []config.Environment{
			{Name: "staging", PatchDir: "patches/staging"},
		},
	}
}

func TestClone_CopiesRecordsToNewEnv(t *testing.T) {
	st := cloneBaseState()
	cfg := cloneCfg()

	result, err := env.Clone(cfg, st, "staging", "qa")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Copied != 2 {
		t.Errorf("expected 2 copied, got %d", result.Copied)
	}
	if result.TargetEnv != "qa" {
		t.Errorf("expected target qa, got %s", result.TargetEnv)
	}

	qaRecords := st.ForEnvironment("qa")
	if len(qaRecords) != 2 {
		t.Errorf("expected 2 qa records, got %d", len(qaRecords))
	}
}

func TestClone_SameEnvReturnsError(t *testing.T) {
	st := cloneBaseState()
	cfg := cloneCfg()

	_, err := env.Clone(cfg, st, "staging", "staging")
	if err == nil {
		t.Fatal("expected error for same source and target")
	}
}

func TestClone_MissingSourceReturnsError(t *testing.T) {
	st := cloneBaseState()
	cfg := cloneCfg()

	_, err := env.Clone(cfg, st, "production", "qa")
	if err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestClone_ExistingTargetReturnsError(t *testing.T) {
	st := cloneBaseState()
	cfg := &config.Config{
		Environments: []config.Environment{
			{Name: "staging", PatchDir: "patches/staging"},
			{Name: "qa", PatchDir: "patches/qa"},
		},
	}

	_, err := env.Clone(cfg, st, "staging", "qa")
	if err == nil {
		t.Fatal("expected error when target already exists in config")
	}
}

func TestClone_EmptySourceReturnsZeroCopied(t *testing.T) {
	st := state.New()
	cfg := cloneCfg()

	result, err := env.Clone(cfg, st, "staging", "qa")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Copied != 0 {
		t.Errorf("expected 0 copied, got %d", result.Copied)
	}
}
