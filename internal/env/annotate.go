package env

import (
	"fmt"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

const annotatePrefix = "annotation:"

// SetAnnotation attaches a free-text annotation to a patch record in the given environment.
func SetAnnotation(st *state.State, env, patch, text string) error {
	recs := st.ForEnvironment(env)
	if len(recs) == 0 {
		return fmt.Errorf("environment %q not found", env)
	}
	found := false
	for _, r := range recs {
		if r.Patch == patch {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("patch %q not found in environment %q", patch, env)
	}
	if strings.ContainsAny(text, "\n\r") {
		return fmt.Errorf("annotation text must not contain newlines")
	}
	key := fmt.Sprintf("%s%s::%s", annotatePrefix, env, patch)
	st.SetMeta(key, text)
	return nil
}

// GetAnnotation returns the annotation for a patch in the given environment.
func GetAnnotation(st *state.State, env, patch string) (string, bool) {
	key := fmt.Sprintf("%s%s::%s", annotatePrefix, env, patch)
	v, ok := st.GetMeta(key)
	return v, ok
}

// RemoveAnnotation clears the annotation for a patch in the given environment.
func RemoveAnnotation(st *state.State, env, patch string) {
	key := fmt.Sprintf("%s%s::%s", annotatePrefix, env, patch)
	st.DeleteMeta(key)
}

// ListAnnotations returns all patch->annotation pairs for the given environment.
func ListAnnotations(st *state.State, env string) map[string]string {
	prefix := fmt.Sprintf("%s%s::", annotatePrefix, env)
	out := map[string]string{}
	for k, v := range st.AllMeta() {
		if strings.HasPrefix(k, prefix) {
			patch := strings.TrimPrefix(k, prefix)
			out[patch] = v
		}
	}
	return out
}
