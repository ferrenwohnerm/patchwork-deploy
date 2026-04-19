package env

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

const throttlePrefix = "__throttle__"

// SetThrottle sets the maximum number of patches that can be applied per run
// for the given environment.
func SetThrottle(st *state.State, env string, limit int) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	if limit <= 0 {
		return fmt.Errorf("throttle limit must be greater than zero")
	}
	key := throttlePrefix + env
	st.SetMeta(key, strconv.Itoa(limit))
	return nil
}

// GetThrottle returns the throttle limit for the environment, or 0 if not set.
func GetThrottle(st *state.State, env string) (int, error) {
	if !st.HasEnvironment(env) {
		return 0, fmt.Errorf("environment %q not found", env)
	}
	key := throttlePrefix + env
	val := st.GetMeta(key)
	if val == "" {
		return 0, nil
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("invalid throttle value: %w", err)
	}
	return n, nil
}

// ClearThrottle removes the throttle limit for the environment.
func ClearThrottle(st *state.State, env string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	key := throttlePrefix + env
	st.DeleteMeta(key)
	return nil
}

// ListThrottles returns a map of env -> limit for all environments with a throttle set.
func ListThrottles(st *state.State) map[string]int {
	result := make(map[string]int)
	for k, v := range st.AllMeta() {
		if !strings.HasPrefix(k, throttlePrefix) {
			continue
		}
		env := strings.TrimPrefix(k, throttlePrefix)
		if n, err := strconv.Atoi(v); err == nil {
			result[env] = n
		}
	}
	return result
}
