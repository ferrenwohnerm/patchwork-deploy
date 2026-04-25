package env

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

var validSignalName = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func signalKey(env, patch, signal string) string {
	return fmt.Sprintf("signal::%s::%s::%s", env, patch, signal)
}

// SetSignal attaches a named signal to a patch within an environment.
// Signals are lightweight markers that indicate readiness, completion, or
// external acknowledgement for a given patch.
func SetSignal(st *state.State, env, patch, signal, value string) error {
	if !hasEnv(st, env) {
		return fmt.Errorf("environment %q not found", env)
	}
	if !patchExistsInEnv(st, env, patch) {
		return fmt.Errorf("patch %q not found in environment %q", patch, env)
	}
	if !validSignalName.MatchString(signal) {
		return fmt.Errorf("signal name %q contains invalid characters", signal)
	}
	if strings.ContainsAny(value, "\n\r") {
		return fmt.Errorf("signal value must not contain newlines")
	}
	st.SetMeta(signalKey(env, patch, signal), value)
	return nil
}

// GetSignal retrieves the value of a named signal for a patch.
func GetSignal(st *state.State, env, patch, signal string) (string, bool) {
	v, ok := st.GetMeta(signalKey(env, patch, signal))
	return v, ok
}

// RemoveSignal deletes a named signal from a patch.
func RemoveSignal(st *state.State, env, patch, signal string) {
	st.DeleteMeta(signalKey(env, patch, signal))
}

// ListSignals returns all signal name→value pairs for the given patch.
func ListSignals(st *state.State, env, patch string) map[string]string {
	prefix := fmt.Sprintf("signal::%s::%s::", env, patch)
	out := map[string]string{}
	for k, v := range st.AllMeta() {
		if strings.HasPrefix(k, prefix) {
			name := strings.TrimPrefix(k, prefix)
			out[name] = v
		}
	}
	return out
}
