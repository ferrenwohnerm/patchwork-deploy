package env

import (
	"fmt"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

const suppressPrefix = "suppress:"

// Suppress marks a patch in an environment as suppressed, preventing it from
// being applied during future runs.
func Suppress(st *state.State, env, patch, reason string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	if !patchExistsInEnv(st, env, patch) {
		return fmt.Errorf("patch %q not found in environment %q", patch, env)
	}
	if strings.ContainsAny(reason, "\n\r") {
		return fmt.Errorf("reason must not contain newlines")
	}
	key := suppressKey(env, patch)
	st.SetMeta(key, reason)
	return nil
}

// Unsuppress removes the suppression marker from a patch in an environment.
func Unsuppress(st *state.State, env, patch string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	st.DeleteMeta(suppressKey(env, patch))
	return nil
}

// IsSuppressed reports whether the given patch is suppressed in the environment.
func IsSuppressed(st *state.State, env, patch string) bool {
	_, ok := st.GetMeta(suppressKey(env, patch))
	return ok
}

// ListSuppressed returns all suppressed patches for an environment along with
// their suppression reasons.
func ListSuppressed(st *state.State, env string) map[string]string {
	prefix := suppressPrefix + env + ":"
	result := make(map[string]string)
	for k, v := range st.AllMeta() {
		if strings.HasPrefix(k, prefix) {
			patch := strings.TrimPrefix(k, prefix)
			result[patch] = v
		}
	}
	return result
}

func suppressKey(env, patch string) string {
	return suppressPrefix + env + ":" + patch
}
