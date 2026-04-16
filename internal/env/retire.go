package env

import (
	"fmt"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

const retiredPrefix = "__retired__"

// Retire marks an environment as retired by moving its records to a retired key
// and recording a sentinel entry. Retired environments are excluded from normal
// operations but their history is preserved.
func Retire(st *state.State, env string) error {
	recs := st.ForEnvironment(env)
	if len(recs) == 0 {
		return fmt.Errorf("environment %q not found", env)
	}

	retiredKey := retiredPrefix + env

	// Copy records to retired key
	for _, r := range recs {
		r.Environment = retiredKey
		st.Add(r)
	}

	// Add sentinel
	st.Add(state.Record{
		Environment: retiredKey,
		Patch:       "__retired_sentinel__",
		AppliedAt:   time.Now(),
	})

	// Remove original records
	st.RemoveEnvironment(env)
	return nil
}

// IsRetired reports whether the given environment has been retired.
func IsRetired(st *state.State, env string) bool {
	retiredKey := retiredPrefix + env
	recs := st.ForEnvironment(retiredKey)
	for _, r := range recs {
		if r.Patch == "__retired_sentinel__" {
			return true
		}
	}
	return false
}

// ListRetired returns the names of all retired environments.
func ListRetired(st *state.State) []string {
	seen := map[string]bool{}
	var out []string
	for _, r := range st.All() {
		if len(r.Environment) > len(retiredPrefix) && r.Environment[:len(retiredPrefix)] == retiredPrefix {
			name := r.Environment[len(retiredPrefix):]
			if !seen[name] {
				seen[name] = true
				out = append(out, name)
			}
		}
	}
	return out
}
