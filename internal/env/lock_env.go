package env

import (
	"errors"
	"fmt"

	"github.com/patchwork-deploy/internal/state"
)

var ErrEnvLocked = errors.New("environment is locked")

// LockEnvironment marks an environment as locked, preventing patch application.
func LockEnvironment(st *state.State, env string) error {
	recs := st.ForEnvironment(env)
	if len(recs) == 0 {
		return fmt.Errorf("environment %q not found", env)
	}
	key := fmt.Sprintf("__lock__%s", env)
	if st.Has(key, "__sentinel__") {
		return nil // already locked
	}
	st.Add(state.Record{
		Environment: key,
		Patch:       "__sentinel__",
	})
	return nil
}

// UnlockEnvironment removes the lock sentinel for an environment.
func UnlockEnvironment(st *state.State, env string) error {
	key := fmt.Sprintf("__lock__%s", env)
	if !st.Has(key, "__sentinel__") {
		return nil // not locked
	}
	st.Remove(key, "__sentinel__")
	return nil
}

// IsEnvLocked reports whether the given environment is locked.
func IsEnvLocked(st *state.State, env string) bool {
	key := fmt.Sprintf("__lock__%s", env)
	return st.Has(key, "__sentinel__")
}
