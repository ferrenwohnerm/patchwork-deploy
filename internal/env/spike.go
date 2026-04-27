package env

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

const spikePrefix = "spike:"

func spikeKey(env, patch string) string {
	return fmt.Sprintf("%s%s:%s", spikePrefix, env, patch)
}

// SetSpike marks a patch in an environment as a spike with a given concurrency limit.
// A spike limit caps how many concurrent applications of the patch may run.
func SetSpike(st *state.State, env, patch string, limit int) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	if !patchExistsInEnv(st, env, patch) {
		return fmt.Errorf("patch %q not found in environment %q", patch, env)
	}
	if limit <= 0 {
		return fmt.Errorf("spike limit must be greater than zero")
	}
	st.SetMeta(spikeKey(env, patch), strconv.Itoa(limit))
	return nil
}

// GetSpike returns the spike limit for a patch, and whether one is set.
func GetSpike(st *state.State, env, patch string) (int, bool) {
	v, ok := st.GetMeta(spikeKey(env, patch))
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return n, true
}

// ClearSpike removes the spike limit for a patch.
func ClearSpike(st *state.State, env, patch string) {
	st.DeleteMeta(spikeKey(env, patch))
}

// ListSpikes returns all spike limits set for an environment as a map of patch -> limit.
func ListSpikes(st *state.State, env string) map[string]int {
	prefix := fmt.Sprintf("%s%s:", spikePrefix, env)
	result := make(map[string]int)
	for _, k := range st.MetaKeys() {
		if !strings.HasPrefix(k, prefix) {
			continue
		}
		patch := strings.TrimPrefix(k, prefix)
		v, ok := st.GetMeta(k)
		if !ok {
			continue
		}
		n, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		result[patch] = n
	}
	return result
}
