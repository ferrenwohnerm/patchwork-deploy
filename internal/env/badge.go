package env

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

var validBadgeKey = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

const badgePrefix = "badge:"

func badgeKey(env, key string) string {
	return fmt.Sprintf("%s%s:%s", badgePrefix, env, key)
}

// SetBadge attaches a display badge (key=value) to an environment.
func SetBadge(st *state.State, env, key, value string) error {
	if !hasEnv(st, env) {
		return fmt.Errorf("environment %q not found", env)
	}
	if !validBadgeKey.MatchString(key) {
		return fmt.Errorf("badge key %q contains invalid characters", key)
	}
	if strings.ContainsAny(value, "\n\r") {
		return fmt.Errorf("badge value must not contain newlines")
	}
	st.Set(badgeKey(env, key), value)
	return nil
}

// GetBadge returns the value of a badge for an environment.
func GetBadge(st *state.State, env, key string) (string, bool) {
	v, ok := st.Get(badgeKey(env, key))
	return v, ok
}

// RemoveBadge removes a badge from an environment.
func RemoveBadge(st *state.State, env, key string) error {
	if !hasEnv(st, env) {
		return fmt.Errorf("environment %q not found", env)
	}
	st.Delete(badgeKey(env, key))
	return nil
}

// ListBadges returns all badges for an environment as a map.
func ListBadges(st *state.State, env string) map[string]string {
	prefix := fmt.Sprintf("%s%s:", badgePrefix, env)
	badges := map[string]string{}
	for _, k := range st.Keys() {
		if strings.HasPrefix(k, prefix) {
			key := strings.TrimPrefix(k, prefix)
			if v, ok := st.Get(k); ok {
				badges[key] = v
			}
		}
	}
	return badges
}
