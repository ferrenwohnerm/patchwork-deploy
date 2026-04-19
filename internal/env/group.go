package env

import (
	"fmt"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

const groupPrefix = "__group__"

func groupKey(name string) string {
	return groupPrefix + name
}

// AddToGroup adds an environment to a named group.
func AddToGroup(st *state.State, group, env string) error {
	if group == "" {
		return fmt.Errorf("group name must not be empty")
	}
	if !validNameRe.MatchString(group) {
		return fmt.Errorf("group name %q contains invalid characters", group)
	}
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	key := groupKey(group)
	members := listGroupMembers(st, key)
	for _, m := range members {
		if m == env {
			return nil // idempotent
		}
	}
	members = append(members, env)
	st.SetMeta(key, strings.Join(members, ","))
	return nil
}

// RemoveFromGroup removes an environment from a named group.
func RemoveFromGroup(st *state.State, group, env string) error {
	key := groupKey(group)
	members := listGroupMembers(st, key)
	updated := members[:0]
	for _, m := range members {
		if m != env {
			updated = append(updated, m)
		}
	}
	if len(updated) == 0 {
		st.DeleteMeta(key)
	} else {
		st.SetMeta(key, strings.Join(updated, ","))
	}
	return nil
}

// ListGroup returns all environments in a named group.
func ListGroup(st *state.State, group string) ([]string, error) {
	if !st.HasMeta(groupKey(group)) {
		return nil, fmt.Errorf("group %q not found", group)
	}
	return listGroupMembers(st, groupKey(group)), nil
}

// ListAllGroups returns all group names.
func ListAllGroups(st *state.State) []string {
	var groups []string
	for _, k := range st.MetaKeys() {
		if strings.HasPrefix(k, groupPrefix) {
			groups = append(groups, strings.TrimPrefix(k, groupPrefix))
		}
	}
	return groups
}

func listGroupMembers(st *state.State, key string) []string {
	v, ok := st.GetMeta(key)
	if !ok || v == "" {
		return nil
	}
	return strings.Split(v, ",")
}
