package env

import (
	"fmt"

	"github.com/patchwork-deploy/internal/config"
	"github.com/patchwork-deploy/internal/state"
)

// CloneResult holds the outcome of a clone operation.
type CloneResult struct {
	SourceEnv string
	TargetEnv string
	Copied    int
}

// Clone duplicates all state records from one environment into a new target
// environment. The target must not already exist in the config.
func Clone(cfg *config.Config, st *state.State, source, target string) (CloneResult, error) {
	if source == target {
		return CloneResult{}, fmt.Errorf("source and target environment must differ: %q", source)
	}

	if _, err := cfg.GetEnvironment(source); err != nil {
		return CloneResult{}, fmt.Errorf("source environment not found: %w", err)
	}

	if _, err := cfg.GetEnvironment(target); err == nil {
		return CloneResult{}, fmt.Errorf("target environment %q already exists in config", target)
	}

	records := st.ForEnvironment(source)
	if len(records) == 0 {
		return CloneResult{Source: source, Target: target}, nil
	}

	for _, r := range records {
		cloned := state.Record{
			Environment: target,
			Patch:       r.Patch,
			AppliedAt:   r.AppliedAt,
		}
		st.Add(cloned)
	}

	return CloneResult{
		SourceEnv: source,
		TargetEnv: target,
		Copied:    len(records),
	}, nil
}
