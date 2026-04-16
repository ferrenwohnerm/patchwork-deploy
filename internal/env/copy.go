package env

import (
	"fmt"

	"github.com/patchwork-deploy/internal/state"
)

// CopyOptions controls which patches are copied between environments.
type CopyOptions struct {
	PatchIDs []string // if empty, copies all patches from source
}

// CopyResult describes the outcome of a Copy operation.
type CopyResult struct {
	Copied  []string
	Skipped []string
}

// Copy duplicates specific (or all) patch records from one environment to
// another without removing the source records.
func Copy(st *state.State, src, dst string, opts CopyOptions) (CopyResult, error) {
	if src == dst {
		return CopyResult{}, fmt.Errorf("source and destination environment must differ: %q", src)
	}

	srcRecords := st.ForEnvironment(src)
	if len(srcRecords) == 0 {
		return CopyResult{}, fmt.Errorf("source environment %q has no records", src)
	}

	// Build a lookup of patches already present in dst.
	dstRecords := st.ForEnvironment(dst)
	alreadyInDst := make(map[string]bool, len(dstRecords))
	for _, r := range dstRecords {
		alreadyInDst[r.PatchID] = true
	}

	// Determine which patches to copy.
	wanted := make(map[string]bool)
	if len(opts.PatchIDs) > 0 {
		for _, id := range opts.PatchIDs {
			wanted[id] = true
		}
	}

	var result CopyResult
	for _, r := range srcRecords {
		if len(wanted) > 0 && !wanted[r.PatchID] {
			continue
		}
		if alreadyInDst[r.PatchID] {
			result.Skipped = append(result.Skipped, r.PatchID)
			continue
		}
		copy := r
		copy.Environment = dst
		st.Add(copy)
		result.Copied = append(result.Copied, r.PatchID)
	}

	return result, nil
}
