package env

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

const weightPrefix = "weight:"

// SetWeight assigns a numeric weight to a patch within an environment.
// Weights are used to influence application ordering when priorities are equal.
func SetWeight(st *state.State, env, patch string, weight int) error {
	if weight < 0 {
		return fmt.Errorf("weight must be non-negative, got %d", weight)
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
	key := weightPrefix + env + ":" + patch
	st.SetMeta(key, strconv.Itoa(weight))
	return nil
}

// GetWeight returns the weight assigned to a patch in an environment.
// Returns 0 and false if no weight is set.
func GetWeight(st *state.State, env, patch string) (int, bool) {
	key := weightPrefix + env + ":" + patch
	val, ok := st.GetMeta(key)
	if !ok {
		return 0, false
	}
	w, err := strconv.Atoi(val)
	if err != nil {
		return 0, false
	}
	return w, true
}

// ClearWeight removes the weight entry for a patch in an environment.
func ClearWeight(st *state.State, env, patch string) {
	key := weightPrefix + env + ":" + patch
	st.DeleteMeta(key)
}

// ListWeights returns all patch->weight mappings for an environment.
func ListWeights(st *state.State, env string) map[string]int {
	result := map[string]int{}
	prefix := weightPrefix + env + ":"
	for k, v := range st.AllMeta() {
		if strings.HasPrefix(k, prefix) {
			patch := strings.TrimPrefix(k, prefix)
			if w, err := strconv.Atoi(v); err == nil {
				result[patch] = w
			}
		}
	}
	return result
}
