package env

import (
	"fmt"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

// PromoteResult holds the outcome of a promotion between environments.
type PromoteResult struct {
	FromEnv   string
	ToEnv     string
	Promoted  []string
	Skipped   []string
	AppliedAt time.Time
}

// Promote copies applied patches from one environment into another,
// skipping any that are already present in the target environment.
func Promote(st *state.State, fromEnv, toEnv string) (*PromoteResult, error) {
	if fromEnv == toEnv {
		return nil, fmt.Errorf("source and target environment must differ")
	}

	source := st.ForEnvironment(fromEnv)
	if len(source) == 0 {
		return nil, fmt.Errorf("no records found for environment %q", fromEnv)
	}

	target := st.ForEnvironment(toEnv)
	targetIndex := make(map[string]struct{}, len(target))
	for _, r := range target {
		targetIndex[r.Patch] = struct{}{}
	}

	result := &PromoteResult{
		FromEnv:   fromEnv,
		ToEnv:     toEnv,
		AppliedAt: time.Now().UTC(),
	}

	for _, r := range source {
		if _, exists := targetIndex[r.Patch]; exists {
			result.Skipped = append(result.Skipped, r.Patch)
			continue
		}
		st.Add(state.Record{
			Environment: toEnv,
			Patch:       r.Patch,
			AppliedAt:   result.AppliedAt,
		})
		result.Promoted = append(result.Promoted, r.Patch)
	}

	return result, nil
}
