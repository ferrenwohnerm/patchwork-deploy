package env

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

const timeoutKeyPrefix = "timeout:"

func timeoutKey(env, patch string) string {
	return fmt.Sprintf("%s%s:%s", timeoutKeyPrefix, env, patch)
}

// SetTimeout records a maximum allowed duration for a patch deployment in a given environment.
func SetTimeout(st *state.State, env, patch string, d time.Duration) error {
	if d <= 0 {
		return fmt.Errorf("timeout duration must be positive")
	}
	records := st.ForEnvironment(env)
	if len(records) == 0 {
		return fmt.Errorf("environment %q not found", env)
	}
	found := false
	for _, r := range records {
		if r.Patch == patch {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("patch %q not found in environment %q", patch, env)
	}
	key := timeoutKey(env, patch)
	st.SetMeta(key, strconv.FormatInt(int64(d), 10))
	return nil
}

// GetTimeout retrieves the timeout duration for a patch in an environment.
// Returns 0 and no error if no timeout is set.
func GetTimeout(st *state.State, env, patch string) (time.Duration, error) {
	key := timeoutKey(env, patch)
	val := st.GetMeta(key)
	if val == "" {
		return 0, nil
	}
	n, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("corrupt timeout value for %s/%s", env, patch)
	}
	return time.Duration(n), nil
}

// ClearTimeout removes the timeout setting for a patch in an environment.
func ClearTimeout(st *state.State, env, patch string) {
	key := timeoutKey(env, patch)
	st.DeleteMeta(key)
}

// ListTimeouts returns all patch→duration mappings for the given environment.
func ListTimeouts(st *state.State, env string) map[string]time.Duration {
	result := make(map[string]time.Duration)
	prefix := fmt.Sprintf("%s%s:", timeoutKeyPrefix, env)
	for _, kv := range st.AllMeta() {
		if !strings.HasPrefix(kv.Key, prefix) {
			continue
		}
		patch := strings.TrimPrefix(kv.Key, prefix)
		n, err := strconv.ParseInt(kv.Value, 10, 64)
		if err != nil {
			continue
		}
		result[patch] = time.Duration(n)
	}
	return result
}
