package env

import (
	"fmt"

	"github.com/patchwork-deploy/internal/state"
)

// ResetResult holds the outcome of a reset operation.
type ResetResult struct {
	Environment string
	Removed     int
}

// Reset removes all state records for the given environment.
// If dryRun is true, the state file is not modified.
func Reset(st *state.State, env string, dryRun bool) (ResetResult, error) {
	records := st.ForEnvironment(env)
	if len(records) == 0 {
		return ResetResult{}, fmt.Errorf("environment %q not found or has no records", env)
	}

	result := ResetResult{
		Environment: env,
		Removed:     len(records),
	}

	if dryRun {
		return result, nil
	}

	for _, r := range records {
		st.Remove(env, r.Patch)
	}

	return result, nil
}
