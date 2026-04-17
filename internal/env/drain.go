package env

import (
	"fmt"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

// DrainResult holds the outcome of a drain operation.
type DrainResult struct {
	Environment string
	Removed     []string
	DryRun      bool
}

// Drain marks an environment as draining and removes all pending (unapplied)
// scheduled entries, returning the patches that were cleared.
func Drain(st *state.State, env string, dryRun bool) (DrainResult, error) {
	records := st.ForEnvironment(env)
	if len(records) == 0 {
		return DrainResult{}, fmt.Errorf("environment %q not found", env)
	}

	scheduledKey := "__scheduled__" + env
	scheduled := st.ForEnvironment(scheduledKey)

	result := DrainResult{
		Environment: env,
		DryRun:      dryRun,
	}

	for _, rec := range scheduled {
		result.Removed = append(result.Removed, rec.Patch)
	}

	if !dryRun {
		st.RemoveEnvironment(scheduledKey)
		st.Add(state.Record{
			Environment: "__drain__" + env,
			Patch:       "drain",
			AppliedAt:   time.Now(),
		})
	}

	return result, nil
}

// IsDrained reports whether the environment has been drained.
func IsDrained(st *state.State, env string) bool {
	return len(st.ForEnvironment("__drain__"+env)) > 0
}

// Undrain clears the drain sentinel for an environment.
func Undrain(st *state.State, env string) error {
	key := "__drain__" + env
	if len(st.ForEnvironment(key)) == 0 {
		return nil
	}
	st.RemoveEnvironment(key)
	return nil
}
