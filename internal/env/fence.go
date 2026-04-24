package env

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

var validFencePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func fenceKey(env, patch string) string {
	return fmt.Sprintf("fence:%s:%s", env, patch)
}

// SetFence attaches a named gate to a patch within an environment.
// Patches with a fence will only be considered ready when the fence is cleared.
func SetFence(st *state.State, env, patch, name string) error {
	if !validFencePattern.MatchString(name) {
		return fmt.Errorf("fence name %q must match [a-zA-Z0-9_-]", name)
	}
	if !envExistsInState(st, env) {
		return fmt.Errorf("environment %q not found", env)
	}
	if !patchExistsInEnv(st, env, patch) {
		return fmt.Errorf("patch %q not found in environment %q", patch, env)
	}
	st.SetMeta(fenceKey(env, patch), name)
	return nil
}

// GetFence returns the fence name for a patch, or empty string if none is set.
func GetFence(st *state.State, env, patch string) string {
	return st.GetMeta(fenceKey(env, patch))
}

// ClearFence removes the fence from a patch.
func ClearFence(st *state.State, env, patch string) error {
	if !envExistsInState(st, env) {
		return fmt.Errorf("environment %q not found", env)
	}
	st.DeleteMeta(fenceKey(env, patch))
	return nil
}

// IsFenced returns true if the given patch has an active fence.
func IsFenced(st *state.State, env, patch string) bool {
	return GetFence(st, env, patch) != ""
}

// ListFences returns all patch->fence mappings for an environment.
func ListFences(st *state.State, env string) map[string]string {
	prefix := fmt.Sprintf("fence:%s:", env)
	result := make(map[string]string)
	for _, key := range st.MetaKeys() {
		if strings.HasPrefix(key, prefix) {
			patch := strings.TrimPrefix(key, prefix)
			result[patch] = st.GetMeta(key)
		}
	}
	return result
}
