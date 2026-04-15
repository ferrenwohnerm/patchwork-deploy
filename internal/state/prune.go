package state

import "time"

// PruneOptions controls which records are removed during a prune operation.
type PruneOptions struct {
	// OlderThan removes records applied before this time. Zero means no time filter.
	OlderThan time.Time
	// Environment restricts pruning to a specific environment. Empty means all.
	Environment string
	// DryRun reports what would be removed without modifying state.
	DryRun bool
}

// PruneResult holds the outcome of a prune operation.
type PruneResult struct {
	Removed []Record
	Retained []Record
}

// PruneByAge removes state records matching the given options.
// It returns a PruneResult describing what was (or would be) removed.
func PruneByAge(st *State, opts PruneOptions) PruneResult {
	var removed, retained []Record

	for _, rec := range st.Records {
		matchEnv := opts.Environment == "" || rec.Environment == opts.Environment
		matchAge := !opts.OlderThan.IsZero() && rec.AppliedAt.Before(opts.OlderThan)

		if matchEnv && matchAge {
			removed = append(removed, rec)
		} else {
			retained = append(retained, rec)
		}
	}

	if !opts.DryRun {
		st.Records = retained
	}

	return PruneResult{
		Removed:  removed,
		Retained: retained,
	}
}
