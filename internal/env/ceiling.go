package env

import (
	"fmt"
	"strconv"

	"github.com/patchwork-deploy/internal/state"
)

const ceilingKeyPrefix = "ceiling:"

func ceilingKey(env string) string {
	return ceilingKeyPrefix + env
}

// SetCeiling sets the maximum number of patches that may be applied to env.
func SetCeiling(st *state.State, env string, limit int) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	if limit <= 0 {
		return fmt.Errorf("ceiling limit must be greater than zero")
	}
	st.SetMeta(ceilingKey(env), strconv.Itoa(limit))
	return nil
}

// GetCeiling returns the ceiling limit for env, or 0 if none is set.
func GetCeiling(st *state.State, env string) (int, error) {
	if !st.HasEnvironment(env) {
		return 0, fmt.Errorf("environment %q not found", env)
	}
	raw, ok := st.GetMeta(ceilingKey(env))
	if !ok {
		return 0, nil
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid ceiling value for %q: %w", env, err)
	}
	return v, nil
}

// ClearCeiling removes the ceiling limit for env.
func ClearCeiling(st *state.State, env string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	st.DeleteMeta(ceilingKey(env))
	return nil
}

// CheckCeiling returns an error if env has reached or exceeded its ceiling.
func CheckCeiling(st *state.State, env string) error {
	limit, err := GetCeiling(st, env)
	if err != nil {
		return err
	}
	if limit == 0 {
		return nil
	}
	records := st.ForEnvironment(env)
	if len(records) >= limit {
		return fmt.Errorf("environment %q has reached its ceiling of %d patches", env, limit)
	}
	return nil
}
