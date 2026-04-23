package env

import (
	"fmt"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

const relayPrefix = "relay:"

// SetRelay configures env to forward patch application events to targetEnv.
func SetRelay(st *state.State, env, patch, targetEnv string) error {
	if env == targetEnv {
		return fmt.Errorf("relay: source and target environment must differ")
	}
	records := st.ForEnvironment(env)
	if len(records) == 0 {
		return fmt.Errorf("relay: environment %q not found", env)
	}
	if !patchExistsInEnv(records, patch) {
		return fmt.Errorf("relay: patch %q not found in environment %q", patch, env)
	}
	if strings.ContainsAny(targetEnv, "\n\r") {
		return fmt.Errorf("relay: target environment name must not contain newlines")
	}
	key := relayKey(env, patch)
	st.SetMeta(key, targetEnv)
	return nil
}

// GetRelay returns the relay target for a patch in an environment, if any.
func GetRelay(st *state.State, env, patch string) (string, bool) {
	key := relayKey(env, patch)
	val, ok := st.GetMeta(key)
	return val, ok
}

// RemoveRelay clears the relay configuration for a patch in an environment.
func RemoveRelay(st *state.State, env, patch string) error {
	records := st.ForEnvironment(env)
	if len(records) == 0 {
		return fmt.Errorf("relay: environment %q not found", env)
	}
	key := relayKey(env, patch)
	st.DeleteMeta(key)
	return nil
}

// ListRelays returns all relay mappings for an environment as patch->target pairs.
func ListRelays(st *state.State, env string) map[string]string {
	result := make(map[string]string)
	prefix := relayPrefix + env + ":"
	for k, v := range st.AllMeta() {
		if strings.HasPrefix(k, prefix) {
			patch := strings.TrimPrefix(k, prefix)
			result[patch] = v
		}
	}
	return result
}

func relayKey(env, patch string) string {
	return relayPrefix + env + ":" + patch
}
