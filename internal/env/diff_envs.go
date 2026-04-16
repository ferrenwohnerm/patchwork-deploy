package env

import (
	"fmt"

	"github.com/patchwork-deploy/internal/diff"
	"github.com/patchwork-deploy/internal/state"
)

// EnvDiffResult holds the comparison between two environments.
type EnvDiffResult struct {
	Source string
	Target string
	Result diff.Result
}

// DiffEnvironments compares the applied patches between two environments
// and returns a structured diff result.
func DiffEnvironments(st *state.State, source, target string) (EnvDiffResult, error) {
	if source == target {
		return EnvDiffResult{}, fmt.Errorf("source and target environments must differ")
	}

	srcRecords := st.ForEnvironment(source)
	if len(srcRecords) == 0 {
		return EnvDiffResult{}, fmt.Errorf("source environment %q not found or has no records", source)
	}

	tgtRecords := st.ForEnvironment(target)
	if len(tgtRecords) == 0 {
		return EnvDiffResult{}, fmt.Errorf("target environment %q not found or has no records", target)
	}

	srcDiff := diff.FromStateRecords(srcRecords)
	tgtDiff := diff.FromStateRecords(tgtRecords)

	result := diff.Compare(srcDiff, tgtDiff)

	return EnvDiffResult{
		Source: source,
		Target: target,
		Result: result,
	}, nil
}
