package env

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

func waveKey(env, patch string) string {
	return fmt.Sprintf("wave:%s:%s", env, patch)
}

// SetWave assigns a numeric deployment wave to a patch within an environment.
// Waves control the order in which patches are applied across a fleet.
func SetWave(st *state.State, env, patch string, wave int) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	if !patchExistsInEnv(st, env, patch) {
		return fmt.Errorf("patch %q not found in environment %q", patch, env)
	}
	if wave < 0 {
		return fmt.Errorf("wave must be a non-negative integer")
	}
	st.SetMeta(waveKey(env, patch), strconv.Itoa(wave))
	return nil
}

// GetWave returns the wave number assigned to a patch, or -1 if none is set.
func GetWave(st *state.State, env, patch string) int {
	v := st.GetMeta(waveKey(env, patch))
	if v == "" {
		return -1
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return -1
	}
	return n
}

// ClearWave removes the wave assignment for a patch.
func ClearWave(st *state.State, env, patch string) {
	st.DeleteMeta(waveKey(env, patch))
}

// WaveEntry holds a patch name and its wave number.
type WaveEntry struct {
	Patch string
	Wave  int
}

// ListWaves returns all wave assignments for the given environment, sorted by wave then patch.
func ListWaves(st *state.State, env string) []WaveEntry {
	prefix := fmt.Sprintf("wave:%s:", env)
	var entries []WaveEntry
	for _, k := range st.MetaKeys() {
		if !strings.HasPrefix(k, prefix) {
			continue
		}
		patch := strings.TrimPrefix(k, prefix)
		n, err := strconv.Atoi(st.GetMeta(k))
		if err != nil {
			continue
		}
		entries = append(entries, WaveEntry{Patch: patch, Wave: n})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Wave != entries[j].Wave {
			return entries[i].Wave < entries[j].Wave
		}
		return entries[i].Patch < entries[j].Patch
	})
	return entries
}
