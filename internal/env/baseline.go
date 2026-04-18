package env

import (
	"fmt"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

const baselineKeyPrefix = "__baseline__"

type BaselineEntry struct {
	Environment string
	Patch       string
	CreatedAt   time.Time
}

// SetBaseline marks a patch as the baseline for an environment.
// All patches at or before this point are considered pre-applied.
func SetBaseline(st *state.State, env, patch string) error {
	recs := st.ForEnvironment(env)
	if len(recs) == 0 {
		return fmt.Errorf("environment %q not found", env)
	}
	found := false
	for _, r := range recs {
		if r.Patch == patch {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("patch %q not found in environment %q", patch, env)
	}
	key := baselineKeyPrefix + env
	st.Remove(key)
	st.Add(state.Record{
		Environment: key,
		Patch:       patch,
		AppliedAt:   time.Now().UTC(),
	})
	return nil
}

// GetBaseline returns the baseline patch for an environment, if set.
func GetBaseline(st *state.State, env string) (BaselineEntry, bool) {
	key := baselineKeyPrefix + env
	recs := st.ForEnvironment(key)
	if len(recs) == 0 {
		return BaselineEntry{}, false
	}
	r := recs[0]
	return BaselineEntry{Environment: env, Patch: r.Patch, CreatedAt: r.AppliedAt}, true
}

// ClearBaseline removes the baseline marker for an environment.
func ClearBaseline(st *state.State, env string) error {
	key := baselineKeyPrefix + env
	recs := st.ForEnvironment(key)
	if len(recs) == 0 {
		return fmt.Errorf("no baseline set for environment %q", env)
	}
	st.Remove(key)
	return nil
}
