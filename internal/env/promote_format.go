package env

import (
	"fmt"
	"strings"
)

// PromoteResult holds the outcome of a promotion between environments.
type PromoteResult struct {
	Source  string
	Target  string
	Applied []string
	Skipped []string
}

// FormatPromoteResult returns a human-readable summary of a promotion.
func FormatPromoteResult(r PromoteResult) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Promote: %s → %s\n", r.Source, r.Target)
	if len(r.Applied) == 0 && len(r.Skipped) == 0 {
		sb.WriteString("  no patches to promote\n")
		return sb.String()
	}
	for _, p := range r.Applied {
		fmt.Fprintf(&sb, "  + %s\n", p)
	}
	for _, p := range r.Skipped {
		fmt.Fprintf(&sb, "  ~ %s (already applied)\n", p)
	}
	fmt.Fprintf(&sb, "Summary: %d applied, %d skipped\n", len(r.Applied), len(r.Skipped))
	return sb.String()
}

// HasChanges reports whether the promotion applied at least one patch.
func (r PromoteResult) HasChanges() bool {
	return len(r.Applied) > 0
}
