package env

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

const priorityKeyPrefix = "priority:"

// SetPriority assigns a numeric priority to a patch within an environment.
// Higher values indicate higher priority.
func SetPriority(st *state.State, env, patch string, priority int) error {
	if priority < 0 {
		return fmt.Errorf("priority must be non-negative, got %d", priority)
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
	key := priorityKeyPrefix + env + ":" + patch
	st.SetMeta(key, strconv.Itoa(priority))
	return nil
}

// GetPriority returns the priority for a patch in an environment.
// Returns 0 and false if no priority is set.
func GetPriority(st *state.State, env, patch string) (int, bool) {
	key := priorityKeyPrefix + env + ":" + patch
	val, ok := st.GetMeta(key)
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return 0, false
	}
	return n, true
}

// ClearPriority removes the priority setting for a patch.
func ClearPriority(st *state.State, env, patch string) {
	key := priorityKeyPrefix + env + ":" + patch
	st.RemoveMeta(key)
}

// ListPriorities returns a map of patch -> priority for all patches in env
// that have an explicit priority set.
func ListPriorities(st *state.State, env string) map[string]int {
	out := map[string]int{}
	prefix := priorityKeyPrefix + env + ":"
	for k, v := range st.AllMeta() {
		if !strings.HasPrefix(k, prefix) {
			continue
		}
		patch := strings.TrimPrefix(k, prefix)
		n, err := strconv.Atoi(v)
		if err == nil {
			out[patch] = n
		}
	}
	return out
}
