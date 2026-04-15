package runner

import "fmt"

// Summary aggregates the results of a Run call into human-readable counts.
type Summary struct {
	Total   int
	Applied int
	Skipped int
	Failed  int
}

// Summarise builds a Summary from a slice of Results.
func Summarise(results []Result) Summary {
	var s Summary
	for _, r := range results {
		s.Total++
		switch {
		case r.Err != nil:
			s.Failed++
		case r.Applied:
			s.Applied++
		case r.Skipped:
			s.Skipped++
		}
	}
	return s
}

// String returns a one-line summary suitable for CLI output.
func (s Summary) String() string {
	return fmt.Sprintf("total=%d applied=%d skipped=%d failed=%d",
		s.Total, s.Applied, s.Skipped, s.Failed)
}
