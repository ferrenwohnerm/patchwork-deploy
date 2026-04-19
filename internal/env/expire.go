package env

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

const expirePrefix = "__expire__"

// SetExpiry sets an expiry time for an environment. After this time,
// the environment is considered expired.
func SetExpiry(st *state.State, env string, at time.Time) error {
	records := st.ForEnvironment(env)
	if len(records) == 0 {
		return fmt.Errorf("environment %q not found", env)
	}
	key := expirePrefix + env
	st.Set(key, at.UTC().Format(time.RFC3339))
	return nil
}

// GetExpiry returns the expiry time for an environment, if set.
func GetExpiry(st *state.State, env string) (time.Time, bool, error) {
	key := expirePrefix + env
	val, ok := st.GetMeta(key)
	if !ok {
		return time.Time{}, false, nil
	}
	t, err := time.Parse(time.RFC3339, val)
	if err != nil {
		return time.Time{}, false, fmt.Errorf("invalid expiry format: %w", err)
	}
	return t, true, nil
}

// ClearExpiry removes the expiry for an environment.
func ClearExpiry(st *state.State, env string) error {
	key := expirePrefix + env
	st.DeleteMeta(key)
	return nil
}

// IsExpired returns true if the environment has an expiry set and it is in the past.
func IsExpired(st *state.State, env string) (bool, error) {
	at, ok, err := GetExpiry(st, env)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	return time.Now().UTC().After(at), nil
}

// ListExpired returns all environments that have passed their expiry time.
func ListExpired(st *state.State) ([]string, error) {
	all := st.AllMeta()
	var expired []string
	for key, val := range all {
		if !strings.HasPrefix(key, expirePrefix) {
			continue
		}
		at, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return nil, errors.New("corrupt expiry entry: " + key)
		}
		if time.Now().UTC().After(at) {
			expired = append(expired, strings.TrimPrefix(key, expirePrefix))
		}
	}
	return expired, nil
}
