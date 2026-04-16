package env

import (
	"fmt"
	"regexp"

	"github.com/patchwork-deploy/internal/state"
)

var validTag = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)

type TagEntry struct {
	Environment string
	Patch       string
	Tag         string
}

func AddTag(st *state.State, env, patch, tag string) error {
	if !validTag.MatchString(tag) {
		return fmt.Errorf("invalid tag %q: only alphanumeric, dash and underscore allowed", tag)
	}
	recs := st.ForEnvironment(env)
	if len(recs) == 0 {
		return fmt.Errorf("environment %q not found", env)
	}
	for _, r := range recs {
		if r.Patch == patch {
			r.Tag = tag
			st.Upsert(r)
			return nil
		}
	}
	return fmt.Errorf("patch %q not found in environment %q", patch, env)
}

func ListTags(st *state.State, env string) ([]TagEntry, error) {
	recs := st.ForEnvironment(env)
	if len(recs) == 0 {
		return nil, fmt.Errorf("environment %q not found", env)
	}
	var out []TagEntry
	for _, r := range recs {
		if r.Tag != "" {
			out = append(out, TagEntry{Environment: r.Environment, Patch: r.Patch, Tag: r.Tag})
		}
	}
	return out, nil
}

func RemoveTag(st *state.State, env, patch string) error {
	recs := st.ForEnvironment(env)
	for _, r := range recs {
		if r.Patch == patch {
			r.Tag = ""
			st.Upsert(r)
			return nil
		}
	}
	return fmt.Errorf("patch %q not found in environment %q", patch, env)
}
