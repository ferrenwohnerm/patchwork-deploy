package env

import (
	"fmt"

	"github.com/patchwork-deploy/internal/state"
)

// FreezeResult holds the outcome of a freeze or unfreeze operation.
type FreezeResult struct {
	Environment string
	Frozen      bool
}

const freezePrefix = "__frozen__"

// Freeze marks an environment as frozen by adding a sentinel record.
// Frozen environments are skipped during patch application.
func Freeze(st *state.State, env string) (FreezeResult, error) {
	records := st.ForEnvironment(env)
	if records == nil {
		return FreezeResult{}, fmt.Errorf("environment %q not found", env)
	}
	for _, r := range records {
		if r.Patch == freezePrefix {
			return FreezeResult{Environment: env, Frozen: true}, nil
		}
	}
	st.Add(state.Record{
		Environment: env,
		Patch:       freezePrefix,
		AppliedAt:   state.Now(),
	})
	return FreezeResult{Environment: env, Frozen: true}, nil
}

// Unfreeze removes the frozen sentinel from an environment.
func Unfreeze(st *state.State, env string) (FreezeResult, error) {
	records := st.ForEnvironment(env)
	if records == nil {
		return FreezeResult{}, fmt.Errorf("environment %q not found", env)
	}
	st.RemoveWhere(func(r state.Record) bool {
		return r.Environment == env && r.Patch == freezePrefix
	})
	return FreezeResult{Environment: env, Frozen: false}, nil
}

// IsFrozen reports whether the given environment is currently frozen.
func IsFrozen(st *state.State, env string) bool {
	for _, r := range st.ForEnvironment(env) {
		if r.Patch == freezePrefix {
			return true
		}
	}
	return false
}
