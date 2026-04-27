package env

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

const surgeKeyPrefix = "surge:"

func surgeKey(env, patch string) string {
	return fmt.Sprintf("%s%s:%s", surgeKeyPrefix, env, patch)
}

// SetSurge records a maximum concurrent-patch surge limit for a patch in an environment.
func SetSurge(st *state.State, env, patch string, limit int) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	if !patchExistsInEnv(st, env, patch) {
		return fmt.Errorf("patch %q not found in environment %q", patch, env)
	}
	if limit <= 0 {
		return fmt.Errorf("surge limit must be greater than zero")
	}
	st.SetMeta(surgeKey(env, patch), strconv.Itoa(limit))
	return nil
}

// GetSurge returns the surge limit for a patch, and whether one is set.
func GetSurge(st *state.State, env, patch string) (int, bool) {
	v, ok := st.GetMeta(surgeKey(env, patch))
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return n, true
}

// ClearSurge removes the surge limit for a patch.
func ClearSurge(st *state.State, env, patch string) {
	st.DeleteMeta(surgeKey(env, patch))
}

// ListSurges returns all surge entries for an environment as a map of patch -> limit.
func ListSurges(st *state.State, env string) map[string]int {
	prefix := fmt.Sprintf("%s%s:", surgeKeyPrefix, env)
	result := make(map[string]int)
	for _, k := range st.MetaKeys() {
		if !strings.HasPrefix(k, prefix) {
			continue
		}
		patch := strings.TrimPrefix(k, prefix)
		if v, ok := st.GetMeta(k); ok {
			if n, err := strconv.Atoi(v); err == nil {
				result[patch] = n
			}
		}
	}
	return result
}
