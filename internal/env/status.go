package env

import (
	"fmt"
	"sort"

	"github.com/patchwork-deploy/internal/state"
)

// EnvStatus summarises the deployment state for a single environment.
type EnvStatus struct {
	Environment string
	Total       int
	Latest      string
	Tags        []string
}

func (s EnvStatus) String() string {
	if s.Total == 0 {
		return fmt.Sprintf("%s: no patches applied", s.Environment)
	}
	return fmt.Sprintf("%s: %d patch(es) applied, latest=%s", s.Environment, s.Total, s.Latest)
}

// Status returns an EnvStatus for the given environment.
func Status(st *state.State, env string) (EnvStatus, error) {
	records := st.ForEnvironment(env)
	if records == nil {
		return EnvStatus{}, fmt.Errorf("environment %q not found in state", env)
	}

	status := EnvStatus{Environment: env, Total: len(records)}
	if len(records) == 0 {
		return status, nil
	}

	// sort by applied-at to find latest
	sort.Slice(records, func(i, j int) bool {
		return records[i].AppliedAt.Before(records[j].AppliedAt)
	})
	status.Latest = records[len(records)-1].Patch

	// collect unique tags
	tagSet := map[string]struct{}{}
	for _, r := range records {
		for _, t := range r.Tags {
			tagSet[t] = struct{}{}
		}
	}
	for t := range tagSet {
		status.Tags = append(status.Tags, t)
	}
	sort.Strings(status.Tags)
	return status, nil
}

// StatusAll returns an EnvStatus for every environment present in state.
func StatusAll(st *state.State) []EnvStatus {
	envs := st.Environments()
	sort.Strings(envs)
	out := make([]EnvStatus, 0, len(envs))
	for _, e := range envs {
		s, err := Status(st, e)
		if err == nil {
			out = append(out, s)
		}
	}
	return out
}
