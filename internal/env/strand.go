package env

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

var validStrandName = regexp.MustCompile(`^[a-z0-9_-]+$`)

func strandKey(env, patch string) string {
	return fmt.Sprintf("strand::%s::%s", env, patch)
}

// SetStrand assigns a named strand (execution branch) to a patch within an environment.
func SetStrand(st *state.State, env, patch, strand string) error {
	if !validStrandName.MatchString(strand) {
		return fmt.Errorf("invalid strand name %q: use only a-z, 0-9, hyphens, underscores", strand)
	}
	if !hasEnv(st, env) {
		return fmt.Errorf("environment %q not found", env)
	}
	if !patchExistsInEnv(st, env, patch) {
		return fmt.Errorf("patch %q not found in environment %q", patch, env)
	}
	st.Set(strandKey(env, patch), strand)
	return nil
}

// GetStrand returns the strand assigned to a patch, or empty string if none.
func GetStrand(st *state.State, env, patch string) string {
	return st.GetMeta(strandKey(env, patch))
}

// ClearStrand removes the strand assignment for a patch.
func ClearStrand(st *state.State, env, patch string) {
	st.Delete(strandKey(env, patch))
}

// ListStrands returns a map of patch -> strand for all strands in an environment.
func ListStrands(st *state.State, env string) map[string]string {
	prefix := fmt.Sprintf("strand::%s::", env)
	result := map[string]string{}
	for _, k := range st.Keys() {
		if strings.HasPrefix(k, prefix) {
			patch := strings.TrimPrefix(k, prefix)
			result[patch] = st.GetMeta(k)
		}
	}
	return result
}
