package state

import (
	"fmt"
	"time"
)

// ImportOptions configures a bulk import of existing patch records.
type ImportOptions struct {
	PatchIDs    []string
	Environment string
	AppliedAt   time.Time
}

// Import bulk-adds patch records without duplicating existing ones.
// Returns the count of newly added records.
func Import(s *State, opts ImportOptions) (int, error) {
	if opts.Environment == "" {
		return 0, fmt.Errorf("environment must not be empty")
	}

	at := opts.AppliedAt
	if at.IsZero() {
		at = time.Now().UTC()
	}

	added := 0
	for _, id := range opts.PatchIDs {
		if id == "" {
			continue
		}
		if s.Has(id, opts.Environment) {
			continue
		}
		s.Records = append(s.Records, Record{
			PatchID:     id,
			Environment: opts.Environment,
			AppliedAt:   at,
		})
		added++
	}
	return added, nil
}

// Prune removes records for patches no longer present in knownIDs.
// Returns the count of removed records.
func Prune(s *State, environment string, knownIDs map[string]struct{}) int {
	var kept []Record
	removed := 0
	for _, r := range s.Records {
		if r.Environment != environment {
			kept = append(kept, r)
			continue
		}
		if _, ok := knownIDs[r.PatchID]; ok {
			kept = append(kept, r)
		} else {
			removed++
		}
	}
	s.Records = kept
	return removed
}
