package env

import (
	"fmt"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

const fluxPrefix = "flux:"

func fluxKey(env, patch string) string {
	return fmt.Sprintf("%s%s:%s", fluxPrefix, env, patch)
}

// SetFlux attaches a flux mode (e.g. "auto", "manual", "gated") to a patch
// within an environment. It controls how the patch transitions between states.
func SetFlux(st *state.State, env, patch, mode string) error {
	if !envExists(st, env) {
		return fmt.Errorf("environment %q not found", env)
	}
	if !patchExistsInEnv(st, env, patch) {
		return fmt.Errorf("patch %q not found in environment %q", patch, env)
	}
	mode = strings.TrimSpace(mode)
	if mode == "" {
		return fmt.Errorf("flux mode must not be empty")
	}
	if strings.ContainsAny(mode, "\n\r") {
		return fmt.Errorf("flux mode must not contain newlines")
	}
	valid := map[string]bool{"auto": true, "manual": true, "gated": true}
	if !valid[mode] {
		return fmt.Errorf("invalid flux mode %q: must be one of auto, manual, gated", mode)
	}
	st.Set(fluxKey(env, patch), mode)
	return nil
}

// GetFlux returns the flux mode for a patch in an environment.
// Returns an empty string if no mode is set.
func GetFlux(st *state.State, env, patch string) string {
	return st.GetMeta(fluxKey(env, patch))
}

// ClearFlux removes the flux mode for a patch in an environment.
func ClearFlux(st *state.State, env, patch string) error {
	if !envExists(st, env) {
		return fmt.Errorf("environment %q not found", env)
	}
	st.DeleteMeta(fluxKey(env, patch))
	return nil
}

// ListFluxes returns all flux entries for the given environment as a map
// of patch name to mode.
func ListFluxes(st *state.State, env string) map[string]string {
	result := make(map[string]string)
	prefix := fmt.Sprintf("%s%s:", fluxPrefix, env)
	for _, key := range st.MetaKeys() {
		if strings.HasPrefix(key, prefix) {
			patch := strings.TrimPrefix(key, prefix)
			result[patch] = st.GetMeta(key)
		}
	}
	return result
}
