package env

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

var validScopeName = regexp.MustCompile(`^[a-z0-9_-]+$`)

func scopeKey(env, patch string) string {
	return fmt.Sprintf("scope:%s:%s", env, patch)
}

// SetScope assigns a named scope to a patch within an environment.
// Scopes allow grouping patches into logical categories (e.g. "infra", "app").
func SetScope(st *state.State, env, patch, scope string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	if !patchExistsInEnv(st, env, patch) {
		return fmt.Errorf("patch %q not found in environment %q", patch, env)
	}
	if !validScopeName.MatchString(scope) {
		return fmt.Errorf("invalid scope name %q: must match [a-z0-9_-]+", scope)
	}
	st.SetMeta(scopeKey(env, patch), scope)
	return nil
}

// GetScope returns the scope assigned to a patch, or empty string if none.
func GetScope(st *state.State, env, patch string) string {
	return st.GetMeta(scopeKey(env, patch))
}

// ClearScope removes the scope assignment from a patch.
func ClearScope(st *state.State, env, patch string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	st.DeleteMeta(scopeKey(env, patch))
	return nil
}

// ListScopes returns all patch→scope mappings for the given environment.
func ListScopes(st *state.State, env string) map[string]string {
	prefix := fmt.Sprintf("scope:%s:", env)
	result := make(map[string]string)
	for k, v := range st.AllMeta() {
		if strings.HasPrefix(k, prefix) {
			patch := strings.TrimPrefix(k, prefix)
			result[patch] = v
		}
	}
	return result
}
