package env

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/patchwork-deploy/internal/config"
)

var validNameRe = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// ValidationResult holds the outcome of validating an environment.
type ValidationResult struct {
	Environment string
	Errors      []string
}

// Valid returns true if no errors were found.
func (v ValidationResult) Valid() bool {
	return len(v.Errors) == 0
}

// String returns a human-readable summary of the result.
func (v ValidationResult) String() string {
	if v.Valid() {
		return fmt.Sprintf("environment %q: OK", v.Environment)
	}
	return fmt.Sprintf("environment %q: %s", v.Environment, strings.Join(v.Errors, "; "))
}

// Validate checks a single environment entry from the config for
// structural correctness and returns a ValidationResult.
func Validate(env config.Environment) ValidationResult {
	result := ValidationResult{Environment: env.Name}

	if strings.TrimSpace(env.Name) == "" {
		result.Errors = append(result.Errors, "name must not be empty")
	} else if !validNameRe.MatchString(env.Name) {
		result.Errors = append(result.Errors,
			"name must contain only alphanumeric characters, hyphens, or underscores")
	}

	if strings.TrimSpace(env.PatchDir) == "" {
		result.Errors = append(result.Errors, "patch_dir must not be empty")
	}

	if strings.TrimSpace(env.StateFile) == "" {
		result.Errors = append(result.Errors, "state_file must not be empty")
	}

	return result
}

// ValidateAll validates every environment in the config and returns
// all results, including passing ones.
func ValidateAll(cfg *config.Config) []ValidationResult {
	results := make([]ValidationResult, 0, len(cfg.Environments))
	for _, env := range cfg.Environments {
		results = append(results, Validate(env))
	}
	return results
}
