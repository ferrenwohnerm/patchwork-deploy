package env

import (
	"fmt"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

// WatchResult holds a snapshot comparison for an environment.
type WatchResult struct {
	Environment string
	PatchCount  int
	LastApplied string
	LastChecked time.Time
	Drifted     bool
	DriftDetail string
}

// Watch checks the current state of an environment against its last snapshot
// and reports whether drift has occurred.
func Watch(st *state.State, env string) (WatchResult, error) {
	records := st.ForEnvironment(env)
	if records == nil {
		return WatchResult{}, fmt.Errorf("environment %q not found", env)
	}

	result := WatchResult{
		Environment: env,
		PatchCount:  len(records),
		LastChecked: time.Now().UTC(),
	}

	if len(records) > 0 {
		latest := records[len(records)-1]
		result.LastApplied = latest.Patch
	}

	snap, err := state.LoadSnapshot(env)
	if err != nil {
		// No snapshot yet — not drifted, just untracked.
		return result, nil
	}

	snapIndex := make(map[string]bool, len(snap))
	for _, r := range snap {
		snapIndex[r.Patch] = true
	}

	for _, r := range records {
		if !snapIndex[r.Patch] {
			result.Drifted = true
			result.DriftDetail = fmt.Sprintf("patch %q applied since last snapshot", r.Patch)
			break
		}
	}

	return result, nil
}
