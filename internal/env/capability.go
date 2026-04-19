package env

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

var validCapabilityKey = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func capabilityKey(env, cap string) string {
	return fmt.Sprintf("__capability__%s__%s", env, cap)
}

// AddCapability registers a named capability for an environment.
func AddCapability(st *state.State, env, cap string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	if !validCapabilityKey.MatchString(cap) {
		return fmt.Errorf("capability name %q contains invalid characters", cap)
	}
	key := capabilityKey(env, cap)
	st.SetMeta(key, "true")
	return nil
}

// RemoveCapability removes a named capability from an environment.
func RemoveCapability(st *state.State, env, cap string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	st.DeleteMeta(capabilityKey(env, cap))
	return nil
}

// HasCapability returns true if the environment has the given capability.
func HasCapability(st *state.State, env, cap string) bool {
	v, ok := st.GetMeta(capabilityKey(env, cap))
	return ok && v == "true"
}

// ListCapabilities returns all capability names registered for an environment.
func ListCapabilities(st *state.State, env string) []string {
	prefix := fmt.Sprintf("__capability__%s__", env)
	var caps []string
	for _, k := range st.MetaKeys() {
		if strings.HasPrefix(k, prefix) {
			caps = append(caps, strings.TrimPrefix(k, prefix))
		}
	}
	return caps
}
