package env

import (
	"fmt"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

const gracePeriodPrefix = "grace:"

// SetGracePeriod sets a grace period duration for a patch in an environment.
// During the grace period, the patch is considered recently applied and may
// be rolled back without confirmation.
func SetGracePeriod(st *state.State, env, patch string, duration time.Duration) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	if !patchExistsInEnv(st, env, patch) {
		return fmt.Errorf("patch %q not found in environment %q", patch, env)
	}
	if duration <= 0 {
		return fmt.Errorf("grace period must be positive")
	}
	key := gracePeriodPrefix + env + ":" + patch
	expiry := time.Now().Add(duration).Format(time.RFC3339)
	st.SetMeta(key, expiry)
	return nil
}

// GetGracePeriod returns the expiry time of the grace period for a patch, if set.
func GetGracePeriod(st *state.State, env, patch string) (time.Time, bool, error) {
	key := gracePeriodPrefix + env + ":" + patch
	val, ok := st.GetMeta(key)
	if !ok {
		return time.Time{}, false, nil
	}
	t, err := time.Parse(time.RFC3339, val)
	if err != nil {
		return time.Time{}, false, fmt.Errorf("invalid grace period value: %w", err)
	}
	return t, true, nil
}

// InGracePeriod returns true if the patch is currently within its grace period.
func InGracePeriod(st *state.State, env, patch string) (bool, error) {
	expiry, ok, err := GetGracePeriod(st, env, patch)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	return time.Now().Before(expiry), nil
}

// ClearGracePeriod removes the grace period for a patch.
func ClearGracePeriod(st *state.State, env, patch string) {
	key := gracePeriodPrefix + env + ":" + patch
	st.DeleteMeta(key)
}
