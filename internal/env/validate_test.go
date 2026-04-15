package env

import (
	"testing"

	"github.com/patchwork-deploy/internal/config"
)

func makeEnvCfg(name, patchDir, stateFile string) config.Environment {
	return config.Environment{
		Name:      name,
		PatchDir:  patchDir,
		StateFile: stateFile,
	}
}

func TestValidate_Valid(t *testing.T) {
	result := Validate(makeEnvCfg("production", "./patches/prod", "state/prod.json"))
	if !result.Valid() {
		t.Fatalf("expected valid, got errors: %v", result.Errors)
	}
}

func TestValidate_EmptyName(t *testing.T) {
	result := Validate(makeEnvCfg("", "./patches", "state.json"))
	if result.Valid() {
		t.Fatal("expected invalid due to empty name")
	}
	if len(result.Errors) == 0 || result.Errors[0] != "name must not be empty" {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
}

func TestValidate_InvalidNameChars(t *testing.T) {
	result := Validate(makeEnvCfg("prod env!", "./patches", "state.json"))
	if result.Valid() {
		t.Fatal("expected invalid due to bad name characters")
	}
}

func TestValidate_EmptyPatchDir(t *testing.T) {
	result := Validate(makeEnvCfg("staging", "", "state.json"))
	if result.Valid() {
		t.Fatal("expected invalid due to empty patch_dir")
	}
}

func TestValidate_EmptyStateFile(t *testing.T) {
	result := Validate(makeEnvCfg("staging", "./patches", ""))
	if result.Valid() {
		t.Fatal("expected invalid due to empty state_file")
	}
}

func TestValidate_String_OK(t *testing.T) {
	result := Validate(makeEnvCfg("dev", "./patches", "state.json"))
	s := result.String()
	if s != `environment "dev": OK` {
		t.Fatalf("unexpected string: %s", s)
	}
}

func TestValidateAll_MixedResults(t *testing.T) {
	cfg := &config.Config{
		Environments: []config.Environment{
			makeEnvCfg("prod", "./patches/prod", "state/prod.json"),
			makeEnvCfg("", "./patches/bad", "state/bad.json"),
			makeEnvCfg("staging", "", "state/staging.json"),
		},
	}

	results := ValidateAll(cfg)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if !results[0].Valid() {
		t.Error("expected first result to be valid")
	}
	if results[1].Valid() {
		t.Error("expected second result to be invalid")
	}
	if results[2].Valid() {
		t.Error("expected third result to be invalid")
	}
}
