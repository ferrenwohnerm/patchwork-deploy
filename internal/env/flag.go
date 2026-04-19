package env

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

var validFlagKey = regexp.MustCompile(`^[a-z][a-z0-9_-]*$`)

func flagKey(env, flag string) string {
	return fmt.Sprintf("__flag__%s__%s", env, flag)
}

// SetFlag attaches a named boolean flag to an environment.
func SetFlag(st *state.State, env, flag string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	if !validFlagKey.MatchString(flag) {
		return fmt.Errorf("flag name %q is invalid: use lowercase letters, digits, hyphens or underscores", flag)
	}
	st.SetMeta(flagKey(env, flag), "true")
	return nil
}

// UnsetFlag removes a named flag from an environment.
func UnsetFlag(st *state.State, env, flag string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	st.DeleteMeta(flagKey(env, flag))
	return nil
}

// HasFlag reports whether the named flag is set on the environment.
func HasFlag(st *state.State, env, flag string) bool {
	v, ok := st.GetMeta(flagKey(env, flag))
	return ok && v == "true"
}

// ListFlags returns all flag names set on the environment.
func ListFlags(st *state.State, env string) []string {
	prefix := fmt.Sprintf("__flag__%s__", env)
	var flags []string
	for _, k := range st.MetaKeys() {
		if strings.HasPrefix(k, prefix) {
			flags = append(flags, strings.TrimPrefix(k, prefix))
		}
	}
	return flags
}
