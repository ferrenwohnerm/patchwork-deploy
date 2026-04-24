package env

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

var validAffinityName = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func affinityKey(env, patch string) string {
	return fmt.Sprintf("affinity:%s:%s", env, patch)
}

// SetAffinity assigns a named affinity group to a patch within an environment.
func SetAffinity(st *state.State, env, patch, group string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	if !patchExistsInEnv(st, env, patch) {
		return fmt.Errorf("patch %q not found in environment %q", patch, env)
	}
	if !validAffinityName.MatchString(group) {
		return fmt.Errorf("affinity group %q contains invalid characters", group)
	}
	st.SetMeta(affinityKey(env, patch), group)
	return nil
}

// GetAffinity returns the affinity group assigned to a patch, or empty string if none.
func GetAffinity(st *state.State, env, patch string) string {
	return st.GetMeta(affinityKey(env, patch))
}

// RemoveAffinity clears the affinity group for a patch.
func RemoveAffinity(st *state.State, env, patch string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	st.DeleteMeta(affinityKey(env, patch))
	return nil
}

// ListAffinities returns a map of patch -> affinity group for all patches in the environment.
func ListAffinities(st *state.State, env string) map[string]string {
	prefix := fmt.Sprintf("affinity:%s:", env)
	result := make(map[string]string)
	for k, v := range st.AllMeta() {
		if strings.HasPrefix(k, prefix) {
			patch := strings.TrimPrefix(k, prefix)
			result[patch] = v
		}
	}
	return result
}
