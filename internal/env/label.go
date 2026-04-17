package env

import (
	"fmt"
	"regexp"

	"github.com/patchwork-deploy/internal/state"
)

var labelKeyRe = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)

// SetLabel attaches a key=value label to an environment in the state.
func SetLabel(st *state.State, env, key, value string) error {
	if env == "" {
		return fmt.Errorf("environment name must not be empty")
	}
	if !labelKeyRe.MatchString(key) {
		return fmt.Errorf("label key %q contains invalid characters", key)
	}
	records := st.ForEnvironment(env)
	if len(records) == 0 {
		return fmt.Errorf("environment %q not found", env)
	}
	sentinel := fmt.Sprintf("__label__%s__%s", key, value)
	// remove existing label with same key
	st.Remove(env, fmt.Sprintf("__label__%s__", key))
	st.Add(state.Record{
		Environment: env,
		Patch:       sentinel,
		AppliedAt:   records[0].AppliedAt,
	})
	return nil
}

// RemoveLabel removes a key label from an environment.
func RemoveLabel(st *state.State, env, key string) error {
	if env == "" {
		return fmt.Errorf("environment name must not be empty")
	}
	prefix := fmt.Sprintf("__label__%s__", key)
	removed := st.RemovePrefix(env, prefix)
	if !removed {
		return fmt.Errorf("label %q not found on environment %q", key, env)
	}
	return nil
}

// ListLabels returns all key=value labels attached to an environment.
func ListLabels(st *state.State, env string) (map[string]string, error) {
	records := st.ForEnvironment(env)
	if len(records) == 0 {
		return nil, fmt.Errorf("environment %q not found", env)
	}
	labels := map[string]string{}
	re := regexp.MustCompile(`^__label__([^_].+?)__(.+)$`)
	for _, r := range records {
		if m := re.FindStringSubmatch(r.Patch); m != nil {
			labels[m[1]] = m[2]
		}
	}
	return labels, nil
}
