package diff

import "github.com/patchwork-deploy/internal/state"

// FromStateRecords converts state.Record slice into diff.Record slice.
func FromStateRecords(records []state.Record) []Record {
	out := make([]Record, 0, len(records))
	for _, r := range records {
		out = append(out, Record{
			Patch:       r.Patch,
			Environment: r.Environment,
			AppliedAt:   r.AppliedAt,
		})
	}
	return out
}
