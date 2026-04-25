package env

import (
	"fmt"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

const chainKeyPrefix = "chain:"

func chainKey(env, patch string) string {
	return fmt.Sprintf("%s%s:%s", chainKeyPrefix, env, patch)
}

// SetChain registers a successor patch that must be applied after the given patch.
func SetChain(st *state.State, env, patch, successor string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	if !patchExistsInEnv(st, env, patch) {
		return fmt.Errorf("patch %q not found in environment %q", patch, env)
	}
	if !patchExistsInEnv(st, env, successor) {
		return fmt.Errorf("successor patch %q not found in environment %q", successor, env)
	}
	if patch == successor {
		return fmt.Errorf("patch cannot chain to itself")
	}
	key := chainKey(env, patch)
	st.SetMeta(key, successor)
	return nil
}

// GetChain returns the successor patch for the given patch, if any.
func GetChain(st *state.State, env, patch string) (string, bool) {
	key := chainKey(env, patch)
	v, ok := st.GetMeta(key)
	return v, ok
}

// RemoveChain clears the chain entry for the given patch.
func RemoveChain(st *state.State, env, patch string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	key := chainKey(env, patch)
	st.DeleteMeta(key)
	return nil
}

// ListChains returns all chain entries for the given environment as patch -> successor pairs.
func ListChains(st *state.State, env string) map[string]string {
	result := make(map[string]string)
	prefix := fmt.Sprintf("%s%s:", chainKeyPrefix, env)
	for _, k := range st.MetaKeys() {
		if strings.HasPrefix(k, prefix) {
			patch := strings.TrimPrefix(k, prefix)
			if v, ok := st.GetMeta(k); ok {
				result[patch] = v
			}
		}
	}
	return result
}
