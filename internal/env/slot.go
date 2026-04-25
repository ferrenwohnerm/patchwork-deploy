package env

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

var validSlotName = regexp.MustCompile(`^[a-z0-9_-]+$`)

func slotKey(env, patch, slot string) string {
	return fmt.Sprintf("slot:%s:%s:%s", env, patch, slot)
}

// SetSlot assigns a named deployment slot to a patch within an environment.
func SetSlot(st *state.State, env, patch, slot, value string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	if !patchExistsInEnv(st, env, patch) {
		return fmt.Errorf("patch %q not found in environment %q", patch, env)
	}
	if !validSlotName.MatchString(slot) {
		return fmt.Errorf("slot name %q is invalid: use only lowercase letters, digits, hyphens, underscores", slot)
	}
	if strings.ContainsAny(value, "\n\r") {
		return fmt.Errorf("slot value must not contain newlines")
	}
	st.SetMeta(slotKey(env, patch, slot), value)
	return nil
}

// GetSlot returns the value of a named slot for a patch in an environment.
func GetSlot(st *state.State, env, patch, slot string) (string, bool) {
	v, ok := st.GetMeta(slotKey(env, patch, slot))
	return v, ok
}

// RemoveSlot deletes a named slot assignment.
func RemoveSlot(st *state.State, env, patch, slot string) {
	st.DeleteMeta(slotKey(env, patch, slot))
}

// ListSlots returns all slot name→value pairs for a patch in an environment.
func ListSlots(st *state.State, env, patch string) map[string]string {
	prefix := fmt.Sprintf("slot:%s:%s:", env, patch)
	result := map[string]string{}
	for k, v := range st.AllMeta() {
		if strings.HasPrefix(k, prefix) {
			slot := strings.TrimPrefix(k, prefix)
			result[slot] = v
		}
	}
	return result
}
