package env

import (
	"fmt"
	"strings"
)

// FormatCompareResult renders a CompareResult as a human-readable string.
func FormatCompareResult(res CompareResult) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "compare: %s vs %s\n", res.Source, res.Target)

	if len(res.InBoth) == 0 && len(res.OnlyInSource) == 0 && len(res.OnlyInTarget) == 0 {
		sb.WriteString("  no differences\n")
		return sb.String()
	}

	for _, p := range res.InBoth {
		fmt.Fprintf(&sb, "  = %s\n", p)
	}
	for _, p := range res.OnlyInSource {
		fmt.Fprintf(&sb, "  < %s (only in %s)\n", p, res.Source)
	}
	for _, p := range res.OnlyInTarget {
		fmt.Fprintf(&sb, "  > %s (only in %s)\n", p, res.Target)
	}

	return sb.String()
}

// SummaryLine returns a one-line summary of the comparison.
func (r CompareResult) SummaryLine() string {
	return fmt.Sprintf("%s vs %s: %d shared, %d only-source, %d only-target",
		r.Source, r.Target, len(r.InBoth), len(r.OnlyInSource), len(r.OnlyInTarget))
}
