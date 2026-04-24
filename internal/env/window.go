package env

import (
	"fmt"
	"strings"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

const windowKeyPrefix = "window:"

func windowKey(env, patch string) string {
	return fmt.Sprintf("%s%s:%s", windowKeyPrefix, env, patch)
}

// SetWindow defines an allowed deployment window for a patch in an environment.
// start and end are in "15:04" format (24h), e.g. "08:00" and "18:00".
func SetWindow(st *state.State, env, patch, start, end string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	if !patchExistsInEnv(st, env, patch) {
		return fmt.Errorf("patch %q not found in environment %q", patch, env)
	}
	if _, err := time.Parse("15:04", start); err != nil {
		return fmt.Errorf("invalid start time %q: must be HH:MM", start)
	}
	if _, err := time.Parse("15:04", end); err != nil {
		return fmt.Errorf("invalid end time %q: must be HH:MM", end)
	}
	if start >= end {
		return fmt.Errorf("start time %q must be before end time %q", start, end)
	}
	st.SetMeta(windowKey(env, patch), start+"-"+end)
	return nil
}

// GetWindow returns the deployment window for a patch, or empty strings if none set.
func GetWindow(st *state.State, env, patch string) (start, end string, ok bool) {
	v, exists := st.GetMeta(windowKey(env, patch))
	if !exists {
		return "", "", false
	}
	parts := strings.SplitN(v, "-", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	return parts[0], parts[1], true
}

// ClearWindow removes the deployment window for a patch.
func ClearWindow(st *state.State, env, patch string) {
	st.DeleteMeta(windowKey(env, patch))
}

// InWindow reports whether the current time falls within the patch's deployment window.
// Returns true if no window is set (unrestricted).
func InWindow(st *state.State, env, patch string) bool {
	start, end, ok := GetWindow(st, env, patch)
	if !ok {
		return true
	}
	now := time.Now().Format("15:04")
	return now >= start && now < end
}

// ListWindows returns all patch→window mappings for an environment.
func ListWindows(st *state.State, env string) map[string][2]string {
	prefix := windowKeyPrefix + env + ":"
	out := map[string][2]string{}
	for k, v := range st.AllMeta() {
		if !strings.HasPrefix(k, prefix) {
			continue
		}
		patch := strings.TrimPrefix(k, prefix)
		parts := strings.SplitN(v, "-", 2)
		if len(parts) == 2 {
			out[patch] = [2]string{parts[0], parts[1]}
		}
	}
	return out
}
