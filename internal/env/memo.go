package env

import (
	"fmt"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

const memoPrefix = "memo:"

// SetMemo attaches a short freeform memo to an environment.
func SetMemo(st *state.State, env, memo string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	if strings.ContainsAny(memo, "\n\r") {
		return fmt.Errorf("memo must not contain newlines")
	}
	if len(memo) > 256 {
		return fmt.Errorf("memo exceeds 256 character limit")
	}
	st.SetMeta(env, memoPrefix+"text", memo)
	return nil
}

// GetMemo returns the memo for an environment, or an empty string if none is set.
func GetMemo(st *state.State, env string) (string, error) {
	if !st.HasEnvironment(env) {
		return "", fmt.Errorf("environment %q not found", env)
	}
	val := st.GetMeta(env, memoPrefix+"text")
	return val, nil
}

// ClearMemo removes the memo from an environment.
func ClearMemo(st *state.State, env string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	st.DeleteMeta(env, memoPrefix+"text")
	return nil
}
