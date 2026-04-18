package env

import (
	"fmt"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

type DepGraph map[string][]string

// SetDep records that patch depends on another patch within an environment.
func SetDep(st *state.State, env, patch, dependsOn string) error {
	records := st.ForEnvironment(env)
	if len(records) == 0 {
		return fmt.Errorf("environment %q not found", env)
	}
	has := func(p string) bool {
		for _, r := range records {
			if r.Patch == p {
				return true
			}
		}
		return false
	}
	if !has(patch) {
		return fmt.Errorf("patch %q not found in environment %q", patch, env)
	}
	if !has(dependsOn) {
		return fmt.Errorf("dependency patch %q not found in environment %q", dependsOn, env)
	}
	key := depKey(env, patch)
	existing := st.GetMeta(key)
	parts := splitLines(existing)
	for _, p := range parts {
		if p == dependsOn {
			return nil // already recorded
		}
	}
	parts = append(parts, dependsOn)
	st.SetMeta(key, strings.Join(parts, "\n"))
	return nil
}

// GetDeps returns the direct dependencies of a patch in an environment.
func GetDeps(st *state.State, env, patch string) ([]string, error) {
	records := st.ForEnvironment(env)
	if len(records) == 0 {
		return nil, fmt.Errorf("environment %q not found", env)
	}
	raw := st.GetMeta(depKey(env, patch))
	if raw == "" {
		return []string{}, nil
	}
	return splitLines(raw), nil
}

// RemoveDep removes a single dependency edge.
func RemoveDep(st *state.State, env, patch, dependsOn string) error {
	key := depKey(env, patch)
	existing := splitLines(st.GetMeta(key))
	filtered := existing[:0]
	for _, p := range existing {
		if p != dependsOn {
			filtered = append(filtered, p)
		}
	}
	st.SetMeta(key, strings.Join(filtered, "\n"))
	return nil
}

func depKey(env, patch string) string {
	return fmt.Sprintf("__dep__%s__%s", env, patch)
}
