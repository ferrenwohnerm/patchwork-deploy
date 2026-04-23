package env

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/your-org/patchwork-deploy/internal/state"
)

const retentionPrefix = "meta:retention:"

// SetRetention sets the maximum number of applied patch records to keep for
// an environment. Records beyond the limit are pruned oldest-first.
func SetRetention(st *state.State, env string, limit int) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	if limit < 1 {
		return fmt.Errorf("retention limit must be at least 1")
	}
	key := retentionPrefix + env
	st.SetMeta(key, strconv.Itoa(limit))
	return nil
}

// GetRetention returns the configured retention limit for an environment.
// Returns 0 and false if no limit is set.
func GetRetention(st *state.State, env string) (int, bool) {
	key := retentionPrefix + env
	val, ok := st.GetMeta(key)
	if !ok || strings.TrimSpace(val) == "" {
		return 0, false
	}
	n, err := strconv.Atoi(val)
	if err != nil || n < 1 {
		return 0, false
	}
	return n, true
}

// ClearRetention removes the retention limit for an environment.
func ClearRetention(st *state.State, env string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	st.DeleteMeta(retentionPrefix + env)
	return nil
}
