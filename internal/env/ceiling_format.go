package env

import (
	"fmt"
	"sort"
	"strings"
)

// CeilingResult holds a single ceiling entry with its current usage.
type CeilingResult struct {
	Patch   string
	Limit   int
	Current int
	Exceeded bool
}

// FormatCeilings renders a list of ceiling results as a human-readable table.
func FormatCeilings(results []CeilingResult) string {
	if len(results) == 0 {
		return "no ceilings configured"
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Patch < results[j].Patch
	})

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-30s %8s %8s %s\n", "PATCH", "LIMIT", "CURRENT", "STATUS"))
	for _, r := range results {
		status := "ok"
		if r.Exceeded {
			status = "EXCEEDED"
		}
		sb.WriteString(fmt.Sprintf("%-30s %8d %8d %s\n", r.Patch, r.Limit, r.Current, status))
	}
	return strings.TrimRight(sb.String(), "\n")
}
