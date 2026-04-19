package env

import (
	"fmt"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

const triggerPrefix = "trigger:"

type Trigger struct {
	Patch string
	Event string
	Action string
}

func SetTrigger(st *state.State, env, patch, event, action string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	if !patchExistsInEnv(st, env, patch) {
		return fmt.Errorf("patch %q not found in environment %q", patch, env)
	}
	event = strings.TrimSpace(event)
	if event == "" {
		return fmt.Errorf("event must not be empty")
	}
	action = strings.TrimSpace(action)
	if action == "" {
		return fmt.Errorf("action must not be empty")
	}
	key := fmt.Sprintf("%s%s:%s:%s", triggerPrefix, patch, event, action)
	st.SetMeta(env, key, "1")
	return nil
}

func ListTriggers(st *state.State, env string) ([]Trigger, error) {
	if !st.HasEnvironment(env) {
		return nil, fmt.Errorf("environment %q not found", env)
	}
	var triggers []Trigger
	for k := range st.MetaFor(env) {
		if !strings.HasPrefix(k, triggerPrefix) {
			continue
		}
		parts := strings.SplitN(strings.TrimPrefix(k, triggerPrefix), ":", 3)
		if len(parts) != 3 {
			continue
		}
		triggers = append(triggers, Trigger{Patch: parts[0], Event: parts[1], Action: parts[2]})
	}
	return triggers, nil
}

func RemoveTrigger(st *state.State, env, patch, event string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	prefix := fmt.Sprintf("%s%s:%s:", triggerPrefix, patch, event)
	for k := range st.MetaFor(env) {
		if strings.HasPrefix(k, prefix) {
			st.DeleteMeta(env, k)
		}
	}
	return nil
}
