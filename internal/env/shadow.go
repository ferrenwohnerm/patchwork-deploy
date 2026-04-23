package env

import (
	"fmt"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

const shadowPrefix = "shadow:"

// SetShadow marks targetEnv as a shadow of sourceEnv for a specific patch.
// A shadow environment mirrors a patch without executing it independently.
func SetShadow(st *state.State, sourceEnv, targetEnv, patch string) error {
	if sourceEnv == targetEnv {
		return fmt.Errorf("source and target environment must differ")
	}
	if strings.ContainsAny(sourceEnv, " \t\n") || strings.ContainsAny(targetEnv, " \t\n") {
		return fmt.Errorf("environment name must not contain whitespace")
	}
	if !st.HasEnvironment(sourceEnv) {
		return fmt.Errorf("source environment %q not found", sourceEnv)
	}
	if !st.HasEnvironment(targetEnv) {
		return fmt.Errorf("target environment %q not found", targetEnv)
	}
	if !patchExistsInEnv(st, sourceEnv, patch) {
		return fmt.Errorf("patch %q not found in environment %q", patch, sourceEnv)
	}
	key := shadowKey(targetEnv, patch)
	st.SetMeta(key, sourceEnv)
	return nil
}

// GetShadow returns the source environment that targetEnv shadows for the given patch.
// Returns an empty string if no shadow relationship is set.
func GetShadow(st *state.State, targetEnv, patch string) string {
	key := shadowKey(targetEnv, patch)
	return st.GetMeta(key)
}

// RemoveShadow clears the shadow relationship for a patch in targetEnv.
func RemoveShadow(st *state.State, targetEnv, patch string) error {
	if !st.HasEnvironment(targetEnv) {
		return fmt.Errorf("environment %q not found", targetEnv)
	}
	key := shadowKey(targetEnv, patch)
	st.DeleteMeta(key)
	return nil
}

// ListShadows returns all patch→source mappings where targetEnv is a shadow.
func ListShadows(st *state.State, targetEnv string) map[string]string {
	result := make(map[string]string)
	prefix := shadowPrefix + targetEnv + ":"
	for k, v := range st.AllMeta() {
		if strings.HasPrefix(k, prefix) {
			patch := strings.TrimPrefix(k, prefix)
			result[patch] = v
		}
	}
	return result
}

func shadowKey(env, patch string) string {
	return shadowPrefix + env + ":" + patch
}
