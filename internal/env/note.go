package env

import (
	"fmt"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

const notePrefix = "note:"

func noteKey(env string) string {
	return fmt.Sprintf("%s%s", notePrefix, env)
}

// SetNote attaches a freeform note to an environment.
func SetNote(st *state.State, env, text string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	if strings.ContainsAny(text, "\n\r") {
		return fmt.Errorf("note must not contain newlines")
	}
	if len(text) > 500 {
		return fmt.Errorf("note exceeds maximum length of 500 characters")
	}
	st.SetMeta(noteKey(env), text)
	return nil
}

// GetNote retrieves the note for an environment.
func GetNote(st *state.State, env string) (string, bool) {
	if !st.HasEnvironment(env) {
		return "", false
	}
	v, ok := st.GetMeta(noteKey(env))
	return v, ok
}

// ClearNote removes the note from an environment.
func ClearNote(st *state.State, env string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	st.DeleteMeta(noteKey(env))
	return nil
}
