package env

import (
	"fmt"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

// ArchiveResult holds the outcome of an archive operation.
type ArchiveResult struct {
	Environment string
	RecordsArchived int
	ArchivedAt      time.Time
}

// Archive marks all state records for an environment as archived by copying
// them into a namespaced archive environment key ("<env>:archived") and
// removing them from the live environment.
func Archive(st *state.State, env string) (ArchiveResult, error) {
	records := st.ForEnvironment(env)
	if len(records) == 0 {
		return ArchiveResult{}, fmt.Errorf("archive: no records found for environment %q", env)
	}

	archiveKey := env + ":archived"
	for _, r := range records {
		ar := state.Record{
			Environment: archiveKey,
			Patch:       r.Patch,
			AppliedAt:   r.AppliedAt,
			Tags:        r.Tags,
		}
		st.Add(ar)
	}

	st.RemoveEnvironment(env)

	return ArchiveResult{
		Environment:     env,
		RecordsArchived: len(records),
		ArchivedAt:      time.Now().UTC(),
	}, nil
}

// ListArchived returns all archived records for the given environment.
func ListArchived(st *state.State, env string) []state.Record {
	return st.ForEnvironment(env + ":archived")
}
