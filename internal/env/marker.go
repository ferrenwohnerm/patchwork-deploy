package env

import (
	"fmt"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

const markerPrefix = "marker:"

func markerKey(env, name string) string {
	return fmt.Sprintf("%s%s:%s", markerPrefix, env, name)
}

// SetMarker attaches a named marker to an environment.
func SetMarker(st *state.State, env, name, value string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	if strings.ContainsAny(name, " \t\n") {
		return fmt.Errorf("marker name must not contain whitespace")
	}
	if strings.ContainsAny(value, "\n") {
		return fmt.Errorf("marker value must not contain newlines")
	}
	st.SetMeta(markerKey(env, name), value)
	return nil
}

// GetMarker returns the value of a named marker for an environment.
func GetMarker(st *state.State, env, name string) (string, bool) {
	v, ok := st.GetMeta(markerKey(env, name))
	return v, ok
}

// RemoveMarker deletes a named marker from an environment.
func RemoveMarker(st *state.State, env, name string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	st.DeleteMeta(markerKey(env, name))
	return nil
}

// ListMarkers returns all markers set for an environment as a map.
func ListMarkers(st *state.State, env string) map[string]string {
	prefix := fmt.Sprintf("%s%s:", markerPrefix, env)
	result := map[string]string{}
	for _, k := range st.MetaKeys() {
		if strings.HasPrefix(k, prefix) {
			name := strings.TrimPrefix(k, prefix)
			v, _ := st.GetMeta(k)
			result[name] = v
		}
	}
	return result
}
