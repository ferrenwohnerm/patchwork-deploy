package env

import (
	"fmt"
	"sort"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

// HistoryRecord represents a single applied patch entry for an environment.
type HistoryRecord struct {
	Patch     string
	AppliedAt time.Time
	Tags      []string
}

// History returns the ordered list of applied patches for the given environment.
func History(st *state.State, env string) ([]HistoryRecord, error) {
	records := st.ForEnvironment(env)
	if records == nil {
		return nil, fmt.Errorf("environment %q not found", env)
	}

	sorted := make([]state.Record, len(records))
	copy(sorted, records)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].AppliedAt.Before(sorted[j].AppliedAt)
	})

	out := make([]HistoryRecord, 0, len(sorted))
	for _, r := range sorted {
		out = append(out, HistoryRecord{
			Patch:     r.Patch,
			AppliedAt: r.AppliedAt,
			Tags:      r.Tags,
		})
	}
	return out, nil
}

// HistorySince returns patches applied on or after the given time.
func HistorySince(st *state.State, env string, since time.Time) ([]HistoryRecord, error) {
	all, err := History(st, env)
	if err != nil {
		return nil, err
	}
	var filtered []HistoryRecord
	for _, r := range all {
		if !r.AppliedAt.Before(since) {
			filtered = append(filtered, r)
		}
	}
	return filtered, nil
}
