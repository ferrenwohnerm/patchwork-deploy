package env

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

var validVectorName = regexp.MustCompile(`^[a-z0-9_-]+$`)

func vectorKey(env, patch string) string {
	return fmt.Sprintf("vector:%s:%s", env, patch)
}

// SetVector assigns a named routing vector to a patch within an environment.
// The vector name must be lowercase alphanumeric with hyphens or underscores.
func SetVector(st *state.State, env, patch, vector string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	if !patchExistsInEnv(st, env, patch) {
		return fmt.Errorf("patch %q not found in environment %q", patch, env)
	}
	if !validVectorName.MatchString(vector) {
		return fmt.Errorf("invalid vector name %q: must match [a-z0-9_-]+", vector)
	}
	st.SetMeta(vectorKey(env, patch), vector)
	return nil
}

// GetVector returns the routing vector assigned to a patch, or empty string if none.
func GetVector(st *state.State, env, patch string) string {
	return st.GetMeta(vectorKey(env, patch))
}

// ClearVector removes the routing vector from a patch.
func ClearVector(st *state.State, env, patch string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	st.DeleteMeta(vectorKey(env, patch))
	return nil
}

// ListVectors returns all patch→vector mappings for the given environment.
func ListVectors(st *state.State, env string) map[string]string {
	prefix := fmt.Sprintf("vector:%s:", env)
	result := make(map[string]string)
	for _, k := range st.MetaKeys() {
		if strings.HasPrefix(k, prefix) {
			patch := strings.TrimPrefix(k, prefix)
			result[patch] = st.GetMeta(k)
		}
	}
	return result
}
