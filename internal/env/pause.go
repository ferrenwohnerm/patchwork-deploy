package env

import (
	"fmt"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

const pauseSentinel = "__paused__"

// Pause marks an environment as paused, preventing new patches from being applied.
func Pause(st *state.State, env string) error {
	recs := st.ForEnvironment(env)
	if len(recs) == 0 {
		return fmt.Errorf("environment %q not found", env)
	}
	if IsPaused(st, env) {
		return fmt.Errorf("environment %q is already paused", env)
	}
	st.Add(state.Record{
		Environment: env,
		Patch:       pauseSentinel,
		AppliedAt:   time.Now(),
	})
	return nil
}

// Unpause removes the paused sentinel from an environment.
func Unpause(st *state.State, env string) error {
	recs := st.ForEnvironment(env)
	if len(recs) == 0 {
		return fmt.Errorf("environment %q not found", env)
	}
	if !IsPaused(st, env) {
		return fmt.Errorf("environment %q is not paused", env)
	}
	st.RemoveWhere(func(r state.Record) bool {
		return r.Environment == env && r.Patch == pauseSentinel
	})
	return nil
}

// IsPaused reports whether the environment is currently paused.
func IsPaused(st *state.State, env string) bool {
	for _, r := range st.ForEnvironment(env) {
		if r.Patch == pauseSentinel {
			return true
		}
	}
	return false
}
