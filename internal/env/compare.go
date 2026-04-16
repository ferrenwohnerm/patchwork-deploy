package env

import (
	"fmt"
	"sort"

	"github.com/patchwork-deploy/internal/state"
)

// CompareResult holds the patch-level comparison between two environments.
type CompareResult struct {
	Source      string
	Target      string
	OnlyInSource []string
	OnlyInTarget []string
	InBoth      []string
}

// CompareEnvironments compares applied patches between two environments.
func CompareEnvironments(st *state.State, source, target string) (CompareResult, error) {
	if source == target {
		return CompareResult{}, fmt.Errorf("source and target environments must differ")
	}

	srcRecs := st.ForEnvironment(source)
	if len(srcRecs) == 0 {
		return CompareResult{}, fmt.Errorf("source environment %q not found", source)
	}

	tgtRecs := st.ForEnvironment(target)
	if len(tgtRecs) == 0 {
		return CompareResult{}, fmt.Errorf("target environment %q not found", target)
	}

	srcIndex := make(map[string]struct{})
	for _, r := range srcRecs {
		srcIndex[r.Patch] = struct{}{}
	}

	tgtIndex := make(map[string]struct{})
	for _, r := range tgtRecs {
		tgtIndex[r.Patch] = struct{}{}
	}

	res := CompareResult{Source: source, Target: target}

	for p := range srcIndex {
		if _, ok := tgtIndex[p]; ok {
			res.InBoth = append(res.InBoth, p)
		} else {
			res.OnlyInSource = append(res.OnlyInSource, p)
		}
	}
	for p := range tgtIndex {
		if _, ok := srcIndex[p]; !ok {
			res.OnlyInTarget = append(res.OnlyInTarget, p)
		}
	}

	sort.Strings(res.InBoth)
	sort.Strings(res.OnlyInSource)
	sort.Strings(res.OnlyInTarget)

	return res, nil
}
