package env

import (
	"fmt"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

const intentPrefix = "intent:"

// SetIntent records a deployment intent (reason/goal) for a patch in an environment.
func SetIntent(st *state.State, env, patch, intent string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	if !patchExistsInEnv(st, env, patch) {
		return fmt.Errorf("patch %q not found in environment %q", patch, env)
	}
	if strings.ContainsAny(intent, "\n\r") {
		return fmt.Errorf("intent must not contain newlines")
	}
	if strings.TrimSpace(intent) == "" {
		return fmt.Errorf("intent must not be empty")
	}
	st.SetMeta(env, intentPrefix+patch, intent)
	return nil
}

// GetIntent returns the recorded intent for a patch in an environment.
func GetIntent(st *state.State, env, patch string) (string, bool) {
	v, ok := st.GetMeta(env, intentPrefix+patch)
	return v, ok
}

// RemoveIntent clears the intent for a patch in an environment.
func RemoveIntent(st *state.State, env, patch string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	st.DeleteMeta(env, intentPrefix+patch)
	return nil
}

// ListIntents returns all patch->intent mappings for an environment.
func ListIntents(st *state.State, env string) (map[string]string, error) {
	if !st.HasEnvironment(env) {
		return nil, fmt.Errorf("environment %q not found", env)
	}
	all := st.ListMeta(env)
	result := make(map[string]string)
	for k, v := range all {
		if strings.HasPrefix(k, intentPrefix) {
			patch := strings.TrimPrefix(k, intentPrefix)
			result[patch] = v
		}
	}
	return result, nil
}
