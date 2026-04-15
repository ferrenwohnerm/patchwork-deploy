package env

import (
	"fmt"

	"github.com/patchwork-deploy/internal/state"
)

// RenameResult holds the outcome of a rename operation.
type RenameResult struct {
	OldName  string
	NewName  string
	Migrated int
}

// Rename copies all state records from srcEnv to dstEnv, then removes
// the srcEnv records from the state file. It returns an error if srcEnv
// does not exist, dstEnv already exists, or the names are equal.
func Rename(st *state.State, srcEnv, dstEnv string) (RenameResult, error) {
	if srcEnv == dstEnv {
		return RenameResult{}, fmt.Errorf("source and destination environment names are identical: %q", srcEnv)
	}

	srcRecords := st.ForEnvironment(srcEnv)
	if len(srcRecords) == 0 {
		return RenameResult{}, fmt.Errorf("source environment %q not found or has no records", srcEnv)
	}

	dstRecords := st.ForEnvironment(dstEnv)
	if len(dstRecords) > 0 {
		return RenameResult{}, fmt.Errorf("destination environment %q already exists with %d record(s)", dstEnv, len(dstRecords))
	}

	for _, r := range srcRecords {
		st.Add(state.Record{
			Environment: dstEnv,
			Patch:       r.Patch,
			AppliedAt:   r.AppliedAt,
		})
	}

	st.RemoveEnvironment(srcEnv)

	return RenameResult{
		OldName:  srcEnv,
		NewName:  dstEnv,
		Migrated: len(srcRecords),
	}, nil
}
