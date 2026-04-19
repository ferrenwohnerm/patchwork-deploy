package env

import (
	"fmt"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

const sealPrefix = "__sealed__"

// Seal marks an environment as sealed, preventing any further patch applications.
// An optional reason may be provided.
func Seal(st *state.State, env, reason string) error {
	records := st.ForEnvironment(env)
	if len(records) == 0 {
		return fmt.Errorf("environment %q not found", env)
	}
	if IsSealed(st, env) {
		return nil
	}
	st.Add(state.Record{
		Environment: env,
		Patch:       sealPrefix,
		AppliedAt:   time.Now().UTC(),
		Note:        reason,
	})
	return nil
}

// Unseal removes the sealed sentinel from an environment.
func Unseal(st *state.State, env string) error {
	records := st.ForEnvironment(env)
	if len(records) == 0 {
		return fmt.Errorf("environment %q not found", env)
	}
	st.RemoveWhere(func(r state.Record) bool {
		return r.Environment == env && r.Patch == sealPrefix
	})
	return nil
}

// IsSealed reports whether the environment is currently sealed.
func IsSealed(st *state.State, env string) bool {
	for _, r := range st.ForEnvironment(env) {
		if r.Patch == sealPrefix {
			return true
		}
	}
	return false
}

// SealReason returns the reason stored when the environment was sealed,
// or an empty string if none was provided or the env is not sealed.
func SealReason(st *state.State, env string) string {
	for _, r := range st.ForEnvironment(env) {
		if r.Patch == sealPrefix {
			return r.Note
		}
	}
	return ""
}
