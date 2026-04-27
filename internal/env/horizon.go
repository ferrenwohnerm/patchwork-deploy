package env

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

const horizonPrefix = "horizon:"

func horizonKey(env, patch string) string {
	return fmt.Sprintf("%s%s:%s", horizonPrefix, env, patch)
}

// SetHorizon sets a maximum number of future patch applications allowed for a
// given patch in an environment. It acts as a forward-looking cap.
func SetHorizon(st *state.State, env, patch string, limit int) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	if !patchExistsInEnv(st, env, patch) {
		return fmt.Errorf("patch %q not found in environment %q", patch, env)
	}
	if limit <= 0 {
		return fmt.Errorf("horizon limit must be greater than zero")
	}
	st.SetMeta(horizonKey(env, patch), strconv.Itoa(limit))
	return nil
}

// GetHorizon returns the horizon limit for a patch in an environment.
// Returns 0 and false if no horizon is set.
func GetHorizon(st *state.State, env, patch string) (int, bool) {
	v, ok := st.GetMeta(horizonKey(env, patch))
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return n, true
}

// ClearHorizon removes the horizon limit for a patch in an environment.
func ClearHorizon(st *state.State, env, patch string) {
	st.DeleteMeta(horizonKey(env, patch))
}

// ListHorizons returns all horizon entries for the given environment as a
// map of patch name to limit.
func ListHorizons(st *state.State, env string) map[string]int {
	result := make(map[string]int)
	prefix := fmt.Sprintf("%s%s:", horizonPrefix, env)
	for _, kv := range st.AllMeta() {
		if !strings.HasPrefix(kv.Key, prefix) {
			continue
		}
		patch := strings.TrimPrefix(kv.Key, prefix)
		n, err := strconv.Atoi(kv.Value)
		if err != nil {
			continue
		}
		result[patch] = n
	}
	return result
}
