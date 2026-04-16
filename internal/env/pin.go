package env

import (
	"fmt"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

const pinSentinel = "__pinned__"

// PinResult holds the outcome of a pin operation.
type PinResult struct {
	Environment string
	Patch       string
	PinnedAt    time.Time
}

// Pin marks a specific patch as the pinned version for an environment.
// Subsequent runs will not apply patches beyond this point.
func Pin(st *state.State, env, patch string) (PinResult, error) {
	records := st.ForEnvironment(env)
	if len(records) == 0 {
		return PinResult{}, fmt.Errorf("environment %q not found", env)
	}

	found := false
	for _, r := range records {
		if r.Patch == patch {
			found = true
			break
		}
	}
	if !found {
		return PinResult{}, fmt.Errorf("patch %q not found in environment %q", patch, env)
	}

	pinnedAt := time.Now().UTC()
	st.Add(state.Record{
		Environment: env,
		Patch:       pinSentinel,
		AppliedAt:   pinnedAt,
		Tag:         patch,
	})

	return PinResult{Environment: env, Patch: patch, PinnedAt: pinnedAt}, nil
}

// Unpin removes the pin sentinel for an environment.
func Unpin(st *state.State, env string) error {
	records := st.ForEnvironment(env)
	if len(records) == 0 {
		return fmt.Errorf("environment %q not found", env)
	}
	st.Remove(env, pinSentinel)
	return nil
}

// PinnedPatch returns the patch name the environment is pinned to, or empty string.
func PinnedPatch(st *state.State, env string) string {
	for _, r := range st.ForEnvironment(env) {
		if r.Patch == pinSentinel {
			return r.Tag
		}
	}
	return ""
}
