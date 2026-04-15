package state

import (
	"fmt"
	"sort"
)

// RollbackPlan describes which patches should be reverted for an environment.
type RollbackPlan struct {
	Environment string
	Patches     []string // patch IDs in reverse-application order
}

// BuildRollbackPlan constructs a rollback plan for the given environment,
// targeting all patches applied after (and including) fromPatch.
// If fromPatch is empty the full history is included.
func BuildRollbackPlan(s *State, env, fromPatch string) (*RollbackPlan, error) {
	records := s.ForEnvironment(env)
	if len(records) == 0 {
		return nil, fmt.Errorf("no applied patches found for environment %q", env)
	}

	// Sort ascending by applied-at timestamp so we can find the cut-point.
	sort.Slice(records, func(i, j int) bool {
		return records[i].AppliedAt.Before(records[j].AppliedAt)
	})

	start := 0
	if fromPatch != "" {
		found := false
		for i, r := range records {
			if r.Patch == fromPatch {
				start = i
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("patch %q not found in applied history for environment %q", fromPatch, env)
		}
	}

	subset := records[start:]
	patches := make([]string, len(subset))
	for i, r := range subset {
		patches[i] = r.Patch
	}

	// Reverse so the most-recently-applied patch comes first.
	for l, r := 0, len(patches)-1; l < r; l, r = l+1, r-1 {
		patches[l], patches[r] = patches[r], patches[l]
	}

	return &RollbackPlan{
		Environment: env,
		Patches:     patches,
	}, nil
}
