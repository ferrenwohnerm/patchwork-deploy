package env

import (
	"fmt"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

const cooldownKeyPrefix = "cooldown:"

func cooldownKey(env, patch string) string {
	return fmt.Sprintf("%s%s:%s", cooldownKeyPrefix, env, patch)
}

// SetCooldown records a cooldown duration for a patch in an environment.
// After a patch is applied, it cannot be re-applied until the cooldown expires.
func SetCooldown(st *state.State, env, patch string, duration time.Duration) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	if !patchExistsInEnv(st, env, patch) {
		return fmt.Errorf("patch %q not found in environment %q", patch, env)
	}
	if duration <= 0 {
		return fmt.Errorf("cooldown duration must be positive")
	}
	expiry := time.Now().Add(duration).Format(time.RFC3339)
	st.SetMeta(cooldownKey(env, patch), expiry)
	return nil
}

// GetCooldown returns the expiry time for a patch cooldown, and whether one is set.
func GetCooldown(st *state.State, env, patch string) (time.Time, bool, error) {
	val, ok := st.GetMeta(cooldownKey(env, patch))
	if !ok {
		return time.Time{}, false, nil
	}
	t, err := time.Parse(time.RFC3339, val)
	if err != nil {
		return time.Time{}, false, fmt.Errorf("invalid cooldown value: %w", err)
	}
	return t, true, nil
}

// InCooldown reports whether the cooldown period for a patch is still active.
func InCooldown(st *state.State, env, patch string) (bool, error) {
	expiry, ok, err := GetCooldown(st, env, patch)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	return time.Now().Before(expiry), nil
}

// ClearCooldown removes the cooldown for a patch in an environment.
func ClearCooldown(st *state.State, env, patch string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	st.DeleteMeta(cooldownKey(env, patch))
	return nil
}
