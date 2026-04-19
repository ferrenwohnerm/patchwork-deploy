package env

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

var validAlias = regexp.MustCompile(`^[a-z0-9_-]+$`)

const aliasPrefix = "alias:"

// SetAlias assigns a short alias to an environment name.
func SetAlias(st *state.State, env, alias string) error {
	if !validAlias.MatchString(alias) {
		return fmt.Errorf("alias %q contains invalid characters (a-z, 0-9, _, - only)", alias)
	}
	records := st.ForEnvironment(env)
	if len(records) == 0 {
		return fmt.Errorf("environment %q not found", env)
	}
	key := aliasPrefix + env
	st.SetMeta(key, alias)
	return nil
}

// GetAlias returns the alias set for an environment, or empty string if none.
func GetAlias(st *state.State, env string) string {
	return st.GetMeta(aliasPrefix + env)
}

// RemoveAlias clears the alias for an environment.
func RemoveAlias(st *state.State, env string) error {
	records := st.ForEnvironment(env)
	if len(records) == 0 {
		return fmt.Errorf("environment %q not found", env)
	}
	st.DeleteMeta(aliasPrefix + env)
	return nil
}

// ResolveAlias returns the canonical env name for a given alias, or empty string.
func ResolveAlias(st *state.State, alias string) string {
	for _, key := range st.MetaKeys() {
		if strings.HasPrefix(key, aliasPrefix) {
			if st.GetMeta(key) == alias {
				return strings.TrimPrefix(key, aliasPrefix)
			}
		}
	}
	return ""
}
